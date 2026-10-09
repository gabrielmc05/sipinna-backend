package main

import (
	//"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
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

	trustedProxies, err := dboper.GetTrustedProxies()
	if err != nil {
		log.Printf("Error at retrieving trusted proxies: %s", err)
	}
	router.SetTrustedProxies(trustedProxies)

	router.GET("/auth/:email", authRequest) 
	router.POST("/auth", auth)
	router.DELETE("/auth/:id", authDelete)
	router.POST("/register", register) // listo
	router.PATCH("/password/:email", password)
	router.GET("/user/:email", user)
	router.PATCH("/user/:email", userEdit)
	router.PATCH("/email", email)
	router.POST("/register/admin/:id", registerAdmin)
	//router.POST("/register/feeder", registerFeeders)

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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User with email address: %s not found.", email),
		})
		return
	}

	if challenge, err = authfuncs.NewAuthTry(uint32(user_id)); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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

func auth(c *gin.Context) {
	var (
		json struct {
			User_id  uint32 `json:"user_id" binding:"required"`
			Response string `json:"response" binding:"required,base64"`
		}
		err      error
		response []byte
		token    string
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if response, err = base64.StdEncoding.DecodeString(json.Response); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"error": fmt.Errorf("Error in JSON body: %s.", err),
		})
		return
	} 
	if len(response) != 32 {
		c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"error": "Response of wrong size.",
		})
		return
	}

	if _, err = authfuncs.Authenticate(json.User_id, [32]byte(response)); err != nil {
		c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}

	if token, err = authfuncs.GenerateToken(json.User_id); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}

	//c.JSON(http.StatusCreated, token)
	c.String(http.StatusCreated, token)
}

type registerJSON struct {
	Email        string  `json:"email" binding:"required,email"`
	Password     string  `json:"password" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Surname      string  `json:"surname" binding:"required"`
	Phone_number *string `json:"phone_number" binding:"e164"`
}

func register(c *gin.Context) {
	var (
		json    registerJSON
		err     error
		user_id uint32
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if user_id, err = dboper.NewUser(
		json.Email, json.Password, json.Name, json.Surname, json.Phone_number); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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

func password(c *gin.Context) {
	var (
		header authHeader
		email  string = c.Param("email")
		json   struct {
			New_password string `json:"new_password" binding:"required"`
		}
		err     error
		user_id uint32
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if _, err = mail.ParseAddress(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User with email address: %s not found.", email),
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, user_id); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	if err = dboper.ChangeUserPassword(user_id, json.New_password); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func user(c *gin.Context) {
	var (
		header    authHeader
		email     string = c.Param("email")
		user_id   uint32
		err       error
		user_data dboper.UserData
	)

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if _, err = mail.ParseAddress(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User with email address: %s not found.", email),
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, user_id); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	if user_data, err = dboper.GetUserData(user_id); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":         user_data.Name,
		"surname":      user_data.Surname,
		"phone_number": user_data.Phone_number,
		"user_id":      user_id,
	})
}

func userEdit(c *gin.Context) {
	var (
		header authHeader
		json   struct {
			Name         *string `json:"name"`
			Surname      *string `json:"surname"`
			Phone_number *string `json:"phone_number"`
		}
		email   string = c.Param("email")
		user_id uint32
		err     error
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
	}

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if _, err = mail.ParseAddress(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address.",
		})
		return
	}

	if user_id, err = dboper.GetUserIDWithEmail(email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User with email address: %s not found.", email),
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, user_id); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	if json.Name != nil {
		if err = dboper.ChangeUserDatum(user_id, "nombre", *json.Name); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
	if json.Surname != nil {
		if err = dboper.ChangeUserDatum(user_id, "apellidos", *json.Surname); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
	if json.Phone_number != nil {
		if err = dboper.ChangeUserDatum(user_id, "telefono", *json.Phone_number); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func email(c *gin.Context) {
	var (
		header authHeader
		json   struct {
			User_id   uint32 `json:"user_id" binding:"required"`
			New_email string `json:"new_email" binding:"required,email"`
		}
		err error
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, json.User_id); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	if err = dboper.ChangeUserEmail(json.User_id, json.New_email); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func authDelete(c *gin.Context) {
	var (
		header  authHeader
		err     error
		user_id uint64
	)

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if user_id, err = strconv.ParseUint(c.Param("id"), 10, 32); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Malformed request parameter.",
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, uint32(user_id)); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	if err = dboper.DeleteUserSignatureKey(uint32(user_id)); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func registerAdmin(c *gin.Context) {
	var (
		header authHeader
		json   struct {
			Institution_id uint32       `json:"institution_id" binding:"required"`
			New_admin      registerJSON `json:"new_admin" binding:"required"`
		}
		err         error
		user_id     uint32
		admin_id    uint64
		user_name   string
		userIsAdmin bool
		newAdmin    dboper.Admin
	)

	if err = c.ShouldBindJSON(&json); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Malformed request body: %s.", err),
		})
		return
	}

	if err = c.ShouldBindHeader(&header); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token expected.",
		})
		return
	}

	if admin_id, err = strconv.ParseUint(c.Param("id"), 10, 32); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Malformed request parameter.",
		})
		return
	}

	if userIsAdmin, err = dboper.UserIsAdmin(uint32(admin_id)); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !userIsAdmin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("User with id %d, is not admin.", admin_id),
		})
		return
	}

	if _, err = authfuncs.VerifyToken(header.JWT, uint32(admin_id)); err != nil {
		c.Error(err)
		c.Header("WWW-Authenticate", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token invalid.",
		})
		return
	}

	newAdmin = dboper.Admin{
		Email:         json.New_admin.Email,
		Password:      json.New_admin.Password,
		Name:          json.New_admin.Name,
		Surname:       json.New_admin.Surname,
		PhoneNumber:   json.New_admin.Phone_number,
		InstitutionId: json.Institution_id,
		Role:          dboper.Administrator,
	}

	if user_id, user_name, err = dboper.NewAdmin(&newAdmin); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "New administrator created successfully.",
		"user_id":   user_id,
		"user_name": user_name,
	})
}

//func registerFeeders(c *gin.Context) {
//	var (
//		header authHeader
//		json   struct {
//			Institution_id uint32         `json:"institution_id" binding:"required"`
//			New_admin      []registerJSON `json:"new_feeders" binding:"required"`
//		}
//		err         error
//		user_id     uint32
//		admin_id    uint64
//		user_name   string
//		userIsAdmin bool
//		newAdmin    dboper.Admin
//	)
//
//	if err = c.ShouldBindJSON(&json); err != nil {
//		c.Error(err)
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
//			"error": fmt.Sprintf("Malformed request body: %s.", err),
//		})
//		return
//	}
//
//	if err = c.ShouldBindHeader(&header); err != nil {
//		c.Error(err)
//		c.Header("WWW-Authenticate", err.Error())
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": "Authorization token expected.",
//		})
//		return
//	}
//
//	if admin_id, err = strconv.ParseUint(c.Param("id"), 10, 32); err != nil {
//		c.Error(err)
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
//			"error": "Malformed request parameter.",
//		})
//		return
//	}
//
//	if userIsAdmin, err = dboper.UserIsAdmin(uint32(admin_id)); err != nil {
//		c.Error(err)
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	if !userIsAdmin {
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": fmt.Sprintf("User with id %d, is not admin.", admin_id),
//		})
//		return
//	}
//
//	if _, err = authfuncs.VerifyToken(header.JWT, uint32(admin_id)); err != nil {
//		c.Error(err)
//		c.Header("WWW-Authenticate", err.Error())
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": "Authorization token invalid.",
//		})
//		return
//	}
//
//	newAdmin = dboper.Admin{
//		Email:         json.New_admin.Email,
//		Password:      json.New_admin.Password,
//		Name:          json.New_admin.Name,
//		Surname:       json.New_admin.Surname,
//		PhoneNumber:   json.New_admin.Phone_number,
//		InstitutionId: json.Institution_id,
//		Role:          dboper.Administrator,
//	}
//
//	if user_id, user_name, err = dboper.NewAdmin(&newAdmin); err != nil {
//		c.Error(err)
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{
//		"message":   "New administrator created successfully.",
//		"user_id":   user_id,
//		"user_name": user_name,
//	})
//}
