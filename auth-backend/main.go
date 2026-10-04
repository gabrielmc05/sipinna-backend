package main

import (
	//"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/gin-gonic/gin"
	// "github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
)

var server *http.Server

/* TODO: add TLS config */
func init() {
	dboper.Setup()
	var router *gin.Engine = gin.Default()

	router.GET("/auth/:email", AuthRequest)
	router.POST("/auth/", Auth)
	//router.POST("/register")

	gin.SetMode(gin.ReleaseMode)

	server = &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	server.ListenAndServe()
}

func AuthRequest(c *gin.Context) {
	var (
		email     string = c.Param("email")
		user_id   uint32
		err       error
		challenge []byte
	)

	if _, err = mail.ParseAddress(email); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("User with email address: %s not found.", email),
		})
		return
	}

	if challenge, err = authfuncs.NewAuthTry(uint32(user_id)); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"user_id":   user_id,
		"message":   "Authentication attempt expected.",
		"challenge": challenge,
	})
}

type authJSON struct {
	User_id  uint32 `form:"user_id" json:"user_id" binding:"required"`
	Response string `form:"response" json:"response" binding:"required"`
}

func Auth(c *gin.Context) {
	var (
		json     authJSON
		err      error
		response []byte
		token    string
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s", err),
		})
		return
	}

	if response, err = base64.StdEncoding.DecodeString(json.Response); err != nil {
		c.Error(err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": fmt.Errorf("Error in JSON body: %s", err),
		})
		return
	}
	if len(response) != 32 {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": "Response of wrong size.",
		})
		return
	}

	if _, err = authfuncs.Authenticate(json.User_id, [32]byte(response)); err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"message": err.Error(),
		})
		return
	}

	if token, err = authfuncs.GenerateToken(json.User_id); err != nil {
		c.Error(err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}

	//c.JSON(http.StatusCreated, token)
	c.String(http.StatusCreated, token)
}

//type registerJSON struct {
//	Email        string `form:"email" json:"email" binding:"required"`
//	Password     string `form:"password" json:"password" binding:"required"`
//	Name         string `form:"name" json:"name" binding:"required"`
//	Surname      string `form:"surname" json:"surname" binding:"required"`
//	Phone_number string `form:"phone_number" json:"phone_number" binding:"required"`
//}
//
//func Register(c *gin.Context) {
//	var (
//		json registerJSON
//		err  error
//	)
//
//	if err = c.ShouldBindJSON(&json); err != nil {
//		c.Error(err)
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": fmt.Sprintf("Malformed request body: %s", err),
//		})
//		return
//	}
//
//}
