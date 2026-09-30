package main

import (
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	//"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"crypto/sha256"
	"time"
)

func main() {
	//dboper.Setup()
	//dboper.Probando(3, "Hola cachorras :3")
	//fmt.Println(dboper.GetUserPassword(3))
	//dboper.CloseDB()

	dboper.Setup()
	var password string = "Hola cachorras :3"
	chap := authfuncs.NewCHAP(dboper.GetUserPassword(3))
	fmt.Printf("%X\n", chap.Challenge())
	fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte(password), chap.Challenge()...))), time.Now().Unix()+500))
	fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte("Hola"), chap.Challenge()...))), time.Now().Unix()))
	fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte(password), chap.Challenge()...))), time.Now().Unix()))
	dboper.CloseDB()

	//var (
	//	password  string = "Hola papu!"
	//	encrypted string
	//	decrypted string
	//	err       error
	//)
	//encrypted, err = dboper.EncryptUserPassword(password)
	//if err != nil {
	//	panic(err)
	//}
	//decrypted, err = dboper.DecryptUserPassword(encrypted)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Password: %s\nEncrypted password: %x\nDecrypted password: %s\n", password, encrypted, decrypted)
}
