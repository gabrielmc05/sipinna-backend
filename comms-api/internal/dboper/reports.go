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
	Latitude        float64
	Longitude       float64
	Address         string
	TerritoryID     uint32
	Minors          []Minor
	Tags            []Tag
}

type Evidence struct {
	Type string `json:"tipo_archivo"`
	Path string `json:"ruta"`
}

type ReportSummary struct {
	Folio       string    `json:"folio"`
	Status      string    `json:"estatus"`
	SubmittedAt timeValue `json:"fecha_envio"`
	Description string    `json:"descripcion"`
	Address     string    `json:"direccion"`
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
		"CALL CrearReporte(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		folio,
		report.UserID,
		report.ContactEmail,
		codeHash,
		report.Description,
		report.ApproximateTime,
		report.Anonymous,
		report.WantsFollowUp,
		report.AdultInvolved,
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
		db.Exec("CALL AsignarEtiquetaNuevoReporte(?,?,?)", folio, tag)
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

func ConsultFolio(folio string, userID uint32) (ReportSummary, error) {
	var report ReportSummary
	err := db.QueryRow("CALL ConsultarReporte(?,?)", folio, userID).Scan(
		&report.Folio,
		&report.Status,
		&report.SubmittedAt,
		&report.Description,
		&report.Address,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReportSummary{}, sql.ErrNoRows
		}
		return ReportSummary{}, fmt.Errorf("consult report: %w", err)
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
	err := db.QueryRow("SELECT llave_firma FROM token WHERE id_usuario=?", userID).Scan(&key)
	if err != nil {
		return nil, fmt.Errorf("get user signing key: %w", err)
	}
	return key, nil
}
