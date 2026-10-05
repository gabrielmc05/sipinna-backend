package main

import (
	"crypto/rand"
	"fmt"

)

func main(){
	buffer := make([]byte, 32) 
	_, err := rand.Read(buffer)
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return
	}
	fmt.Printf("Random number: %x\n", buffer)
}
