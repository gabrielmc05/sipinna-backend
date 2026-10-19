package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// autorizarGestionAlimentadores valida al administrador solicitante.
//
// Obtiene el ID declarado en el JWT, verifica su firma y expiración,
// comprueba que corresponda a ese ID y consulta su rol real.
//
// El ID obtenido antes de verificar la firma nunca se considera
// una identidad autenticada.
func autorizarGestionAlimentadores(c *gin.Context) bool {
	header := strings.TrimSpace(c.GetHeader("Authorization"))

	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Se requiere un token Bearer.",
		})
		return false
	}

	tokenString := strings.TrimSpace(
		strings.TrimPrefix(header, "Bearer "),
	)

	claims := &jwt.RegisteredClaims{}

	_, _, err := jwt.NewParser().ParseUnverified(
		tokenString,
		claims,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token inválido.",
		})
		return false
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil || id == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Identidad de sesión inválida.",
		})
		return false
	}

	adminID := uint32(id)

	valido, err := authfuncs.VerifyToken(
		header,
		adminID,
	)
	if err != nil || !valido {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Sesión inválida o expirada.",
		})
		return false
	}

	esAdmin, err := dboper.UserIsAdmin(adminID)
	if err != nil || !esAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Solo un administrador puede gestionar alimentadores.",
		})
		return false
	}

	return true
}

// obtenerIDAlimentador valida el ID de la cuenta de destino
// recibido como parámetro de la ruta.
func obtenerIDAlimentador(c *gin.Context) (uint32, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ID de alimentador inválido.",
		})
		return 0, false
	}

	return uint32(id), true
}

// editarAlimentadorAdmin modifica datos e institución.
//
// PATCH /admin/alimentadores/:id
//
// El cuerpo JSON contiene:
// nombre, apellidos, correo, telefono e institution_id.
func editarAlimentadorAdmin(c *gin.Context) {
	if !autorizarGestionAlimentadores(c) {
		return
	}

	id, ok := obtenerIDAlimentador(c)
	if !ok {
		return
	}

	var body struct {
		Nombre        string  `json:"nombre" binding:"required,min=3,max=30"`
		Apellidos     string  `json:"apellidos" binding:"required,min=3,max=100"`
		Correo        string  `json:"correo" binding:"required,email"`
		Telefono      *string `json:"telefono"`
		InstitucionID uint32  `json:"institution_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Datos inválidos: " + err.Error(),
		})
		return
	}

	body.Nombre = strings.TrimSpace(body.Nombre)
	body.Apellidos = strings.TrimSpace(body.Apellidos)
	body.Correo = strings.ToLower(strings.TrimSpace(body.Correo))

	if len([]rune(body.Nombre)) < 3 ||
		len([]rune(body.Apellidos)) < 3 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Nombre y apellidos deben tener al menos 3 caracteres.",
		})
		return
	}

	if body.Telefono != nil {
		telefono := strings.TrimSpace(*body.Telefono)
		if telefono == "" {
			body.Telefono = nil
		} else {
			// Formato E.164 para teléfonos no vacíos.
			if len(telefono) < 9 || len(telefono) > 16 ||
				!strings.HasPrefix(telefono, "+") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Utiliza un teléfono en formato internacional E.164.",
				})
				return
			}
			for _, digito := range telefono[1:] {
				if digito < '0' || digito > '9' {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "El teléfono debe contener únicamente dígitos después del signo +.",
					})
					return
				}
			}
			body.Telefono = &telefono
		}
	}

	err := dboper.EditarAlimentador(
		id,
		dboper.DatosAlimentadorEditado{
			Nombre:        body.Nombre,
			Apellidos:     body.Apellidos,
			Correo:        body.Correo,
			Telefono:      body.Telefono,
			InstitucionID: body.InstitucionID,
		},
	)

	if errors.Is(err, dboper.ErrAlimentadorNoEncontrado) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "No se encontró la cuenta alimentadora.",
		})
		return
	}

	if errors.Is(err, dboper.ErrInstitucionNoDisponible) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "La institución no existe o está inactiva.",
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "No se pudo editar el alimentador. Revisa que el correo no esté registrado.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alimentador actualizado.",
	})
}

// cambiarEstadoAlimentadorAdmin cambia es_activo
// y revoca cualquier clave JWT anterior.
//
// PATCH /admin/alimentadores/:id/estado
//
// Body: {"es_activo": false}
func cambiarEstadoAlimentadorAdmin(c *gin.Context) {
	if !autorizarGestionAlimentadores(c) {
		return
	}

	id, ok := obtenerIDAlimentador(c)
	if !ok {
		return
	}

	var body struct {
		EsActivo *bool `json:"es_activo" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil ||
		body.EsActivo == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Se requiere es_activo con true o false.",
		})
		return
	}

	err := dboper.CambiarEstadoAlimentador(id, *body.EsActivo)

	if errors.Is(err, dboper.ErrAlimentadorNoEncontrado) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "No se encontró la cuenta alimentadora.",
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo actualizar el estado.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estado del alimentador actualizado.",
	})
}
