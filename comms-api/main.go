package main

import (
	"log"
	"net/http"

	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

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
