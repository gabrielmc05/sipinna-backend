package dboper

import (
	"log"
	"strconv"
	"strings"
)

func Probando() {
	const query string = `
SELECT id_usuario, correo, contrasenna_encriptada FROM usuario;
`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error in query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id_usuario  uint32
			correo      string
			contrasenna string
		)
		if err := rows.Scan(&id_usuario, &correo, &contrasenna); err != nil {
			log.Fatal(rows.Err())
		}
		log.Printf("id_usuario: %d, correo: %s, contraseña encriptada: %s\n", id_usuario, correo, contrasenna)
	}
}

func GetUserPassword(user_id uint32) string {
	var (
		query                  strings.Builder
		contrasenna_encriptada string
	)
	query.WriteString("SELECT contrasenna_encriptada FROM usuario where id_usuario=")
	query.WriteString(strconv.FormatUint(uint64(user_id), 10))
	query.WriteString(";")

	err := db.QueryRow(query.String()).Scan(&contrasenna_encriptada)
	if err != nil {
		log.Fatalf("Error: %s, in query: %s", err, query.String())
	}

	// TODO: desencriptar contraseña
	return contrasenna_encriptada
}
