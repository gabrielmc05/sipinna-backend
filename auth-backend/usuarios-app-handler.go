package main

import (
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/authfuncs"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/gin-gonic/gin"
)

// adminPrincipalID identifica al administrador principal SIPINNA.
// Esta autorización debe permanecer también en el servidor.
const adminPrincipalID uint32 = 3

// autorizarAdminPrincipal valida la identidad y el rol
// de la cuenta principal del sistema.
//
// No confía en identificadores enviados por el frontend.
//
// Retorna true únicamente cuando la solicitud está autorizada.
func autorizarAdminPrincipal(c *gin.Context) bool {
	authorization := strings.TrimSpace(
		c.GetHeader("Authorization"),
	)

	if !strings.HasPrefix(authorization, "Bearer ") {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Se requiere un token.",
			},
		)
		return false
	}

	valido, err := authfuncs.VerifyToken(
		authorization,
		adminPrincipalID,
	)

	if err != nil || !valido {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error": "Acceso restringido al administrador principal.",
			},
		)
		return false
	}

	esAdmin, err := dboper.UserIsAdmin(adminPrincipalID)

	if err != nil || !esAdmin {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error": "Sin permisos administrativos.",
			},
		)
		return false
	}

	return true
}

// listarUsuariosApp consulta las cuentas registradas desde la app.
//
// GET /admin/usuarios-app
//
// Respuesta:
//
//	{
//	  "usuarios": [...],
//	  "total": 7
//	}
func listarUsuariosApp(c *gin.Context) {
	if !autorizarAdminPrincipal(c) {
		return
	}

	usuarios, err := dboper.ListarUsuariosApp()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "No se pudieron consultar los usuarios.",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"usuarios": usuarios,
			"total":    len(usuarios),
		},
	)
}

// idUsuarioApp obtiene y valida el identificador
// numérico de usuario recibido en la URL.
func idUsuarioApp(c *gin.Context) (uint32, bool) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		32,
	)

	if err != nil || id == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "ID inválido.",
			},
		)
		return 0, false
	}

	return uint32(id), true
}

// editarUsuarioApp actualiza los datos de una cuenta
// registrada desde la aplicación.
//
// PATCH /admin/usuarios-app/:id
//
// Campos:
//   - nombre
//   - apellidos
//   - correo
//   - telefono (opcional)
func editarUsuarioApp(c *gin.Context) {
	if !autorizarAdminPrincipal(c) {
		return
	}

	id, ok := idUsuarioApp(c)
	if !ok {
		return
	}

	var datos struct {
		Nombre string `json:"nombre" binding:"required,min=3,max=30"`

		Apellidos string `json:"apellidos" binding:"required,min=3,max=100"`

		Correo string `json:"correo" binding:"required,email"`

		Telefono *string `json:"telefono"`
	}

	if err := c.ShouldBindJSON(&datos); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Datos inválidos: " + err.Error(),
			},
		)
		return
	}

	datos.Nombre = strings.TrimSpace(datos.Nombre)
	datos.Apellidos = strings.TrimSpace(datos.Apellidos)
	datos.Correo = strings.ToLower(
		strings.TrimSpace(datos.Correo),
	)

	if len(datos.Nombre) < 3 || len(datos.Apellidos) < 3 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Nombre o apellidos demasiado cortos.",
			},
		)
		return
	}

	if _, err := mail.ParseAddress(datos.Correo); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Correo inválido.",
			},
		)
		return
	}

	if datos.Telefono != nil &&
		strings.TrimSpace(*datos.Telefono) == "" {
		datos.Telefono = nil
	}

	actualizado, err := dboper.EditarUsuarioApp(
		id,
		datos.Nombre,
		datos.Apellidos,
		datos.Correo,
		datos.Telefono,
	)

	if err != nil {
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": "No se pudo actualizar. Verifica que el correo no esté registrado.",
			},
		)
		return
	}

	if !actualizado {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "Usuario de la aplicación no encontrado o sin cambios.",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Usuario actualizado.",
		},
	)
}

// cambiarEstadoUsuarioApp activa o desactiva una cuenta.
//
// PATCH /admin/usuarios-app/:id/estado
//
// Body:
//
//	{
//	  "es_activo": false
//	}
//
// La persistencia y revocación JWT se realizan
// mediante una transacción en dboper.
func cambiarEstadoUsuarioApp(c *gin.Context) {
	if !autorizarAdminPrincipal(c) {
		return
	}

	id, ok := idUsuarioApp(c)
	if !ok {
		return
	}

	var datos struct {
		EsActivo *bool `json:"es_activo" binding:"required"`
	}

	if err := c.ShouldBindJSON(&datos); err != nil ||
		datos.EsActivo == nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Se requiere es_activo true o false.",
			},
		)
		return
	}

	actualizado, err := dboper.CambiarEstadoUsuarioApp(
		id,
		*datos.EsActivo,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "No se pudo cambiar el estado.",
			},
		)
		return
	}

	if !actualizado {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "Usuario de la aplicación no encontrado.",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Estado actualizado.",
		},
	)
}
