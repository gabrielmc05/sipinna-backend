package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// El salt es un valor aleatorio único que evita que dos usuarios con la misma
	// contraseña tengan el mismo hash. Se guarda junto al hash, no es secreto.
	saltLength = 16
	// Cantidad de bytes del resultado de Argon2.
	keyLength = 32
	// Parámetros de costo de Argon2id: iteraciones, memoria en KiB y paralelismo.
	// Aumentarlos encarece el cálculo tanto para quien verifica como para atacantes.
	argonTime    = 3
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
)

// HashPassword transforma una contraseña en un texto que contiene el algoritmo,
// sus parámetros, el salt y el hash. Devuelve error si no pudo generar aleatoriedad.
func HashPassword(password string) (string, error) {
	// make crea un slice (una secuencia de bytes) con espacio para el salt.
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Argon2id deriva un hash usando la contraseña, el salt y los parámetros definidos.
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, keyLength)

	// Base64 permite guardar bytes binarios como texto. RawStdEncoding omite el
	// relleno final (=), y el formato final incluye lo necesario para verificar luego.
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads, encodedSalt, encodedHash), nil
}

// VerifyPassword calcula de nuevo el hash con los parámetros y el salt guardados,
// y devuelve true solo si coincide con el hash almacenado.
func VerifyPassword(password, encodedHash string) (bool, error) {
	// El texto creado arriba se separa por "$"; se esperan seis partes del formato.
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("formato de hash invalido")
	}

	var memory, timeArg uint32
	var threads uint8
	// Recupera los parámetros guardados en el propio texto del hash.
	fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &timeArg, &threads)

	// Decodifica el salt y el hash esperado desde Base64 de vuelta a bytes.
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Se usan los mismos parámetros, salt y longitud que al crear el hash original.
	computedHash := argon2.IDKey([]byte(password), salt, timeArg, memory, threads, uint32(len(expectedHash)))

	// La comparación en tiempo constante reduce filtraciones sobre cuánto coincide
	// cada intento de contraseña. El resultado de Compare es 1 si son iguales.
	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1, nil
}

// Este main no inicia la aplicación: el paquete de este archivo es auth, no main.
// Parece una prueba manual incompleta, porque ignora tanto los resultados como los errores.
func main() {
	HashPassword("123456")
	VerifyPassword("123456", "$argon2id$v=19$m=65536,t=3,p=4$5a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p$7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f")

}
