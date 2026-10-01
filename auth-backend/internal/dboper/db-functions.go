package dboper

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

func Probando(user_id uint32, newPassword string) {
	var (
		query        strings.Builder
		err          error
		result       sql.Result
		lastInsertID int64
		rowsAffected int64
	)
	newPassword, err = EncryptHashUserPassword(newPassword)
	if err != nil {
		log.Println(err)
	}

	query.WriteString("UPDATE usuario SET hash_contrasenna = UNHEX('")
	query.WriteString(newPassword)
	query.WriteString("') where id_usuario=")
	query.WriteString(strconv.FormatUint(uint64(user_id), 10))
	query.WriteByte(';')

	result, err = db.Exec(query.String())
	if err != nil {
		log.Printf("Error in query: %s", err)
	}

	lastInsertID, err = result.LastInsertId()
	if err != nil {
		log.Printf("Error in query: %s", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Printf("Error in query: %s", err)
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)
}

func GetUserPasswordHash(user_id uint32) (password string, err error) {
	var (
		query            strings.Builder
		hash_contrasenna []byte
	)
	query.WriteString("SELECT hash_contrasenna FROM usuario where id_usuario=")
	query.WriteString(strconv.FormatUint(uint64(user_id), 10))
	query.WriteByte(';')

	err = db.QueryRow(query.String()).Scan(&hash_contrasenna)
	if err != nil {
		log.Printf("Error: %s, in query: %s", err, query.String())
		return "", err
	}

	password, err = DecryptUserPasswordHash(hash_contrasenna)
	if err != nil {
		log.Fatal(err)
	}

	return password, nil
}
