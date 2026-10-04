package main

import (
	//"crypto/tls"
	"encoding/base64"
	"fmt"
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
	//gin.SetMode(gin.ReleaseMode)
	var router *gin.Engine = gin.Default()

	router.GET("/auth/:email", authRequest)
	router.POST("/auth/", auth)
	router.POST("/register", register)
	router.PUT("/password/:email", password)

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

func authRequest(c *gin.Context) {
	var (
		email     string = c.Param("email")
		user_id   uint32
		err       error
		challenge []byte
	)

	if _, err = mail.ParseAddress(email); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
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
	User_id  uint32 `json:"user_id" binding:"required"`
	Response string `json:"response" binding:"required"`
}

func auth(c *gin.Context) {
	var (
		json     authJSON
		err      error
		response []byte
		token    string
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if response, err = base64.StdEncoding.DecodeString(json.Response); err != nil {
		c.Error(err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": fmt.Errorf("Error in JSON body: %s.", err),
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
			"error": err.Error(),
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

type registerJSON struct {
	Email        string  `json:"email" binding:"required"`
	Password     string  `json:"password" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Surname      string  `json:"surname" binding:"required"`
	Phone_number *string `json:"phone_number"`
}

func register(c *gin.Context) {
	var (
		json    registerJSON
		err     error
		user_id uint32
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if _, err = mail.ParseAddress(json.Email); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.NewUser(
		json.Email, json.Password, json.Name, json.Surname, json.Phone_number); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New user created successfully.",
		"user_id": user_id,
	})
}

type authHeader struct {
	JWT string `header:"Authorization" binding:"required"`
}

type passwordJSON struct {
	New_password string `json:"new_password" binding:"required"`
}

func password(c *gin.Context) {
	var (
		header  authHeader
		email   string = c.Param("email")
		json    passwordJSON
		err     error
		user_id uint32
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", "lol")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if _, err = mail.ParseAddress(email); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User with email address: %s not found.", email),
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, user_id); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", "lol")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	if err = dboper.ChangeUserPassword(user_id, json.New_password); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Password changed for user with id: %d.", user_id),
	})
}
