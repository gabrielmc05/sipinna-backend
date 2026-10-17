package dboper

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrTerritoryNotFound = sql.ErrNoRows

type Territory struct {
	ID           uint32 `json:"id_territorio"`
	State        string `json:"estado"`
	Municipality string `json:"nombre"`
}

type TerritoryOptions struct {
	States      []string    `json:"estados"`
	Territories []Territory `json:"territorios"`
}

type Feeder struct {
	ID            uint32  `json:"id_usuario"`
	Name          string  `json:"nombre"`
	Surname       string  `json:"apellidos"`
	Email         string  `json:"correo"`
	Phone         *string `json:"telefono"`
	InstitutionID uint32  `json:"id_institucion"`
	Institution   string  `json:"institucion"`
}

type MapReport struct {
	Folio           string    `json:"folio"`
	Latitude        float64   `json:"latitud"`
	Longitude       float64   `json:"longitud"`
	Status          string    `json:"estatus"`
	SubmittedAt     timeValue `json:"fecha_envio"`
	TerritoryID     uint32    `json:"id_territorio"`
	State           string    `json:"estado"`
	Municipality    string    `json:"municipio"`
	ImmediateDanger string    `json:"peligro_inmediato"`
}

func GetTerritoryOptions(state, municipality *string) (TerritoryOptions, error) {
	options := TerritoryOptions{
		States:      make([]string, 0),
		Territories: make([]Territory, 0),
	}

	stateRows, err := db.Query("CALL ObtenerEstados()")
	if err != nil {
		return TerritoryOptions{}, fmt.Errorf("get states: %w", err)
	}
	for stateRows.Next() {
		var value string
		if err := stateRows.Scan(&value); err != nil {
			stateRows.Close()
			return TerritoryOptions{}, fmt.Errorf("scan state: %w", err)
		}
		options.States = append(options.States, value)
	}
	if err := stateRows.Err(); err != nil {
		stateRows.Close()
		return TerritoryOptions{}, fmt.Errorf("iterate states: %w", err)
	}
	if err := stateRows.Close(); err != nil {
		return TerritoryOptions{}, fmt.Errorf("close states result: %w", err)
	}

	rows, err := db.Query("CALL ObtenerTerritorios(?,?)", state, municipality)
	if err != nil {
		return TerritoryOptions{}, fmt.Errorf("get territories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var territory Territory
		if err := rows.Scan(&territory.ID, &territory.State, &territory.Municipality); err != nil {
			return TerritoryOptions{}, fmt.Errorf("scan territory: %w", err)
		}
		options.Territories = append(options.Territories, territory)
	}
	if err := rows.Err(); err != nil {
		return TerritoryOptions{}, fmt.Errorf("iterate territories: %w", err)
	}
	return options, nil
}

func GetTerritory(id uint32) (Territory, error) {
	var territory Territory
	err := db.QueryRow("CALL ObtenerTerritorioId(?)", id).Scan(
		&territory.ID,
		&territory.State,
		&territory.Municipality,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Territory{}, ErrTerritoryNotFound
	}
	if err != nil {
		return Territory{}, fmt.Errorf("get territory: %w", err)
	}
	return territory, nil
}

func GetFeeders() ([]Feeder, error) {
	rows, err := db.Query("CALL ObtenerAlimentadores()")
	if err != nil {
		return nil, fmt.Errorf("get feeders: %w", err)
	}
	defer rows.Close()

	feeders := make([]Feeder, 0)
	for rows.Next() {
		var feeder Feeder
		if err := rows.Scan(
			&feeder.ID,
			&feeder.Name,
			&feeder.Surname,
			&feeder.Email,
			&feeder.Phone,
			&feeder.InstitutionID,
			&feeder.Institution,
		); err != nil {
			return nil, fmt.Errorf("scan feeder: %w", err)
		}
		feeders = append(feeders, feeder)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate feeders: %w", err)
	}
	return feeders, nil
}

func GetMapReports(userID uint32) ([]MapReport, error) {
	rows, err := db.Query("CALL ObtenerPuntosMapa(?)", userID)
	if err != nil {
		return nil, fmt.Errorf("get map reports: %w", err)
	}
	defer rows.Close()

	reports := make([]MapReport, 0)
	for rows.Next() {
		var report MapReport
		if err := rows.Scan(
			&report.Folio,
			&report.Latitude,
			&report.Longitude,
			&report.Status,
			&report.SubmittedAt,
			&report.TerritoryID,
			&report.State,
			&report.Municipality,
			&report.ImmediateDanger,
		); err != nil {
			return nil, fmt.Errorf("scan map report: %w", err)
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate map reports: %w", err)
	}
	return reports, nil
}
