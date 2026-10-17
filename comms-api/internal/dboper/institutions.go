package dboper

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrInstitutionNotFound          = sql.ErrNoRows
	ErrDuplicateInstitution         = errors.New("institution already exists")
	ErrInstitutionTerritoryNotFound = errors.New("institution territory not found")
)

type Institution struct {
	ID        uint32  `json:"id_institucion"`
	Name      string  `json:"nombre"`
	Type      string  `json:"tipo"`
	Email     *string `json:"correo"`
	Phone     *string `json:"telefono"`
	Address   *string `json:"direccion"`
	IsActive  bool    `json:"es_activo"`
	Territory struct {
		ID    uint32 `json:"id_territorio"`
		State string `json:"estado"`
		Name  string `json:"nombre"`
	} `json:"territorio"`
}

type InstitutionInput struct {
	TerritoryID uint32
	Name        string
	Type        string
	Email       *string
	Phone       *string
	Address     *string
}

type InstitutionUpdate struct {
	TerritoryID *uint32
	Name        *string
	Type        *string
	Email       *string
	Phone       *string
	Address     *string
}

func Institutions() ([]Institution, error) {
	rows, err := db.Query("CALL ObtenerInstituciones()")
	if err != nil {
		return nil, fmt.Errorf("get institutions: %w", err)
	}
	defer rows.Close()

	institutions := make([]Institution, 0)
	for rows.Next() {
		institution, err := scanInstitution(rows)
		if err != nil {
			return nil, fmt.Errorf("scan institution: %w", err)
		}
		institutions = append(institutions, institution)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate institutions: %w", err)
	}
	return institutions, nil
}

func InstitutionByID(id uint32) (Institution, error) {
	rows, err := db.Query("CALL ObtenerInstitucionPorId(?)", id)
	if err != nil {
		return Institution{}, fmt.Errorf("get institution: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return Institution{}, fmt.Errorf("read institution: %w", err)
		}
		return Institution{}, ErrInstitutionNotFound
	}
	institution, err := scanInstitution(rows)
	if err != nil {
		return Institution{}, fmt.Errorf("scan institution: %w", err)
	}
	return institution, nil
}

func CreateInstitution(input InstitutionInput) (uint32, error) {
	var id uint32
	err := db.QueryRow(
		"CALL RegistrarInstitucion(?,?,?,?,?,?)",
		input.TerritoryID,
		input.Name,
		input.Type,
		input.Email,
		input.Phone,
		input.Address,
	).Scan(&id)
	if err != nil {
		if isInstitutionDuplicateError(err) {
			return 0, ErrDuplicateInstitution
		}
		if isInstitutionTerritoryError(err) {
			return 0, ErrInstitutionTerritoryNotFound
		}
		return 0, fmt.Errorf("create institution: %w", err)
	}
	return id, nil
}

func UpdateInstitution(id uint32, update InstitutionUpdate) error {
	var found, duplicate, territoryFound bool
	err := db.QueryRow(
		"CALL ActualizarInstitucion(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		id,
		update.TerritoryID,
		update.Name,
		update.Type,
		update.Email,
		update.Phone,
		update.Address,
		update.TerritoryID != nil,
		update.Name != nil,
		update.Type != nil,
		update.Email != nil,
		update.Phone != nil,
		update.Address != nil,
	).Scan(&found, &duplicate, &territoryFound)
	if err != nil {
		return fmt.Errorf("update institution: %w", err)
	}
	if !found {
		return ErrInstitutionNotFound
	}
	if duplicate {
		return ErrDuplicateInstitution
	}
	if update.TerritoryID != nil && !territoryFound {
		return ErrInstitutionTerritoryNotFound
	}
	return nil
}

func UpdateInstitutionStatus(id uint32, active bool) error {
	var found bool
	err := db.QueryRow("CALL ActualizarEstadoInstitucion(?,?)", id, active).Scan(&found)
	if err != nil {
		return fmt.Errorf("update institution status: %w", err)
	}
	if !found {
		return ErrInstitutionNotFound
	}
	return nil
}

type institutionScanner interface {
	Scan(dest ...any) error
}

func scanInstitution(scanner institutionScanner) (Institution, error) {
	var institution Institution
	var state sql.NullString
	err := scanner.Scan(
		&institution.ID,
		&institution.Name,
		&institution.Type,
		&institution.Email,
		&institution.Phone,
		&institution.Address,
		&institution.IsActive,
		&institution.Territory.ID,
		&state,
		&institution.Territory.Name,
	)
	if err != nil {
		return Institution{}, err
	}
	if state.Valid {
		institution.Territory.State = state.String
	}
	return institution, nil
}

func isInstitutionDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && (mysqlErr.Number == 1062 || mysqlErr.Message == "INSTITUTION_DUPLICATE")
}

func isInstitutionTerritoryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Message == "INSTITUTION_TERRITORY_NOT_FOUND"
}
