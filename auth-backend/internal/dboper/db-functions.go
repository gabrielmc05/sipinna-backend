package dboper

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func Probando(user_id uint32, newPassword string) {
	var (
		err          error
		result       sql.Result
		lastInsertID int64
		rowsAffected int64
	)
	newPassword, err = EncryptHashUserPassword(newPassword)
	if err != nil {
		log.Println(err)
	}

	if result, err = db.Exec("UPDATE usuario SET hash_contrasenna=UNHEX(?) where id_usuario=?",
		newPassword, user_id); err != nil {
		log.Printf("Error in query: %s", err)
		return
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return 0, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return 0, err
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)
}

func GetUserPasswordHash(user_id uint32) (password string, err error) {
	var (
		hash_contrasenna []byte
	)

	if err = db.QueryRow("SELECT hash_contrasenna FROM usuario where id_usuario=?",
		user_id).Scan(&hash_contrasenna); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return "", err
	}

	if password, err = DecryptUserPasswordHash(hash_contrasenna); err != nil {
		log.Fatal(err)
	}

	return password, nil
}

func GetUserSignatureKey(user_id uint32) (key []byte, err error) {
	if err = db.QueryRow("SELECT llave_firma FROM usuario where id_usuario=?",
		user_id).Scan(&key); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return nil, err
	}

	return key, nil
}

func GetUserIDWithEmail(email string) (user_id uint32, err error) {
	if err = db.QueryRow("SELECT id_usuario FROM usuario WHERE correo=?",
		email).Scan(&user_id); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return 0, err
	}

	return user_id, nil
}

func NewUser(email, password, name, surname string, phoneNumber ...string) (user_id uint32, err error) {
	var (
		encryptedHash string
		result        sql.Result
		lastInsertID  int64
		rowsAffected  int64
	)

	if encryptedHash, err = EncryptHashUserPassword(password); err != nil {
		log.Printf("Error at creating encrypting user's password hash: %s", err)
		return 0, err
	}

	if len(phoneNumber) > 0 {
		result, err = db.Exec("INSERT INTO usuario(correo,hash_contrasenna,nombre,apellidos,telefono)VALUES(?,UNHEX(?),?,?,?)",
			email, encryptedHash, name, surname, phoneNumber[0])
	} else {
		result, err = db.Exec("INSERT INTO usuario(correo,hash_contrasenna,nombre,apellidos)VALUES(?,UNHEX(?),?,?)",
			email, encryptedHash, name, surname)
	}

	if err != nil {
		err = fmt.Errorf("Error: %s, inserting new user.", err)
		log.Print(err)
		return 0, err
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return 0, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		err = fmt.Errorf("Error: in query: %s", err)
		log.Print(err)
		return 0, err
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)

	return uint32(lastInsertID), nil
}
