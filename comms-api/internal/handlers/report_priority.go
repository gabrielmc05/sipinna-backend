package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/token"
	"github.com/gin-gonic/gin"
)

// GetReportPriority devuelve la prioridad de un expediente
// al personal institucional autorizado.
func GetReportPriority(c *gin.Context) {
	userID, err := token.Authenticate(
		c.GetHeader("Authorization"),
	)
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido o faltante",
			})
			return
		}

		log.Printf("authenticate priority request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no fue posible validar la autorización",
		})
		return
	}

	folio := c.Param("folio")

	priority, err := dboper.GetReportPriority(userID, folio)

	if err != nil {
		switch {
		case errors.Is(err, dboper.ErrReportNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "reporte no encontrado",
			})
		case errors.Is(err, dboper.ErrReportAccessDenied):
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no tienes permiso para consultar este expediente",
			})
		default:
			log.Printf("get report priority: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "no fue posible consultar la prioridad",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folio":     folio,
		"prioridad": priority,
	})
}

// UpdateReportPriority modifica la prioridad de un reporte.
// La autorización administrativa se valida también
// en la capa de base de datos.
func UpdateReportPriority(c *gin.Context) {
	userID, err := token.Authenticate(
		c.GetHeader("Authorization"),
	)
	if err != nil {
		if errors.Is(err, token.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido o faltante",
			})
			return
		}

		log.Printf("authenticate priority update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "no fue posible validar la autorización",
		})
		return
	}

	var request struct {
		Priority string `json:"prioridad" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "prioridad requerida",
		})
		return
	}

	request.Priority = strings.TrimSpace(request.Priority)

	if request.Priority != "Alta" &&
		request.Priority != "Media" &&
		request.Priority != "Baja" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "prioridad inválida",
		})
		return
	}

	folio := c.Param("folio")

	err = dboper.UpdateReportPriority(
		userID,
		folio,
		request.Priority,
	)

	if err != nil {
		switch {
		case errors.Is(err, dboper.ErrReportNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "reporte no encontrado",
			})
		case errors.Is(err, dboper.ErrReportAccessDenied):
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no tienes autorización para cambiar esta prioridad",
			})
		default:
			log.Printf("update report priority: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "no fue posible actualizar la prioridad",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":   "Prioridad actualizada correctamente",
		"folio":     folio,
		"prioridad": request.Priority,
	})
}
