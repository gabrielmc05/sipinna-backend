package main

import (
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	// "github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"github.com/gin-gonic/gin"
	"net/http"
	//"reflect"
	"encoding/base64"
	"strconv"
)

var r *gin.Engine

func init() {
	dboper.Setup()
	r = gin.Default()

	r.GET("/auth/:id", AuthRequest)
	r.POST("/auth/", Auth)
}

func main() {
	r.Run()
}

func AuthRequest(c *gin.Context) {
	var (
		id        string = c.Param("id")
		user_id   uint64
		err       error
		challenge []byte
	)

	if user_id, err = strconv.ParseUint(id, 10, 32); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request (expected unsigned integer).",
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
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//response = make([]byte, base64.RawStdEncoding.DecodedLen(len(json.Response)))
	//copy([]byte(json.Response), response)
	if response, err = base64.StdEncoding.DecodeString(json.Response); err != nil {
		c.Error(err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(response) != 32 {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"message": "Response of wrong size.",
		})
		return
	}

	if _, err = authfuncs.Authenticate(json.User_id, [32]byte(response)); err != nil {
		c.Error(err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": "",
	})
}
