package dboper

import (
	"database/sql"
	"fmt"
	"log"
)

func GetUserPasswordHash(user_id uint32) (passwordHash []byte, err error) {
	var (
		passwordHashString string
		hash_contrasenna   []byte
	)

	if err = db.QueryRow("SELECT hash_contrasenna FROM usuario where id_usuario=?",
		user_id).Scan(&hash_contrasenna); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}

	if passwordHashString, err = DecryptUserPasswordHash(hash_contrasenna); err != nil {
		log.Fatal(err)
	}

	return []byte(passwordHashString), nil
}

func GetUserSignatureKey(user_id uint32) (key []byte, err error) {
	if err = db.QueryRow("SELECT llave_firma FROM token where id_usuario=?",
		user_id).Scan(&key); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}

	return key, nil
}

func GetUserIDWithEmail(email string) (user_id uint32, err error) {
	if err = db.QueryRow("SELECT id_usuario FROM usuario WHERE correo=?",
		email).Scan(&user_id); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return 0, err
	}

	return user_id, nil
}

func NewUser(email, password, name, surname string, phoneNumber *string) (user_id uint32, err error) {
	var (
		encryptedHash string
		result        sql.Result
		lastInsertID  int64
		rowsAffected  int64
	)

	if encryptedHash, err = EncryptHashUserPassword(password); err != nil {
		log.Printf("Error at creating encrypting user's password hash: %w", err)
		return 0, err
	}

	if phoneNumber != nil {
		result, err = db.Exec("INSERT INTO usuario(correo,hash_contrasenna,nombre,apellidos,telefono)VALUES(?,UNHEX(?),?,?,?)",
			email, encryptedHash, name, surname, *phoneNumber)
	} else {
		result, err = db.Exec("INSERT INTO usuario(correo,hash_contrasenna,nombre,apellidos)VALUES(?,UNHEX(?),?,?)",
			email, encryptedHash, name, surname)
	}

	if err != nil {
		err = fmt.Errorf("Error: %w, inserting new user.", err)
		log.Print(err)
		return 0, err
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return 0, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return 0, err
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)

	return uint32(lastInsertID), nil
}

func ChangeUserPassword(user_id uint32, newPassword string) (err error) {
	var (
		result       sql.Result
		lastInsertID int64
		rowsAffected int64
	)

	if newPassword, err = EncryptHashUserPassword(newPassword); err != nil {
		log.Println(err)
	}

	if result, err = db.Exec("UPDATE usuario SET hash_contrasenna=UNHEX(?) where id_usuario=?",
		newPassword, user_id); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return err
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return err
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)

	return nil
}

type UserData struct {
	Name         string
	Surname      string
	Phone_number string
}

func GetUserData(user_id uint32) (user_data UserData, err error) {
	if err = db.QueryRow("SELECT nombre, apellidos, telefono FROM usuario WHERE id_usuario=?",
		user_id).Scan(&user_data.Name, &user_data.Surname, &user_data.Phone_number); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return UserData{}, err
	}

	return user_data, nil
}
