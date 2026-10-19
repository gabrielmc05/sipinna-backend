// handlers.go es para definir los controladores de la API REST relacionados con los reportes,
// incluyendo la creación, consulta y actualización de reportes, así como la carga de evidencias
// y la obtención de instituciones.

package handlers

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/token"
	"github.com/gin-gonic/gin"
)

const (
	maxEvidenceSize = 50 << 20
	multipartMemory = 8 << 20
	privateCodeSize = 8
	folioAlphabet   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type reportRequest struct {
	UserID          *uint32        `json:"id_usuario"`
	ContactEmail    *string        `json:"correo_contacto"`
	Description     string         `json:"descripcion" binding:"required"`
	ApproximateTime string         `json:"horario_aproximado"`
	Anonymous       bool           `json:"es_anonimo"`
	WantsFollowUp   bool           `json:"desea_seguimiento"`
	AdultInvolved   bool           `json:"hay_adulto_involucrado"`
	ImmediateDanger string         `json:"peligro_inmediato" binding:"required"`
	Location        locationInput  `json:"ubicacion" binding:"required"`
	Minors          []dboper.Minor `json:"menores" binding:"required,min=1"`
	Tags            []dboper.Tag   `json:"etiquetas" binding:"required"`
}

type locationInput struct {
	Latitude    *float64 `json:"latitud" binding:"required"`
	Longitude   *float64 `json:"longitud" binding:"required"`
	Address     string   `json:"direccion_texto"`
	TerritoryID uint32   `json:"id_territorio" binding:"required"`
}

type reportConsultRequest struct {
	Folio       string `json:"folio" binding:"required"`
	PrivateCode string `json:"codigo_acceso" binding:"required"`
}

type statusRequest struct {
	Status       string  `json:"estado" binding:"required"`
	Observations *string `json:"observaciones"`
}

type assignReportRequest struct {
	FeederID uint32 `json:"id_alimentador" binding:"required"`
}

type institutionCreateRequest struct {
	TerritoryID uint32  `json:"id_territorio" binding:"required"`
	Name        string  `json:"nombre" binding:"required"`
	Type        string  `json:"tipo" binding:"required"`
	Email       *string `json:"correo"`
	Phone       *string `json:"telefono"`
	Address     *string `json:"direccion"`
}

type institutionUpdateRequest struct {
	TerritoryID *uint32 `json:"id_territorio"`
	Name        *string `json:"nombre"`
	Type        *string `json:"tipo"`
	Email       *string `json:"correo"`
	Phone       *string `json:"telefono"`
	Address     *string `json:"direccion"`
}

type institutionStatusRequest struct {
	Active *bool `json:"es_activo" binding:"required"`
}

// CreateReport maneja la creación de un nuevo reporte.
// Valida la solicitud, genera un folio único y un código privado,
// y guarda el reporte en la base de datos.
func CreateReport(c *gin.Context) {
	var request reportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de reporte inválido"})
		return
	}
	if err := validateReportRequest(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.UserID != nil && *request.UserID != 0 {
		userID, err := token.Authenticate(c.GetHeader("Authorization"))
		if errors.Is(err, token.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "se requiere un token válido del usuario indicado"})
			return
		}
		if err != nil {
			log.Printf("authenticate report owner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible validar el usuario"})
			return
		}
		if userID != *request.UserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "el usuario autenticado no coincide con id_usuario"})
			return
		}
	} else {
		request.UserID = nil
	}

	privateCode, err := randomString(privateCodeSize, folioAlphabet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible generar el código de acceso"})
		return
	}
	codeHash := sha256.Sum256([]byte(privateCode))
	requestData := dboper.ReportCreation{
		UserID:          request.UserID,
		ContactEmail:    request.ContactEmail,
		Description:     request.Description,
		ApproximateTime: optionalString(request.ApproximateTime),
		Anonymous:       request.Anonymous,
		WantsFollowUp:   request.WantsFollowUp,
		AdultInvolved:   request.AdultInvolved,
		ImmediateDanger: request.ImmediateDanger,
		Latitude:        *request.Location.Latitude,
		Longitude:       *request.Location.Longitude,
		Address:         request.Location.Address,
		TerritoryID:     request.Location.TerritoryID,
		Minors:          request.Minors,
	}
	var folio string
	for range 3 {
		folio, err = newFolio()
		if err != nil {
			log.Printf("generate report folio: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible generar el folio"})
			return
		}
		err = dboper.CreateReport(requestData, folio, hex.EncodeToString(codeHash[:]))
		if err == nil {
			break
		}
		if !dboper.IsDuplicateKeyError(err) {
			log.Printf("create report: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible guardar el reporte"})
			return
		}
	}
	if err != nil {
		log.Printf("could not allocate a unique report folio after retries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible generar un folio único"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"folio":          folio,
		"codigo_privado": privateCode,
	})
}

// UploadEvidence maneja la carga de evidencias para un reporte específico.
// Valida el folio, procesa los archivos cargados y los guarda en el sistema de archivos,
// registrando la información en la base de datos.
func UploadEvidence(c *gin.Context) {
	folio := strings.TrimSpace(c.Param("folio"))
	if folio == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folio requerido"})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxEvidenceSize)
	if err := c.Request.ParseMultipartForm(multipartMemory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formulario multipart inválido o supera el límite permitido"})
		return
	}
	defer c.Request.MultipartForm.RemoveAll()
	files := c.Request.MultipartForm.File["archivos"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se requiere al menos un archivo en el campo archivos"})
		return
	}
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.Join("uploads", "evidencias")
	}
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible preparar el almacenamiento de evidencias"})
		return
	}

	evidence := make([]dboper.Evidence, 0, len(files))
	cleanup := func() {
		for _, item := range evidence {
			if err := os.Remove(item.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
				log.Printf("remove unlinked evidence after upload failure: %v", err)
			}
		}
	}
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			cleanup()
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible abrir uno de los archivos"})
			return
		}
		contentType, err := detectEvidenceType(file)
		if err != nil {
			file.Close()
			cleanup()
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			file.Close()
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible leer uno de los archivos"})
			return
		}

		fileName, err := randomString(32, folioAlphabet)
		if err != nil {
			file.Close()
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible generar el nombre del archivo"})
			return
		}
		path := filepath.Join(uploadDir, fileName+extensionForType(contentType))
		destination, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			file.Close()
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible guardar la evidencia"})
			return
		}
		_, copyErr := io.Copy(destination, file)
		closeDestinationErr := destination.Close()
		closeFileErr := file.Close()
		if copyErr != nil || closeDestinationErr != nil || closeFileErr != nil {
			if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				log.Printf("remove incomplete evidence upload: %v", err)
			}
			cleanup()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible guardar la evidencia"})
			return
		}
		fileType := "Imagen"
		if strings.HasPrefix(contentType, "video/") {
			fileType = "Video"
		}
		evidence = append(evidence, dboper.Evidence{Type: fileType, Path: path})
	}

	if err := dboper.CreateEvidences(folio, evidence); err != nil {
		log.Printf("register report evidences: %v", err)
		cleanup()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible registrar las evidencias"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "Evidencias guardadas y vinculadas exitosamente."})
}

// ConsultFolio maneja la consulta de un reporte específico utilizando el folio y
// el token de autorización del usuario.
// Valida la solicitud y devuelve la información del reporte si el usuario tiene acceso.
// Si el reporte no se encuentra o el usuario no tiene acceso, devuelve un error adecuado.
func ConsultFolio(c *gin.Context) {
	userID, err := token.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorización inválido o faltante"})
			return
		}
		log.Printf("authenticate report list request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible validar la autorización"})
		return
	}

	folio := c.Param("folio")
	report, err := dboper.ConsultFolio(folio, userID)
	if err != nil {
		if errors.Is(err, dboper.ErrReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No fue posible consultar el reporte (no autorizado)."})
		return
	}

	c.JSON(http.StatusOK, report)
}

// ConsultReport maneja la consulta de un reporte específico utilizando el folio y un código privado.
// Valida la solicitud y devuelve la información del reporte si el código es correcto.
// Si el reporte no se encuentra o el código es incorrecto, devuelve un error adecuado.
func ConsultReport(c *gin.Context) {
	var request reportConsultRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folio y codigo_acceso son requeridos"})
		return
	}
	hash := sha256.Sum256([]byte(request.PrivateCode))
	report, err := dboper.ConsultReport(request.Folio, hex.EncodeToString(hash[:]))
	if err != nil {
		if errors.Is(err, dboper.ErrReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar el reporte"})
		return
	}
	c.JSON(http.StatusOK, report)
}

// ListUserReports maneja la consulta de todos los reportes asociados a un usuario autenticado.
// Valida el token de autorización y devuelve la lista de reportes si el usuario es válido.
// Si el token es inválido o no se puede consultar la base de datos, devuelve un error adecuado.
func ListUserReports(c *gin.Context) {
	userID, err := token.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorización inválido o faltante"})
			return
		}
		log.Printf("authenticate report list request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible validar la autorización"})
		return
	}
	reports, err := dboper.ListUserReports(userID)
	if err != nil {
		log.Printf("list user reports: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar los reportes"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// ReportHistory devuelve el histórico global para administradores y solo los casos asignados para alimentadores.
func ReportHistory(c *gin.Context) {
	userID, err := authenticateStaff(c)
	if err != nil {
		return
	}
	role, err := dboper.GetStaffRole(userID)
	if err != nil {
		respondStaffRoleError(c, err)
		return
	}
	if role != "Administrador" && role != "Alimentador" {
		c.JSON(http.StatusForbidden, gin.H{"error": "rol sin permisos para consultar el histórico"})
		return
	}

	var status *string
	if value := strings.TrimSpace(c.Query("estado")); value != "" {
		if !isHistoryStatus(value) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "estado de expediente inválido"})
			return
		}
		status = &value
	}

	var territoryID *uint32
	if value := c.Query("id_territorio"); value != "" {
		parsedID, err := strconv.ParseUint(value, 10, 32)
		if err != nil || parsedID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id_territorio debe ser un entero positivo"})
			return
		}
		parsedTerritoryID := uint32(parsedID)
		territoryID = &parsedTerritoryID
	}

	dateOrder := strings.ToLower(strings.TrimSpace(c.DefaultQuery("orden_fecha", "desc")))
	if dateOrder != "asc" && dateOrder != "desc" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orden_fecha debe ser asc o desc"})
		return
	}
	reports, err := dboper.GetReportHistory(userID, status, territoryID, dateOrder)
	if err != nil {
		log.Printf("get report history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar el histórico"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// AssignReport asigna un reporte a un alimentador; la base de datos valida el rol de ambos usuarios.
func AssignReport(c *gin.Context) {
	userID, err := authenticateStaff(c)
	if err != nil {
		return
	}
	var request assignReportRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.FeederID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_alimentador debe ser un entero positivo"})
		return
	}

	err = dboper.AssignReport(userID, c.Param("folio"), request.FeederID)
	switch {
	case errors.Is(err, dboper.ErrReportNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
	case errors.Is(err, dboper.ErrAdminRoleRequired):
		c.JSON(http.StatusForbidden, gin.H{"error": "solo un administrador puede asignar reportes"})
	case errors.Is(err, dboper.ErrFeederRoleRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_alimentador no corresponde a un alimentador registrado"})
	case err != nil:
		log.Printf("assign report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible asignar el reporte"})
	default:
		c.JSON(http.StatusOK, gin.H{"mensaje": "Reporte asignado correctamente"})
	}
}

// GetEvidence devuelve las evidencias si el usuario es administrador o tiene asignado el reporte.
func GetEvidence(c *gin.Context) {
	userID, err := authenticateStaff(c)
	if err != nil {
		return
	}
	evidences, err := dboper.GetReportEvidences(userID, c.Param("folio"))
	switch {
	case errors.Is(err, dboper.ErrReportNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
	case errors.Is(err, dboper.ErrReportAccessDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": "no tienes permiso para consultar las evidencias"})
	case err != nil:
		log.Printf("get report evidences: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar las evidencias"})
	default:
		c.JSON(http.StatusOK, evidences)
	}
}

/*
DownloadEvidence entrega el archivo de una evidencia registrada.

Verifica que:
  - El usuario esté autenticado y autorizado.
  - La evidencia pertenezca al expediente solicitado.
  - La ruta corresponda al directorio de evidencias.
  - El archivo exista físicamente.

El archivo se entrega sin publicar el directorio de almacenamiento.
*/
func DownloadEvidence(c *gin.Context) {
	userID, err := authenticateStaff(c)
	if err != nil {
		return
	}

	folio := strings.TrimSpace(c.Param("folio"))

	evidenceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || evidenceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id de evidencia inválido",
		})
		return
	}

	evidences, err := dboper.GetReportEvidences(userID, folio)

	switch {
	case errors.Is(err, dboper.ErrReportNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "reporte no encontrado",
		})
		return

	case errors.Is(err, dboper.ErrReportAccessDenied):
		c.JSON(http.StatusForbidden, gin.H{
			"error": "no tienes permiso para consultar esta evidencia",
		})
		return

	case err != nil:
		log.Printf("download evidence authorization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no fue posible consultar la evidencia",
		})
		return
	}

	var evidencePath string

	for _, evidence := range evidences {
		if uint64(evidence.ID) == evidenceID {
			evidencePath = evidence.Path
			break
		}
	}

	if evidencePath == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "evidencia no encontrada",
		})
		return
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.Join("uploads", "evidencias")
	}

	basePath, err := filepath.Abs(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "almacenamiento no disponible",
		})
		return
	}

	filePath, err := filepath.Abs(evidencePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ruta de evidencia inválida",
		})
		return
	}

	relativePath, err := filepath.Rel(basePath, filePath)
	if err != nil ||
		relativePath == ".." ||
		filepath.IsAbs(relativePath) ||
		strings.HasPrefix(relativePath, ".."+string(os.PathSeparator)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "ruta de evidencia no permitida",
		})
		return
	}

	// Evitar seguir enlaces simbólicos.
	fileInfo, err := os.Lstat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "el archivo de evidencia ya no está disponible",
		})
		return
	}
	if err != nil {
		log.Printf("stat evidence file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no fue posible consultar el archivo",
		})
		return
	}
	if !fileInfo.Mode().IsRegular() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "archivo de evidencia no permitido",
		})
		return
	}

	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "private, no-store")
	c.Header("Content-Security-Policy", "sandbox")

	c.File(filePath)
}

// GetTerritories devuelve los estados disponibles y los municipios, con filtros opcionales.
func GetTerritories(c *gin.Context) {
	var state, municipality *string
	if value := strings.TrimSpace(c.Query("estado")); value != "" {
		state = &value
	}
	if value := strings.TrimSpace(c.Query("municipio")); value != "" {
		municipality = &value
	}

	options, err := dboper.GetTerritoryOptions(state, municipality)
	if err != nil {
		log.Printf("get territory options: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar los territorios"})
		return
	}
	c.JSON(http.StatusOK, options)
}

// GetTerritory devuelve un municipio registrado mediante su ID.
func GetTerritory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de territorio inválido"})
		return
	}
	territory, err := dboper.GetTerritory(uint32(id))
	if errors.Is(err, dboper.ErrTerritoryNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "territorio no encontrado"})
		return
	}
	if err != nil {
		log.Printf("get territory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar el territorio"})
		return
	}
	c.JSON(http.StatusOK, territory)
}

// GetFeeders devuelve las cuentas activas de alimentadores; solo está disponible para administradores.
func GetFeeders(c *gin.Context) {
	userID, err := authenticateStaff(c)
	if err != nil {
		return
	}
	role, err := dboper.GetStaffRole(userID)
	if err != nil {
		respondStaffRoleError(c, err)
		return
	}
	if role != "Administrador" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo un administrador puede consultar alimentadores"})
		return
	}

	feeders, err := dboper.GetFeeders()
	if err != nil {
		log.Printf("get feeders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar los alimentadores"})
		return
	}
	c.JSON(http.StatusOK, feeders)
}

// GetMapReports devuelve las ubicaciones de reportes visibles para el personal autenticado.
func GetMapReports(c *gin.Context) {
	userID, err := authenticateStaff(c)
	if err != nil {
		return
	}
	role, err := dboper.GetStaffRole(userID)
	if err != nil {
		respondStaffRoleError(c, err)
		return
	}
	if role != "Administrador" && role != "Alimentador" {
		c.JSON(http.StatusForbidden, gin.H{"error": "rol sin permisos para consultar el mapa"})
		return
	}

	reports, err := dboper.GetMapReports(userID)
	if err != nil {
		log.Printf("get map reports: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar los puntos del mapa"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

func authenticateStaff(c *gin.Context) (uint32, error) {
	userID, err := token.Authenticate(c.GetHeader("Authorization"))
	if err == nil {
		return userID, nil
	}
	if errors.Is(err, token.ErrInvalidToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorización inválido o faltante"})
		return 0, err
	}
	log.Printf("authenticate staff request: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible validar la autorización"})
	return 0, err
}

func respondStaffRoleError(c *gin.Context, err error) {
	if errors.Is(err, dboper.ErrStaffRoleRequired) {
		c.JSON(http.StatusForbidden, gin.H{"error": "se requiere rol de administrador o alimentador"})
		return
	}
	log.Printf("get staff role: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible validar el rol"})
}

// Institutions maneja la consulta de todas las instituciones disponibles en la base de datos.
// Devuelve la lista de instituciones en formato JSON. Si ocurre un error al consultar la base de datos,
// devuelve un error adecuado.
func Institutions(c *gin.Context) {
	institutions, err := dboper.Institutions()
	if err != nil {
		log.Printf("Error obteniendo reportes: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar las instituciones"})
		return
	}

	c.JSON(http.StatusOK, institutions)
}

// InstitutionByID devuelve el detalle de una institución.
func InstitutionByID(c *gin.Context) {
	id, ok := parseInstitutionID(c)
	if !ok {
		return
	}
	institution, err := dboper.InstitutionByID(id)
	if errors.Is(err, dboper.ErrInstitutionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "institución no encontrada"})
		return
	}
	if err != nil {
		log.Printf("get institution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible consultar la institución"})
		return
	}
	c.JSON(http.StatusOK, institution)
}

// CreateInstitution registra una institución; solo los administradores pueden realizar esta operación.
func CreateInstitution(c *gin.Context) {
	if !requireAdministrator(c) {
		return
	}
	var request institutionCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de institución inválidos"})
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	request.Type = strings.TrimSpace(request.Type)
	if request.Name == "" || request.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nombre y tipo son obligatorios"})
		return
	}
	if err := validateInstitutionFields(request.Name, request.Type, request.Email, request.Phone, request.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.TerritoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_territorio debe ser un entero positivo"})
		return
	}

	id, err := dboper.CreateInstitution(dboper.InstitutionInput{
		TerritoryID: request.TerritoryID,
		Name:        request.Name,
		Type:        request.Type,
		Email:       request.Email,
		Phone:       request.Phone,
		Address:     request.Address,
	})
	switch {
	case errors.Is(err, dboper.ErrDuplicateInstitution):
		c.JSON(http.StatusConflict, gin.H{"error": "ya existe una institución con ese nombre, tipo y territorio"})
	case errors.Is(err, dboper.ErrInstitutionTerritoryNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "el territorio indicado no existe"})
	case err != nil:
		log.Printf("create institution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible registrar la institución"})
	default:
		c.JSON(http.StatusCreated, gin.H{"id_institucion": id, "mensaje": "Institución registrada correctamente"})
	}
}

// UpdateInstitution aplica los campos enviados sin modificar los que se omitan.
func UpdateInstitution(c *gin.Context) {
	if !requireAdministrator(c) {
		return
	}
	id, ok := parseInstitutionID(c)
	if !ok {
		return
	}
	var request institutionUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de institución inválidos"})
		return
	}
	if request.TerritoryID == nil && request.Name == nil && request.Type == nil &&
		request.Email == nil && request.Phone == nil && request.Address == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se requiere al menos un campo para actualizar"})
		return
	}
	if request.Name != nil {
		trimmed := strings.TrimSpace(*request.Name)
		request.Name = &trimmed
	}
	if request.Type != nil {
		trimmed := strings.TrimSpace(*request.Type)
		request.Type = &trimmed
	}
	if request.Name != nil && *request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nombre no puede quedar vacío"})
		return
	}
	if request.Type != nil && *request.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo no puede quedar vacío"})
		return
	}
	if err := validateInstitutionFields(optionalValue(request.Name), optionalValue(request.Type), request.Email, request.Phone, request.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.TerritoryID != nil && *request.TerritoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_territorio debe ser un entero positivo"})
		return
	}

	err := dboper.UpdateInstitution(id, dboper.InstitutionUpdate{
		TerritoryID: request.TerritoryID,
		Name:        request.Name,
		Type:        request.Type,
		Email:       request.Email,
		Phone:       request.Phone,
		Address:     request.Address,
	})
	switch {
	case errors.Is(err, dboper.ErrInstitutionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "institución no encontrada"})
	case errors.Is(err, dboper.ErrDuplicateInstitution):
		c.JSON(http.StatusConflict, gin.H{"error": "ya existe una institución con ese nombre, tipo y territorio"})
	case errors.Is(err, dboper.ErrInstitutionTerritoryNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": "el territorio indicado no existe"})
	case err != nil:
		log.Printf("update institution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible actualizar la institución"})
	default:
		c.JSON(http.StatusOK, gin.H{"mensaje": "Institución actualizada correctamente"})
	}
}

// UpdateInstitutionStatus activa o desactiva una institución sin eliminarla.
func UpdateInstitutionStatus(c *gin.Context) {
	if !requireAdministrator(c) {
		return
	}
	id, ok := parseInstitutionID(c)
	if !ok {
		return
	}
	var request institutionStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "es_activo debe ser true o false"})
		return
	}
	if err := dboper.UpdateInstitutionStatus(id, *request.Active); err != nil {
		if errors.Is(err, dboper.ErrInstitutionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "institución no encontrada"})
			return
		}
		log.Printf("update institution status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible actualizar el estado de la institución"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"es_activo": *request.Active, "mensaje": "Estado de la institución actualizado correctamente"})
}

func parseInstitutionID(c *gin.Context) (uint32, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id de institución inválido"})
		return 0, false
	}
	return uint32(id), true
}

func requireAdministrator(c *gin.Context) bool {
	userID, err := authenticateStaff(c)
	if err != nil {
		return false
	}
	role, err := dboper.GetStaffRole(userID)
	if err != nil {
		respondStaffRoleError(c, err)
		return false
	}
	if role != "Administrador" {
		c.JSON(http.StatusForbidden, gin.H{"error": "solo un administrador puede gestionar instituciones"})
		return false
	}
	return true
}

func validateInstitutionFields(name, institutionType string, email, phone, address *string) error {
	if name != "" && len(name) > 50 {
		return errors.New("nombre debe tener máximo 50 caracteres")
	}
	if institutionType != "" && len(institutionType) > 50 {
		return errors.New("tipo debe tener máximo 50 caracteres")
	}
	if email != nil && strings.TrimSpace(*email) != "" {
		if len(*email) > 254 {
			return errors.New("correo debe tener máximo 254 caracteres")
		}
		parsed, err := mail.ParseAddress(*email)
		if err != nil || parsed.Address != *email {
			return errors.New("correo inválido")
		}
	}
	if phone != nil && len(*phone) > 50 {
		return errors.New("telefono debe tener máximo 50 caracteres")
	}
	if address != nil && len(*address) > 100 {
		return errors.New("direccion debe tener máximo 100 caracteres")
	}
	return nil
}

func optionalValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// UpdateReportStatus maneja la actualización del estado de un reporte específico.
// Valida el token de autorización, el folio y el nuevo estado, y actualiza la información en la base de datos.
// Si el usuario no tiene permiso para actualizar el reporte o si ocurre un error, devuelve un error adecuado.
func UpdateReportStatus(c *gin.Context) {
	userID, err := token.Authenticate(c.GetHeader("Authorization"))
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorización inválido o faltante"})
			return
		}
		log.Printf("authenticate report status request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible validar la autorización"})
		return
	}
	var request statusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estado requerido"})
		return
	}
	if !isReportStatus(request.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estado de expediente inválido"})
		return
	}
	if err := dboper.UpdateReportStatus(userID, c.Param("folio"), request.Status, request.Observations); err != nil {
		if errors.Is(err, dboper.ErrReportAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "no tienes autorización para actualizar este expediente"})
			return
		}
		// Registra el error técnico únicamente en el servidor.
		// No expone detalles de MariaDB al navegador.
		log.Printf(
			"update report status failed: %v",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "no fue posible actualizar el estado del reporte"},
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "Estado del expediente actualizado a: " + request.Status})
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// validateReportRequest valida los campos de un reporte antes de su creación.
func validateReportRequest(request reportRequest) error {
	if strings.TrimSpace(request.Description) == "" {
		return errors.New("descripcion no puede estar vacía")
	}
	if !isImmediateDanger(request.ImmediateDanger) {
		return errors.New("peligro_inmediato debe ser Sí, No o No sé")
	}
	if len(request.Description) > 65535 {
		return errors.New("descripcion excede el tamaño permitido")
	}
	if request.Location.Latitude == nil || *request.Location.Latitude < -90 || *request.Location.Latitude > 90 {
		return errors.New("latitud fuera de rango")
	}
	if request.Location.Longitude == nil || *request.Location.Longitude < -180 || *request.Location.Longitude > 180 {
		return errors.New("longitud fuera de rango")
	}
	if len(request.Location.Address) > 255 {
		return errors.New("direccion_texto excede el tamaño permitido")
	}
	if request.ApproximateTime != "" {
		if _, err := time.Parse("15:04", request.ApproximateTime); err != nil {
			return errors.New("horario_aproximado debe usar el formato HH:MM")
		}
	}
	if len(request.Minors) == 0 || len(request.Minors) > 255 {
		return errors.New("menores debe incluir entre 1 y 255 elementos")
	}
	for index, minor := range request.Minors {
		if minor.Gender != "Masculino" && minor.Gender != "Femenino" && minor.Gender != "No identificado" {
			return fmt.Errorf("genero inválido para el menor %d", index+1)
		}
		if len(minor.PhysicalDetails) > 255 {
			return fmt.Errorf("descripcion_fisica excede el tamaño permitido para el menor %d", index+1)
		}
	}
	if request.ContactEmail != nil && len(*request.ContactEmail) > 254 {
		return errors.New("correo_contacto excede el tamaño permitido")
	}
	if request.ContactEmail != nil && strings.TrimSpace(*request.ContactEmail) != "" {
		address, err := mail.ParseAddress(*request.ContactEmail)
		if err != nil || address.Address != *request.ContactEmail {
			return errors.New("correo_contacto no es válido")
		}
	}
	return nil
}

/*
isReportStatus valida que el estado solicitado corresponda
a uno de los valores oficiales definidos en expediente.estado.

Se utiliza para evitar enviar estados incompatibles a MariaDB.
*/
func isReportStatus(status string) bool {
	switch status {
	case "Registrado",
		"En revisión",
		"En seguimiento",
		"Canalizado",
		"Concluido",
		"Archivado",
		"Cancelado",
		"Reincidente":
		return true
	default:
		return false
	}
}

func isImmediateDanger(value string) bool {
	switch value {
	case "Sí", "No", "No sé":
		return true
	default:
		return false
	}
}

func isHistoryStatus(status string) bool {
	switch status {
	case "Registrado", "En revisión", "En seguimiento", "Canalizado", "Concluido", "Archivado", "Cancelado", "Reincidente":
		return true
	default:
		return false
	}
}

// detectEvidenceType determina el tipo de contenido de un archivo de evidencia basado en su encabezado.
// Devuelve el tipo MIME si es válido, o un error si el tipo no es permitido.
func detectEvidenceType(file io.Reader) (string, error) {
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	contentType := http.DetectContentType(header[:n])
	if n >= 12 && bytes.Equal(header[:4], []byte("RIFF")) && bytes.Equal(header[8:12], []byte("WEBP")) {
		return "image/webp", nil
	}
	if n >= 4 && bytes.Equal(header[:4], []byte{0x1a, 0x45, 0xdf, 0xa3}) {
		return "video/webm", nil
	}
	if n >= 12 && bytes.Equal(header[4:8], []byte("ftyp")) {
		if bytes.Equal(header[8:12], []byte("qt  ")) {
			return "video/quicktime", nil
		}
		return "video/mp4", nil
	}
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return contentType, nil
	case "video/mp4", "video/webm", "video/quicktime":
		return contentType, nil
	default:
		return "", errors.New("solo se permiten imágenes JPEG, PNG o WebP y videos MP4, WebM o QuickTime")
	}
}

// extensionForType devuelve la extensión de archivo correspondiente a un tipo MIME dado.
func extensionForType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "video/quicktime":
		return ".mov"
	default:
		return ""
	}
}

// newFolio genera un folio único para un reporte, compuesto por el prefijo "REP-",
// el año actual y una cadena aleatoria de seis caracteres.
func newFolio() (string, error) {
	randomPart, err := randomString(6, folioAlphabet)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("REP-%d-%s", time.Now().Year(), randomPart), nil
}

// randomString genera una cadena aleatoria de una longitud específica utilizando un alfabeto dado.
func randomString(length int, alphabet string) (string, error) {
	var result strings.Builder
	result.Grow(length)
	max := big.NewInt(int64(len(alphabet)))
	for range length {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result.WriteByte(alphabet[index.Int64()])
	}
	return result.String(), nil
}
