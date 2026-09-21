package dboper

import (
	"database/sql"
	"log"
)

func Probando(db *sql.DB) {
	const query string = `
SELECT id_usuario, correo, contrasenna_encriptada FROM 'usuario';
`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id_usuario             int32
			correo                 string
			contrasenna_encriptada string
		)
		if err := rows.Scan(&id_usuario, &correo, &contrasenna_encriptada); err != nil {
			log.Fatal(err)
		}
		log.Printf("id_usuario: %d, correo: %s, contraseña encriptada: %s\n", id_usuario, correo, contrasenna_encriptada)
	}
	if !rows.NextResultSet() {
		log.Fatalf("expected more result sets: %v", rows.Err())
	}
}
