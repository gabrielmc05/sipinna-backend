package dboper

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"log"
	"strconv"
	"strings"
)

// EncryptHashUserPassword se utiliza para aplicarles el algoritmo Argon2id
// y luego, encriptar las contraseñas con la llave simétrica de la base de
// datos, se usa en el registro.
func EncryptHashUserPassword(password string) (encryptedHash string, err error) {
	var (
		passwordBytes [32]byte = sha256.Sum256([]byte(password))
		nonce         []byte
	)

	nonce = make([]byte, keys.GetDBCipherGCM().NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		log.Printf("Nonce (random 32 bytes) could not be generated: %s.", err)
		return "", err
	}

	return hex.EncodeToString(keys.GetDBCipherGCM().Seal(nonce, nonce, passwordBytes[:], nil)), err
}

// DecryptUserPasswordHash se utiliza para desencriptar los Hash (Argon2id)
// de contraseñas con la llave simétrica de la base de datos.
func DecryptUserPasswordHash(encrypted []byte) (password string, err error) {
	var passwordBytes []byte

	if len(encrypted) < keys.GetDBCipherGCM().NonceSize() {
		return "", fmt.Errorf("EncryptedBytes too short.")
	}

	passwordBytes, err = keys.GetDBCipherGCM().Open(
		nil,
		encrypted[:keys.GetDBCipherGCM().NonceSize()],
		encrypted[keys.GetDBCipherGCM().NonceSize():],
		nil,
	)
	return string(passwordBytes), err
}

func NewUserSignatureKey(user_id uint32) (key []byte, err error) {
	var (
		query        strings.Builder
		result       sql.Result
		lastInsertID int64
		rowsAffected int64
		llave_firma  string
	)

	key = make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		log.Printf("Error at creating signature key for user with id %d: %s", user_id, err)
		return nil, err
	}

	llave_firma = hex.EncodeToString(key)
	query.WriteString("INSERT INTO token(id_usuario,llave_firma)VALUES(")
	query.WriteString(strconv.FormatUint(uint64(user_id), 10))
	query.WriteString(",UNHEX('")
	query.WriteString(llave_firma)
	query.WriteString("'))ON DUPLICATE KEY UPDATE llave_firma=UNHEX('")
	query.WriteString(llave_firma)
	query.WriteString("');")

	if result, err = db.Exec(query.String()); err != nil {
		log.Printf("Error in query: %s", err)
		return nil, err
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		log.Printf("Error in query: %s", err)
		return nil, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		log.Printf("Error in query: %s", err)
		return nil, err
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)

	return key, nil
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
