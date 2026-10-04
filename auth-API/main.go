package main

import (
	"net/http"

	"github.com/gabrielmc05/sipinna-backend/auth-API/internal/authbackend"
	"github.com/gin-gonic/gin"
)


type LoginResponseInput struct {
	UserID   uint32 `json:"user_id"`
	Response string `json:"response"`
}

type RegistroInput struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Correo   string `json:"correo"`
	Password string `json:"contraseña"`
	Telefono string `json:"numero_de_telefono"`
}

type CambiarContrasenaInput struct {
	NuevaContrasena string `json:"nueva_contraseña"`
}


// func registro(c *gin.Context) {
// 	var input RegistroInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro invalidos"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		// mensaja de usuario creado de manera exitosa
// 		"mensaje": "Usuario creado de manera exitosa",
// 	})
// }

// primera funcion login encargada de recibir el correo y enviar el challenge al cliente
func loginChallenge(c *gin.Context) {

	correo := c.Param("correo")

	result, err := authbackend.RequestChallenge(correo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales invalidas"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"challenge": result.Challenge,
		"mensaje":   "Intento aceptado",
		"user_id":   result.UserID,
	})
}

// segunda funcion login encargada de recibir la respuesta del challenge y enviar el token al cliente 
func loginVerify(c *gin.Context) {
	
	var input LoginResponseInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// El cliente calcula la respuesta al reto, aquí se envía al servicio de autenticación para validarla y obtener un token.
	token, err := authbackend.VerifyChallenge(input.UserID, input.Response)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales invalidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func registro(c *gin.Context) {
	var input RegistroInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro invalidos"})
		return
	}

	// El servicio de autenticación realiza el registro y devuelve el ID nuevo.
	userID, err := authbackend.Register(input.Nombre, input.Apellido, input.Correo, input.Password, input.Telefono)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo completar el registro"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
		"mensaje": "Usuario creado de manera exitosa",
	})
}

func cambiarContrasenna(c *gin.Context) {
	
	correo := c.Param("correo")

	var input CambiarContrasenaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// El token de autorización en el encabezado Authorization (por verificar)
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta token de autenticacion"})
		return
	}

	// Se reenvían el correo, la nueva contraseña y el token.
	if err := authbackend.ChangePassword(correo, input.NuevaContrasena, token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo cambiar la contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Contraseña actualizada correctamente"})
}

func main() {
	
	router := gin.Default()

	router.GET("/Login/:correo", loginChallenge)
	router.POST("/Login/", loginVerify)
	router.POST("/Registro/", registro)
	router.PUT("/Contrasenna/:correo", cambiarContrasenna)

	router.Run()
}