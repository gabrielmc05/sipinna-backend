package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "os"
)

func main() {
    passwordHash := sha256.Sum256([]byte(os.Args[1]))
    challenge, _ := base64.StdEncoding.DecodeString(os.Args[2])

    mac := hmac.New(sha256.New, passwordHash[:])
    mac.Write(challenge)

    var response [sha256.Size]byte
    copy(response[:], mac.Sum(nil))

    fmt.Printf("Hex: %x\nBase64: %s\n", response, base64.StdEncoding.EncodeToString(response[:]))
}