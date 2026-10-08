package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gabrielmc05/sipinna-backend/auth-API/internal/authbackend"
	"github.com/gin-gonic/gin"
)

type LoginResponseInput struct {
	UserID   uint32 `json:"user_id" binding:"required"`
	Response string `json:"response" binding:"required"`
}

type RegistroInput struct {
	Nombre   string  `json:"nombre" binding:"required"`
	Apellido string  `json:"apellido" binding:"required"`
	Correo   string  `json:"correo" binding:"required,email"`
	Password string  `json:"contraseña" binding:"required"`
	Telefono *string `json:"numero_de_telefono"`
}

type CambiarContrasenaInput struct {
	NuevaContrasena string `json:"new_password" binding:"required"`
}

type EditarUsuarioInput struct {
	Nombre   *string `json:"nombre"`
	Apellido *string `json:"apellido"`
	Telefono *string `json:"numero_de_telefono"`
}

type CambiarCorreoInput struct {
	UserID      uint32 `json:"user_id" binding:"required"`
	NuevoCorreo string `json:"nuevo_correo" binding:"required,email"`
}

// Middleware: corta la petición si no trae el header Authorization.
// Solo revisa que exista; quien valida el JWT de verdad es auth-backend.
func requireAuthHeader(c *gin.Context) {
	if c.GetHeader("Authorization") == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Falta token de autenticacion"})
		return
	}
	c.Next()
}

// Traduce el error de auth-backend a una respuesta pública segura para que no se reenvía el mensaje interno
func respondBackendError(c *gin.Context, err error) {
	var backendErr *authbackend.BackendError
	if errors.As(err, &backendErr) {
		switch backendErr.Status {
		case http.StatusBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		case http.StatusUnauthorized, http.StatusNotFound:
			// 404 y 401 se responden igual a propósito: así un usuario con sesión no puede descubrir qué correos existen según el código que reciba
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		default:
			c.JSON(http.StatusBadGateway, gin.H{"error": "Error en el servicio de autenticacion"})
		}
		return
	}
	c.JSON(http.StatusBadGateway, gin.H{"error": "No se pudo contactar el servicio de autenticacion"})
}

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

// funcion que recibe los datos de registro y los envia al servicio de autenticacion para que ahí se cifre y guarde la contraseña
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

// funcion para cambiar contraseña
func cambiarContrasenna(c *gin.Context) {

	var input CambiarContrasenaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Se reenvían el correo, la nueva contraseña y el token.
	if err := authbackend.ChangePassword(c.Param("correo"), input.NuevaContrasena, c.GetHeader("Authorization")); err != nil {
		respondBackendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Contraseña actualizada correctamente"})
}

// auth-backend borra la llave de firma del usuario, así que sus tokens dejan de servir.
func logout(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identificador invalido"})
		return
	}
	if err := authbackend.Logout(uint32(id), c.GetHeader("Authorization")); err != nil {
		respondBackendError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// consultar los datos del usuario (solo con el token del propio usuario)
func getUsuario(c *gin.Context) {
	data, err := authbackend.GetUser(c.Param("correo"), c.GetHeader("Authorization"))
	if err != nil {
		respondBackendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":            data.UserID,
		"nombre":             data.Name,
		"apellido":           data.Surname,
		"numero_de_telefono": data.PhoneNumber,
	})
}

// editar nombre, apellido y/o teléfono; solo se mandan los campos que vengan en el JSON
func editarUsuario(c *gin.Context) {
	var input EditarUsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}
	if input.Nombre == nil && input.Apellido == nil && input.Telefono == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay datos para actualizar"})
		return
	}

	err := authbackend.EditUser(c.Param("correo"), authbackend.EditUserFields{
		Name:        input.Nombre,
		Surname:     input.Apellido,
		PhoneNumber: input.Telefono,
	}, c.GetHeader("Authorization"))
	if err != nil {
		respondBackendError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// cambiar el correo del usuario
func cambiarCorreo(c *gin.Context) {
	var input CambiarCorreoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}
	if err := authbackend.ChangeEmail(input.UserID, input.NuevoCorreo, c.GetHeader("Authorization")); err != nil {
		respondBackendError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func main() {

	router := gin.Default()

	router.GET("/Login/:correo", loginChallenge)                                //listo
	router.POST("/Login/", loginVerify)                                         // listo
	router.POST("/Registro/", registro)                                         // listo
	router.PATCH("/Contrasenna/:correo", requireAuthHeader, cambiarContrasenna) //listo

	router.DELETE("/Login/:id", requireAuthHeader, logout)
	router.GET("/Usuario/:correo", requireAuthHeader, getUsuario)
	router.PATCH("/Usuario/:correo", requireAuthHeader, editarUsuario)
	router.PATCH("/Correo/", requireAuthHeader, cambiarCorreo)

	router.Run()
}
