package dboper

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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
