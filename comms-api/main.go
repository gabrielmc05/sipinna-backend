package main

import (
	"log"
	"net/http"

	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

/**
 * corsMiddleware configura las políticas CORS para permitir que
 * la aplicación web local consuma la API desde el navegador.
 *
 * Se permiten únicamente los orígenes utilizados durante el
 * desarrollo local de la aplicación.
 */
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Orígenes permitidos durante el desarrollo local.
		if origin == "http://localhost:5173" ||
			origin == "http://127.0.0.1:5173" {

			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		}

		// Responder correctamente a las solicitudes preflight.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	if err := dboper.Setup(); err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer func() {
		if err := dboper.CloseDB(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	router := gin.Default()

	router.Use(corsMiddleware())

	api := router.Group("/api/v1")

	// POST para crear un reporte con la información básica. La evidencia se sube en otro endpoint.
	// CreateReport es para la creación de un nuevo reporte. Valida la solicitud, genera un folio único 
	// y un código privado, y guarda el reporte en la base de datos.
	api.POST("/reportes", handlers.CreateReport)

	// POST para subir la evidencia
	// UploadEvidence maneja la carga de evidencias para un reporte específico.
	// Valida el folio, procesa los archivos cargados y los guarda en el sistema de archivos, 
	// registrando la información en la base de datos.
	api.POST("/reportes/:folio/evidencias", handlers.UploadEvidence)
	
	// GET para consultar un solo reporte con su folio y codigo
	// ConsultReport maneja la consulta de un reporte específico utilizando el folio y un código privado.
	// Valida la solicitud y devuelve la información del reporte si el código es correcto.
	// Si el reporte no se encuentra o el código es incorrecto, devuelve un error adecuado
	api.GET("/reportes/consultar", handlers.ConsultReport)

	// GET para consultar reporte solo con folio pero con usuario autenticado
	// ConsultFolio maneja la consulta de un reporte específico utilizando el folio y 
	// el token de autorización del usuario.
	// Valida la solicitud y devuelve la información del reporte si el usuario tiene acceso.
	// Si el reporte no se encuentra o el usuario no tiene acceso, devuelve un error adecuado.
	api.GET("/reportes/:folio", handlers.ConsultFolio)
	
	// GET para consultar la lista de reporte de un usuario en específico
	// ListUserReports maneja la consulta de todos los reportes asociados a un usuario autenticado.
	// Valida el token de autorización y devuelve la lista de reportes si el usuario es válido.
	// Si el token es inválido o no se puede consultar la base de datos, devuelve un error adecuado.
	api.GET("/reportes", handlers.ListUserReports)

	// PATCH para modificar estatus de reporte
	// UpdateReportStatus maneja la actualización del estado de un reporte específico.
	// Valida el token de autorización, el folio y el nuevo estado, y actualiza la información en la base de datos.
	// Si el usuario no tiene permiso para actualizar el reporte o si ocurre un error, devuelve un error adecuado.
	api.PATCH("/reportes/:folio/estatus", handlers.UpdateReportStatus)

	// GET de lista de instituciones
	// Institutions maneja la consulta de todas las instituciones disponibles en la base de datos.
	// Devuelve la lista de instituciones en formato JSON. Si ocurre un error al consultar la base de datos, 
	// devuelve un error adecuado.
	api.GET("/instituciones", handlers.Institutions)

	if err := http.ListenAndServe(":8085", router); err != nil {
		log.Fatalf("start HTTP server: %v", err)
	}
}
