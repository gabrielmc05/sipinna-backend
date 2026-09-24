package main

import (
	//"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"fmt"
	//"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	//"github.com/ovh/symmecrypt"
)

func main() {
	//dboper.Setup()
	//fmt.Println(dboper.GetUserPassword(3))
	//dboper.CloseDB()

	//chap := authfuncs.NewCHAP("Blablablabla1284023k")
	//fmt.Printf("%X\n", chap.Challenge())
	//fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte("Blablablabla1284023k"), chap.Challenge()...)))))
	//fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte("Hola"), chap.Challenge()...)))))

	if k, err := symmecrypt.NewRandomKey("hmac"); err != nil {
		panic(err)
	} else {
		fmt.Println(k.String())
	}
}
