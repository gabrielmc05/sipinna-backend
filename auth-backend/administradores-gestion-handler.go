package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/gin-gonic/gin"
)

// autorizarAdministradores comprueba JWT y rol Administrador
// reutilizando el mecanismo ya implementado en auth-backend.
func autorizarAdministradores(c *gin.Context) bool {
	return autorizarGestionAlimentadores(c)
}

// obtenerIDAdministrador extrae y valida el parámetro :id.
func obtenerIDAdministrador(c *gin.Context) (uint32, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ID de administrador inválido.",
		})
		return 0, false
	}

	return uint32(id), true
}

// listarAdministradoresAdmin consulta las cuentas reales.
//
// GET /admin/administradores
func listarAdministradoresAdmin(c *gin.Context) {
	if !autorizarAdministradores(c) {
		return
	}

	administradores, err := dboper.ListarAdministradores()
	if err != nil {
		log.Printf("listar administradores: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No fue posible consultar los administradores.",
		})
		return
	}

	c.JSON(http.StatusOK, administradores)
}

// editarAdministradorAdmin actualiza datos e institución.
//
// PATCH /admin/administradores/:id
func editarAdministradorAdmin(c *gin.Context) {
	if !autorizarAdministradores(c) {
		return
	}

	id, ok := obtenerIDAdministrador(c)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Los datos del administrador son inválidos.",
		})
		return
	}

	body.Nombre = strings.TrimSpace(body.Nombre)
	body.Apellidos = strings.TrimSpace(body.Apellidos)
	body.Correo = strings.ToLower(strings.TrimSpace(body.Correo))

	if body.Telefono != nil {
		telefono := strings.TrimSpace(*body.Telefono)

		if telefono == "" {
			body.Telefono = nil
		} else {
			if len(telefono) < 9 ||
				len(telefono) > 16 ||
				!strings.HasPrefix(telefono, "+") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "El teléfono debe tener formato E.164.",
				})
				return
			}

			for _, caracter := range telefono[1:] {
				if caracter < '0' || caracter > '9' {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "El teléfono debe contener solo dígitos después del signo +.",
					})
					return
				}
			}

			body.Telefono = &telefono
		}
	}

	err := dboper.EditarAdministradorActual(
		id,
		dboper.DatosAlimentadorEditado{
			Nombre:        body.Nombre,
			Apellidos:     body.Apellidos,
			Correo:        body.Correo,
			Telefono:      body.Telefono,
			InstitucionID: body.InstitucionID,
		},
	)

	switch {
	case errors.Is(err, dboper.ErrAdministradorNoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No se encontró el administrador.",
		})

	case errors.Is(err, dboper.ErrInstitucionNoDisponible):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La institución no existe o está inactiva.",
		})

	case err != nil:
		log.Printf("editar administrador: %v", err)
		c.JSON(http.StatusConflict, gin.H{
			"error": "No se pudo actualizar el administrador. Comprueba que el correo sea único.",
		})

	default:
		c.JSON(http.StatusOK, gin.H{
			"message": "Administrador actualizado en MariaDB.",
		})
	}
}

// cambiarEstadoAdministradorAdmin actualiza el estado
// y revoca la firma JWT de la cuenta seleccionada.
//
// PATCH /admin/administradores/:id/estado
func cambiarEstadoAdministradorAdmin(c *gin.Context) {
	if !autorizarAdministradores(c) {
		return
	}

	id, ok := obtenerIDAdministrador(c)
	if !ok {
		return
	}

	var body struct {
		Activo *bool `json:"es_activo" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil ||
		body.Activo == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "es_activo debe ser true o false.",
		})
		return
	}

	err := dboper.CambiarEstadoAdministradorActual(
		id,
		*body.Activo,
	)

	switch {
	case errors.Is(err, dboper.ErrAdministradorProtegido):
		c.JSON(http.StatusForbidden, gin.H{
			"error": "La cuenta administradora principal no puede desactivarse.",
		})

	case errors.Is(err, dboper.ErrAdministradorNoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No se encontró el administrador.",
		})

	case err != nil:
		log.Printf("cambiar estado administrador: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo cambiar el estado.",
		})

	default:
		c.JSON(http.StatusOK, gin.H{
			"message": "Estado del administrador actualizado.",
		})
	}
}
