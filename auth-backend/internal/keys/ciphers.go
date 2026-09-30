package keys

import (
	"crypto/cipher"
)

var (
	dbBlock cipher.Block
	dbGCM   cipher.AEAD
)

// GetDBCipherBlock regresa el bloque de cifrado con la llave
// simétrica de la base de datos.
func GetDBCipherBlock() cipher.Block {
	return dbBlock
}

// GetDBCipherGCM regresa el bloque de cifrado con la llave
// simétrica de la base de datos, envuelto en Galois Counter Mode,
// se utiliza para la encriptación y desencriptación.
func GetDBCipherGCM() cipher.AEAD {
	return dbGCM
}
