package main

import (
// "crypto/sha256"
// "fmt"
// "github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
// "github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
// "github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
// "time"
)

func main() {
	//dboper.Setup()
	//var (
	//	password     [32]byte = sha256.Sum256([]byte("Hola cachorras :3"))
	//	passwordHash string
	//)
	//passwordHash, _ = dboper.GetUserPasswordHash(3)
	//chap := authfuncs.NewCHAP(passwordHash)
	//fmt.Printf("%X\n", chap.Challenge())
	//fmt.Println(chap.Authenticate(sha256.Sum256((append(password[:], chap.Challenge()...))), time.Now().Add(time.Duration(500*time.Second))))
	//fmt.Println(chap.Authenticate(sha256.Sum256((append([]byte("Hola"), chap.Challenge()...))), time.Now()))
	//fmt.Println(chap.Authenticate(sha256.Sum256((append(password[:], chap.Challenge()...))), time.Now()))
	//dboper.CloseDB()
}
