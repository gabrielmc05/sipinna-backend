package dboper

import (
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"strings"
)

// GetUserPasswordHash obtiene de la base de datos la salida Hash de
// la contraseña, y la desencripta usando la llave simétrica.
func GetUserPasswordHash(user_id uint32) (passwordHash []byte, err error) {
	var (
		passwordHashString string
		hash_contrasenna   []byte
	)

	if err = db.QueryRow("SELECT hash_contrasenna FROM usuario where id_usuario=? AND es_activo=1",
		user_id).Scan(&hash_contrasenna); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}

	if passwordHashString, err = DecryptUserPasswordHash(hash_contrasenna); err != nil {
		log.Printf("Error al descifrar el hash del usuario %d: %v", user_id, err)
		return nil, fmt.Errorf("no fue posible descifrar el hash de contraseña: %w", err)
	}

	return []byte(passwordHashString), nil
}

// GetUserSignatureKey obtiene la llave de firma de los JWT, se usa
// para confirmación de identidad.
func GetUserSignatureKey(user_id uint32) (key []byte, err error) {
	if err = db.QueryRow("SELECT t.llave_firma FROM token t JOIN usuario u ON u.id_usuario = t.id_usuario WHERE t.id_usuario = ? AND u.es_activo = 1",
		user_id).Scan(&key); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}

	return key, nil
}

// GetUserIDWithEmail recibe una dirección de correo electrónico
// y consulta si hay un usuario registrado en la base datos con esa
// dirección.
func GetUserIDWithEmail(email string) (user_id uint32, err error) {
	if err = db.QueryRow("SELECT id_usuario FROM usuario WHERE correo=?",
		email).Scan(&user_id); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return 0, err
	}

	return user_id, nil
}

// UserExists recibe un ID y consulta a la base de datos si existe.
func UserExists(user_id uint32) (ok bool, err error) {
	if err = db.QueryRow("SELECT ExisteUsuario(?)", user_id).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

// UserIsAdmin recibe un ID y consulta a la base de datos si es un
// usuario administrador.
func UserIsAdmin(user_id uint32) (ok bool, err error) {
	if ok, err = UserExists(user_id); err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("User doesn't exist.")
	}
	if err = db.QueryRow("SELECT EsAdministrador(?)", user_id).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

// NewUser registra un usuario en la base de datos, recibiendo sus
// datos personales, regresa un error si no pudo realizarse el
// registro.
func NewUser(email, password, name, surname string, phoneNumber *string) (user_id uint32, err error) {
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

	if result, err = db.Exec(
		"INSERT INTO usuario(correo,hash_contrasenna,nombre,apellidos,telefono)VALUES(?,UNHEX(?),?,?,?)",
		email, encryptedHash, name, surname, *phoneNumber); err != nil {
		err = fmt.Errorf("Error: %w, inserting new user.", err)
		log.Print(err)
		return 0, err
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

// Role es ara definir constantes Rol al registrar en la base de datos
type Role uint8

const (
	Administrator Role = 0
	Feeder        Role = 1
	User          Role = 2
)

// Admin se usa para registrar un usuario administrador
type Admin struct {
	Email         string
	Password      string
	Name          string
	Surname       string
	PhoneNumber   *string
	InstitutionId uint32
	Role          Role
}

// NewAdmin registra un usuario administrador en la base de datos,
// recibiendo sus datos personales y profesionales, regresa un error
// si no pudo realizarse el registro.
func NewAdmin(admin *Admin) (user_id uint32, user_name string, err error) {
	var (
		encryptedHash string
		query         string
	)

	if _, err = mail.ParseAddress(admin.Email); err != nil {
		err = fmt.Errorf("Email address invalid: %w.", err)
		return 0, "", err
	}

	if encryptedHash, err = EncryptHashUserPassword(admin.Password); err != nil {
		log.Printf("Error at creating encrypting user's password hash: %w", err)
		return 0, "", err
	}

	switch admin.Role {
	case Administrator:
		query = "CALL RegistrarAdministrador(?,?,?,?,?,?)"
	case Feeder:
		query = "CALL RegistrarAlimentador(?,?,?,?,?,?)"
	default:
		return 0, "", fmt.Errorf("Invalid role, expected values: 'Administrator', 'Feeder'.")
	}

	if err = db.QueryRow(query,
		admin.Email, encryptedHash, admin.Name, admin.Surname, admin.PhoneNumber, admin.InstitutionId,
	).Scan(&user_id, &user_name); err != nil {
		err = fmt.Errorf("Error: %w, inserting new admin.", err)
		log.Print(err)
		return 0, "", err
	}

	return user_id, user_name, nil
}

// ChangeUserPassword actualiza la contraseña de un usuario.
//
// Primero genera el hash SHA-256 y lo cifra usando
// el mecanismo existente de autenticación.
//
// Si el procesamiento falla, no modifica MariaDB.

func ChangeUserPassword(user_id uint32, newPassword string) error {
	encryptedHash, err := EncryptHashUserPassword(newPassword)
	if err != nil {
		return fmt.Errorf(
			"error al procesar la nueva contraseña: %w",
			err,
		)
	}

	_, err = db.Exec(
		"UPDATE usuario SET hash_contrasenna = UNHEX(?) WHERE id_usuario = ?",
		encryptedHash,
		user_id,
	)

	if err != nil {
		return fmt.Errorf(
			"error al actualizar la contraseña: %w",
			err,
		)
	}

	return nil
}

func ChangeUserEmail(user_id uint32, newEmail string) (err error) {
	if _, err = db.Exec("UPDATE usuario SET correo=? where id_usuario=?",
		newEmail, user_id); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return err
	}
	return nil
}

// ChangeUserDatum es para poder cambiar número de teléfono, nombre,
// o apellidos de un usuario o administrador.
func ChangeUserDatum(user_id uint32, field, newDatum string) (err error) {
	var query strings.Builder

	query.WriteString("UPDATE usuario SET ")
	query.WriteString(field)
	query.WriteString("=? where id_usuario=?")

	if _, err = db.Exec(query.String(),
		newDatum, user_id); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return err
	}
	return nil
}

/*
UserData contiene la información del perfil de un usuario.

Los campos de organización y municipio son opcionales porque
no todos los usuarios necesariamente pertenecen a una institución.
*/
type UserData struct {
	Name         string
	Surname      string
	Phone_number *string
	Organizacion *string
	Municipio    *string
}

/*
GetUserData consulta los datos del usuario autenticado.

Relaciona la cuenta con su institución y el territorio
correspondiente, sin modificar información en MariaDB.
*/
func GetUserData(user_id uint32) (user_data UserData, err error) {
	const query = `
		SELECT
			u.nombre,
			u.apellidos,
			u.telefono,
			i.nombre,
			t.nombre
		FROM usuario AS u
		LEFT JOIN administrador AS a
			ON a.id_usuario = u.id_usuario
		LEFT JOIN institucion AS i
			ON i.id_institucion = a.id_institucion
		LEFT JOIN territorio AS t
			ON t.id_territorio = i.id_territorio
		WHERE u.id_usuario = ?
	`

	err = db.QueryRow(query, user_id).Scan(
		&user_data.Name,
		&user_data.Surname,
		&user_data.Phone_number,
		&user_data.Organizacion,
		&user_data.Municipio,
	)

	if err != nil {
		err = fmt.Errorf("error al consultar datos del usuario: %w", err)
		log.Print(err)
		return UserData{}, err
	}

	return user_data, nil
}

// DeleteUserSignatureKey borra la llave de firma de JWT actual,
// se utiliza al hacer un cerrado de sesión
func DeleteUserSignatureKey(user_id uint32) (err error) {
	if _, err = db.Exec("UPDATE token SET llave_firma=NULL where id_usuario=?",
		user_id); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return err
	}
	return nil
}

// GetUserRole recibe un ID y regresa el rol del usuario con ese ID,
// esto solo se cumple si el usuario es administrador.
func GetUserRole(user_id uint32) (role string, err error) {
	if err = db.QueryRow("SELECT rol from administrador where id_usuario = ?", user_id).Scan(&role); err != nil {
		return "", err
	}

	return role, nil
}
