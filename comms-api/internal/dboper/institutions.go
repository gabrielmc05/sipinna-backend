// institutions.go es para definir la estructura de datos y las funciones relacionadas con las instituciones en la base de datos.

package dboper

import (
	"database/sql"
	//"encoding/json"
	//"errors"
	"fmt"
)

//var (
//// ErrReportNotFound     = sql.ErrNoRows
//// ErrReportAccessDenied = errors.New("user cannot update this report")
//)

//type Territory struct {
//
//}

type Institution struct {
	Id_institucion uint32 `json:"id_institucion" binding:"number"`
	Nombre         string `json:"nombre"`
	Tipo           string `json:"tipo"`
	Correo         string `json:"correo" binding:"email"`
	Telefono       string `json:"telefono" binding:"e164"`
	Territorio     struct {
		Id_Territorio uint32 `json:"id_territorio" binding:"number"`
		Estado        string `json:"estado"`
		Nombre        string `json:"nombre"`
	} `json:"territorio"`
}


func Institutions() (institutions []Institution, err error) {
	var rows *sql.Rows
	if rows, err = db.Query(
		`SELECT i.id_institucion, i.nombre, i.tipo, i.correo, i.telefono, t.id_territorio, t.estado, t.nombre
		FROM institucion i JOIN territorio t`); err != nil {
		return nil, fmt.Errorf("Error retrieving institutions: %w", err)
	}

	for rows.Next() {
		var institution Institution

		if err = rows.Scan(
			&institution.Id_institucion,
			&institution.Nombre,
			&institution.Tipo,
			&institution.Correo,
			&institution.Telefono,
			&institution.Territorio.Id_Territorio,
			&institution.Territorio.Estado,
			&institution.Territorio.Nombre); err != nil {
			return nil, fmt.Errorf("Error scanning query results: %w", err)
		}

		institutions = append(institutions, institution)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return institutions, nil
}
