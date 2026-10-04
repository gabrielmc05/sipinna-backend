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

	if result, err = db.Exec(query.String()); err != nil {
		log.Printf("Error in query: %s", err)
		return
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		log.Printf("Error in query: %s", err)
		return
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		log.Printf("Error in query: %s", err)
		return
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

	if err = db.QueryRow(query.String()).Scan(&hash_contrasenna); err != nil {
		log.Printf("Error: %s, in query: %s", err, query.String())
		return "", err
	}

	if password, err = DecryptUserPasswordHash(hash_contrasenna); err != nil {
		log.Fatal(err)
	}

	return password, nil
}

func GetUserSignatureKey(user_id uint32) (key []byte, err error) {
	var query strings.Builder

	query.WriteString("SELECT llave_firma FROM usuario where id_usuario=")
	query.WriteString(strconv.FormatUint(uint64(user_id), 10))
	query.WriteByte(';')

	if err = db.QueryRow(query.String()).Scan(&key); err != nil {
		log.Printf("Error: %s, in query: %s", err, query.String())
		return nil, err
	}

	return key, nil
}

func GetUserIDWithEmail(email string) (user_id uint32, err error) {
	var query strings.Builder

	query.WriteString("SELECT id_usuario FROM usuario WHERE correo='")
	query.WriteString(email)
	query.WriteString("';")

	if err = db.QueryRow(query.String()).Scan(&user_id); err != nil {
		log.Printf("Error: %s, in query: %s", err, query.String())
		return 0, err
	}

	return user_id, nil
}

//func NewUser(email, password, name, surname, phoneNumber string) (err error) {
//
//}
