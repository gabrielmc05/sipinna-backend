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

func loginVerify(c *gin.Context) {
	
	var input LoginResponseInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// El cliente calcula la respuesta al reto; aquí se envía al servicio
	// de autenticación para validarla y obtener un token.
	token, err := authbackend.VerifyChallenge(input.UserID, input.Response)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales invalidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func registro(c *gin.Context) {
	var input RegistroInput
	// Primero se comprueba que el cuerpo tenga un JSON válido con los campos
	// que corresponden a RegistroInput.
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

	// StatusCreated (201) indica que el recurso, en este caso el usuario,
	// fue creado correctamente.
	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
		"mensaje": "Usuario creado de manera exitosa",
	})
}

func cambiarContrasenna(c *gin.Context) {
	// El correo se obtiene de :correo en la ruta.
	correo := c.Param("correo")

	var input CambiarContrasenaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// El token de autorización normalmente llega en el encabezado HTTP
	// Authorization. El servicio de autenticación lo validará.
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
	// gin.Default crea el enrutador e incluye middleware básico, como
	// registro de solicitudes y recuperación ante errores inesperados.
	router := gin.Default()

	// Cada ruta conecta un método HTTP y una URL con la función que la atiende.
	// :correo es un parámetro de ruta que se puede leer con c.Param("correo").
	router.GET("/Login/:correo", loginChallenge)
	router.POST("/Login/", loginVerify)
	router.POST("/Registro/", registro)
	router.PUT("/Contrasenna/:correo", cambiarContrasenna)

	// Inicia el servidor HTTP. Sin argumentos, Gin escucha en el puerto 8080.
	router.Run()
}