package main

import (
	//"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	//"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
)

func main() {
	//dboper.Setup()
	//fmt.Println(dboper.GetUserPassword(3))
	//dboper.CloseDB()

	//chap := authfuncs.NewCHAP("Blablablabla1284023k")
	//fmt.Printf("%X\n", chap.Challenge())
	//fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte("Blablablabla1284023k"), chap.Challenge()...)))))
	//fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte("Hola"), chap.Challenge()...)))))

	var (
		password  string = "Hola putitas!"
		encrypted string
		decrypted string
		err       error
	)
	encrypted, err = dboper.EncryptUserPassword(password)
	if err != nil {
		panic(err)
	}
	decrypted, err = dboper.DecryptUserPassword(encrypted)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Password: %s\nEncrypted password: %x\nDecrypted password: %s\n", password, encrypted, decrypted)
}
