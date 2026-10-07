package dboper

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"log"
)

// EncryptHashUserPassword se utiliza para aplicarles el algoritmo sha256
// y luego, encriptar las contraseñas con la llave simétrica de la base de
// datos, se usa en el registro.
func EncryptHashUserPassword(password string) (encryptedHash string, err error) {
	var (
		passwordBytes [32]byte = sha256.Sum256([]byte(password))
		nonce         []byte
	)

	nonce = make([]byte, keys.GetDBCipherGCM().NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		log.Printf("Nonce (random 32 bytes) could not be generated: %w.", err)
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
		result       sql.Result
		lastInsertID int64
		rowsAffected int64
		llave_firma  string
	)

	key = make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		log.Printf("Error at creating signature key for user with id %d: %w", user_id, err)
		return nil, err
	}

	llave_firma = hex.EncodeToString(key)
	if result, err = db.Exec(
		"INSERT INTO token (id_usuario,llave_firma) VALUES (?,UNHEX(?)) ON DUPLICATE KEY UPDATE llave_firma=UNHEX(?)",
		user_id, llave_firma, llave_firma); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}

	if lastInsertID, err = result.LastInsertId(); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		err = fmt.Errorf("Error: in query: %w", err)
		log.Print(err)
		return nil, err
	}
	log.Printf("Last Insert ID: %d, Rows Affected: %d", lastInsertID, rowsAffected)

	return key, nil
}

func GetTrustedProxies() (result []string, err error) {
	var (
		proxy string
		rows  *sql.Rows
	)

	if rows, err = db.Query("SELECT IP FROM TRUSTED_PROXIES"); err != nil {
		err = fmt.Errorf("Error at retrieving trusted proxies: %w", err)
		log.Fatal(err)
	}

	for rows.Next() {
		if err = rows.Scan(&proxy); err != nil {
			return result, err
		}
		result = append(result, proxy)
	}
	if err = rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
