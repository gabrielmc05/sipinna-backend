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
		if errors.Is(err, dboper.ErrReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible actualizar el estado del reporte"})
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

// isReportStatus verifica si el estado proporcionado es uno de los estados válidos para un reporte.
func isReportStatus(status string) bool {
	switch status {
	case "Nuevo", "En revisión", "Asignado", "En investigación", "Canalizado a Procuraduría", "Cerrado", "Descartado":
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
