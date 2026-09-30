package keys

import (
	"crypto/cipher"
)

var dbBlock cipher.Block

func GetDBCipherBlock() cipher.Block {
	return dbBlock
}
