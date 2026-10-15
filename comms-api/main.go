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
	api.POST("/reportes", handlers.CreateReport)
	api.POST("/reportes/:folio/evidencias", handlers.UploadEvidence)
	api.GET("/reportes/consultar", handlers.ConsultReport)
	api.GET("/reportes/:folio", handlers.ConsultFolio)
	api.GET("/reportes", handlers.ListUserReports)
	api.PATCH("/reportes/:folio/estatus", handlers.UpdateReportStatus)

	api.GET("/instituciones", handlers.Institutions)

	if err := http.ListenAndServe(":8085", router); err != nil {
		log.Fatalf("start HTTP server: %v", err)
	}
}
