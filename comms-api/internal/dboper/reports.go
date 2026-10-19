package dboper

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrReportNotFound     = sql.ErrNoRows
	ErrReportAccessDenied = errors.New("user cannot update this report")
	ErrStaffRoleRequired  = errors.New("user does not have a staff role")
	ErrAdminRoleRequired  = errors.New("user is not an administrator")
	ErrFeederRoleRequired = errors.New("target user is not an alimentador")
)

type Minor struct {
	ApproximateAge  *uint8 `json:"edad_aproximada"`
	Gender          string `json:"genero"`
	PhysicalDetails string `json:"descripcion_fisica"`
}

type Tag struct {
	Name string `json:"nombre_etiqueta" binding:"required,max=50"`
	Type string `json:"tipo" binding:"required"`
}

type ReportCreation struct {
	UserID          *uint32
	ContactEmail    *string
	Description     string
	ApproximateTime *string
	Anonymous       bool
	WantsFollowUp   bool
	AdultInvolved   bool
	ImmediateDanger string
	Latitude        float64
	Longitude       float64
	Address         string
	TerritoryID     uint32
	Minors          []Minor
	Tags            []Tag
}

type Evidence struct {
	ID   uint32 `json:"id_evidencia"`
	Type string `json:"tipo_archivo"`
	Path string `json:"ruta"`
}

type HistoricalReport struct {
	Folio            string    `json:"folio"`
	Status           string    `json:"estado"`
	SubmittedAt      timeValue `json:"fecha_envio"`
	Description      string    `json:"descripcion"`
	Address          string    `json:"direccion"`
	TerritoryID      uint32    `json:"id_territorio"`
	AssignedFeederID *uint32   `json:"id_alimentador"`
	ImmediateDanger  string    `json:"peligro_inmediato"`
	Tags             []Tag     `json:"etiquetas"`
}

type ReportSummary struct {
	Folio           string    `json:"folio"`
	Status          string    `json:"estatus"`
	SubmittedAt     timeValue `json:"fecha_envio"`
	Description     string    `json:"descripcion"`
	Address         string    `json:"direccion"`
	ImmediateDanger string    `json:"peligro_inmediato"`
}

// timeValue marshals database timestamps using the ISO-8601 format expected by clients.
type timeValue struct {
	sql.NullTime
}

func (value timeValue) MarshalJSON() ([]byte, error) {
	if !value.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(value.Time.UTC().Format("2006-01-02T15:04:05Z"))
}

// CreateReport es para crear un nuevo reporte en la base de datos.
// Toma un objeto ReportCreation, un folio único y un hash de código como parámetros.
func CreateReport(report ReportCreation, folio, codeHash string) error {
	minorsJSON, err := json.Marshal(report.Minors)
	if err != nil {
		return fmt.Errorf("encode report minors: %w", err)
	}
	_, err = db.Exec(
		"CALL CrearReporte(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		folio,
		report.UserID,
		report.ContactEmail,
		codeHash,
		report.Description,
		report.ApproximateTime,
		report.Anonymous,
		report.WantsFollowUp,
		report.AdultInvolved,
		report.ImmediateDanger,
		report.Latitude,
		report.Longitude,
		report.Address,
		report.TerritoryID,
		string(minorsJSON),
	)
	if err != nil {
		return fmt.Errorf("create report: %w", err)
	}

	for _, tag := range report.Tags {
		db.Exec("CALL AsignarEtiquetaNuevoReporte(?,?,?)", folio, tag.Name, tag.Type)
	}

	return nil
}

func IsDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func CreateEvidences(folio string, evidence []Evidence) error {
	evidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("encode evidence list: %w", err)
	}
	_, err = db.Exec("CALL AgregarEvidencias(?,?)", folio, string(evidenceJSON))
	if err != nil {
		return fmt.Errorf("create evidences: %w", err)
	}
	return nil
}

// ConsultFolio obtiene la información de un reporte
// para un administrador o alimentador autenticado.
//
// El procedimiento almacenado verifica que:
//   - Los administradores puedan consultar todos los reportes.
//   - Los alimentadores únicamente consulten los asignados.
//
// No utiliza el código privado del reportante.
func ConsultFolio(folio string, userID uint32) (ReportSummary, error) {
	var report ReportSummary

	err := db.QueryRow(
		"CALL ConsultarReportePersonal(?,?)",
		folio,
		userID,
	).Scan(
		&report.Folio,
		&report.Status,
		&report.SubmittedAt,
		&report.Description,
		&report.Address,
		&report.ImmediateDanger,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReportSummary{}, sql.ErrNoRows
		}

		return ReportSummary{}, fmt.Errorf(
			"consult staff report: %w",
			err,
		)
	}

	return report, nil
}

func ConsultReport(folio, codeHash string) (ReportSummary, error) {
	var report ReportSummary
	err := db.QueryRow("CALL ConsultarReporte(?,?)", folio, codeHash).Scan(
		&report.Folio,
		&report.Status,
		&report.SubmittedAt,
		&report.Description,
		&report.Address,
		&report.ImmediateDanger,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReportSummary{}, sql.ErrNoRows
		}
		return ReportSummary{}, fmt.Errorf("consult report: %w", err)
	}
	return report, nil
}

func ListUserReports(userID uint32) ([]ReportSummary, error) {
	rows, err := db.Query("CALL ListarReportesUsuario(?)", userID)
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()

	reports := make([]ReportSummary, 0)
	for rows.Next() {
		var report ReportSummary
		if err := rows.Scan(
			&report.Folio,
			&report.Status,
			&report.SubmittedAt,
			&report.Description,
			&report.Address,
			&report.ImmediateDanger,
		); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reports: %w", err)
	}
	return reports, nil
}

func GetStaffRole(userID uint32) (string, error) {
	var role string
	err := db.QueryRow(
		"SELECT rol FROM administrador WHERE id_usuario=?",
		userID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrStaffRoleRequired
	}
	if err != nil {
		return "", fmt.Errorf("get staff role: %w", err)
	}
	return role, nil
}

func GetReportHistory(userID uint32, status *string, territoryID *uint32, dateOrder string) ([]HistoricalReport, error) {
	rows, err := db.Query("CALL ObtenerReportes(?,?,?,?)", userID, status, territoryID, dateOrder)
	if err != nil {
		return nil, fmt.Errorf("get report history: %w", err)
	}
	defer rows.Close()

	reports := make([]HistoricalReport, 0)
	for rows.Next() {
		var report HistoricalReport
		var assignedFeederID sql.NullInt64
		if err := rows.Scan(
			&report.Folio,
			&report.Status,
			&report.SubmittedAt,
			&report.Description,
			&report.Address,
			&report.TerritoryID,
			&assignedFeederID,
			&report.ImmediateDanger,
		); err != nil {
			return nil, fmt.Errorf("scan report history row: %w", err)
		}
		if assignedFeederID.Valid {
			id := uint32(assignedFeederID.Int64)
			report.AssignedFeederID = &id
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate report history: %w", err)
	}
	// Cerramos el resultado del procedimiento antes de realizar otra consulta.
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close report history: %w", err)
	}
	if len(reports) == 0 {
		return reports, nil
	}

	// Solo consultamos etiquetas de folios ya autorizados por ObtenerReportes.
	// La consulta es única para evitar una petición por cada expediente.
	placeholders := make([]string, len(reports))
	arguments := make([]any, len(reports))
	indices := make(map[string]int, len(reports))
	for index := range reports {
		reports[index].Tags = make([]Tag, 0)
		placeholders[index] = "?"
		arguments[index] = reports[index].Folio
		indices[reports[index].Folio] = index
	}
	tagRows, err := db.Query(
		"SELECT folio, nombre_etiqueta, tipo FROM reporteetiqueta WHERE folio IN ("+
			strings.Join(placeholders, ",")+") ORDER BY tipo, nombre_etiqueta",
		arguments...,
	)
	if err != nil {
		return nil, fmt.Errorf("get report history tags: %w", err)
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var folio string
		var tag Tag
		if err := tagRows.Scan(&folio, &tag.Name, &tag.Type); err != nil {
			return nil, fmt.Errorf("scan report history tag: %w", err)
		}
		if index, ok := indices[folio]; ok {
			reports[index].Tags = append(reports[index].Tags, tag)
		}
	}
	if err := tagRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate report history tags: %w", err)
	}
	return reports, nil
}

func AssignReport(actorID uint32, folio string, feederID uint32) error {
	var reportFound, actorIsAdmin, targetIsFeeder bool
	err := db.QueryRow(
		"CALL AsignarReporte(?,?,?)",
		actorID,
		folio,
		feederID,
	).Scan(&reportFound, &actorIsAdmin, &targetIsFeeder)
	if err != nil {
		return fmt.Errorf("assign report: %w", err)
	}
	if !reportFound {
		return ErrReportNotFound
	}
	if !actorIsAdmin {
		return ErrAdminRoleRequired
	}
	if !targetIsFeeder {
		return ErrFeederRoleRequired
	}
	return nil
}

func GetReportEvidences(userID uint32, folio string) ([]Evidence, error) {
	rows, err := db.Query("CALL ObtenerEvidencias(?,?)", folio, userID)
	if err != nil {
		return nil, fmt.Errorf("get report evidences: %w", err)
	}
	defer rows.Close()

	evidences := make([]Evidence, 0)
	var reportFound, userAuthorized bool
	for rows.Next() {
		var evidenceID sql.NullInt64
		var fileType, path sql.NullString
		if err := rows.Scan(&reportFound, &userAuthorized, &evidenceID, &fileType, &path); err != nil {
			return nil, fmt.Errorf("scan report evidence row: %w", err)
		}
		if evidenceID.Valid {
			evidences = append(evidences, Evidence{
				ID:   uint32(evidenceID.Int64),
				Type: fileType.String,
				Path: path.String,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate report evidences: %w", err)
	}
	if !reportFound {
		return nil, ErrReportNotFound
	}
	if !userAuthorized {
		return nil, ErrReportAccessDenied
	}
	return evidences, nil
}

func UpdateReportStatus(userID uint32, folio, status string, observations *string) error {
	var found, allowed bool
	err := db.QueryRow("CALL ActualizarEstadoReporte(?,?,?,?)", userID, folio, status, observations).Scan(&found, &allowed)
	if err != nil {
		return fmt.Errorf("update report status: %w", err)
	}
	if !found {
		return ErrReportNotFound
	}
	if !allowed {
		return ErrReportAccessDenied
	}
	return nil
}

func GetUserSignatureKey(userID uint32) ([]byte, error) {
	var key []byte
	err := db.QueryRow("SELECT t.llave_firma FROM token t INNER JOIN usuario u ON u.id_usuario = t.id_usuario WHERE t.id_usuario = ? AND u.es_activo = 1",
		userID).Scan(&key)

	if err != nil {
		return nil, fmt.Errorf("get user signing key: %w", err)
	}
	return key, nil
}

// EtiquetaCatalogo representa una etiqueta del catálogo
// institucional de MariaDB.
type EtiquetaCatalogo struct {
	Name   string `json:"nombre_etiqueta"`
	Type   string `json:"tipo"`
	Active bool   `json:"es_activo"`
}

// ListarCatalogoEtiquetas devuelve el catálogo completo,
// incluidas las etiquetas que aún no tienen reportes.
func ListarCatalogoEtiquetas() ([]EtiquetaCatalogo, error) {
	rows, err := db.Query(
		`SELECT nombre_etiqueta, tipo, es_activo
		 FROM etiqueta
		 ORDER BY tipo, nombre_etiqueta`,
	)
	if err != nil {
		return nil, fmt.Errorf("consultar etiquetas: %w", err)
	}
	defer rows.Close()

	resultado := make([]EtiquetaCatalogo, 0)

	for rows.Next() {
		var etiqueta EtiquetaCatalogo
		if err := rows.Scan(
			&etiqueta.Name,
			&etiqueta.Type,
			&etiqueta.Active,
		); err != nil {
			return nil, fmt.Errorf("leer etiqueta: %w", err)
		}
		resultado = append(resultado, etiqueta)
	}

	return resultado, rows.Err()
}

// CrearEtiquetaCatalogo registra una etiqueta nueva.
// La restricción de clave primaria evita duplicados.
func CrearEtiquetaCatalogo(nombre, tipo string) error {
	_, err := db.Exec(
		`INSERT INTO etiqueta
		 (nombre_etiqueta, tipo, es_activo)
		 VALUES (?, ?, 1)`,
		nombre,
		tipo,
	)
	if err != nil {
		return fmt.Errorf("crear etiqueta: %w", err)
	}
	return nil
}

// RenombrarEtiquetaCatalogo cambia el nombre de una etiqueta
// manteniendo las asociaciones de los reportes.
//
// Se inserta primero la nueva etiqueta porque la clave
// foránea de reporteetiqueta no tiene ON UPDATE CASCADE.
func RenombrarEtiquetaCatalogo(
	tipo, nombreAnterior, nombreNuevo string,
) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar edición: %w", err)
	}
	defer tx.Rollback()

	var activa bool
	err = tx.QueryRow(
		`SELECT es_activo
		 FROM etiqueta
		 WHERE tipo = ? AND nombre_etiqueta = ?
		 FOR UPDATE`,
		tipo,
		nombreAnterior,
	).Scan(&activa)

	if err != nil {
		return fmt.Errorf("buscar etiqueta: %w", err)
	}

	if nombreAnterior == nombreNuevo {
		return nil
	}

	_, err = tx.Exec(
		`INSERT INTO etiqueta
		 (nombre_etiqueta, tipo, es_activo)
		 VALUES (?, ?, ?)`,
		nombreNuevo,
		tipo,
		activa,
	)
	if err != nil {
		return fmt.Errorf("registrar nuevo nombre: %w", err)
	}

	_, err = tx.Exec(
		`UPDATE reporteetiqueta
		 SET nombre_etiqueta = ?
		 WHERE nombre_etiqueta = ? AND tipo = ?`,
		nombreNuevo,
		nombreAnterior,
		tipo,
	)
	if err != nil {
		return fmt.Errorf("actualizar asociaciones: %w", err)
	}

	_, err = tx.Exec(
		`DELETE FROM etiqueta
		 WHERE nombre_etiqueta = ? AND tipo = ?`,
		nombreAnterior,
		tipo,
	)
	if err != nil {
		return fmt.Errorf("retirar nombre anterior: %w", err)
	}

	return tx.Commit()
}

// AplicarEtiquetaReportes aplica una etiqueta existente
// a varios expedientes en una transacción.
//
// Las etiquetas de prioridad son excluyentes: antes de
// aplicar una nueva prioridad, sustituye la anterior.
//
// Categoría y Riesgo admiten varias etiquetas.
func AplicarEtiquetaReportes(
	folios []string,
	nombre, tipo string,
) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar aplicación: %w", err)
	}
	defer tx.Rollback()

	var activa bool
	err = tx.QueryRow(
		`SELECT es_activo
		 FROM etiqueta
		 WHERE nombre_etiqueta = ? AND tipo = ?`,
		nombre,
		tipo,
	).Scan(&activa)

	if err != nil {
		return fmt.Errorf("consultar etiqueta: %w", err)
	}
	if !activa {
		return fmt.Errorf("la etiqueta está inactiva")
	}

	// Comprobar todos los folios antes de modificar datos.
	for _, folio := range folios {
		var existente string
		err = tx.QueryRow(
			`SELECT folio FROM reporte
			 WHERE folio = ? FOR UPDATE`,
			folio,
		).Scan(&existente)
		if err != nil {
			return fmt.Errorf("folio %s: %w", folio, err)
		}
	}

	for _, folio := range folios {
		if tipo == "Prioridad" {
			_, err = tx.Exec(
				`DELETE FROM reporteetiqueta
				 WHERE folio = ? AND tipo = 'Prioridad'`,
				folio,
			)
			if err != nil {
				return fmt.Errorf("reemplazar prioridad: %w", err)
			}
		}

		_, err = tx.Exec(
			`INSERT IGNORE INTO reporteetiqueta
			 (folio, nombre_etiqueta, tipo)
			 VALUES (?, ?, ?)`,
			folio,
			nombre,
			tipo,
		)
		if err != nil {
			return fmt.Errorf("asignar etiqueta: %w", err)
		}
	}

	return tx.Commit()
}
