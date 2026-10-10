// Endpoints para comms-api
package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init(){
	router = gin.Default()

	// Endpoints que se pidio para la web

	// Endpoints de reportes 
	router.GET("/reportes", func(c *gin.Context){}) // Este es el de GET / reportes
	router.GET("/reportes/:folio", func(c *gin.Context){}) // Este es el endpoint de GET /reportes/{folio}
	router.PATCH("/reportes/:folio", func(c *gin.Context){}) // Este es el endpoint de PATCH / reportes /{folio}

	// Endpoints de evidencias 
	router.POST("/reportes/:folio/evidencias", func(c *gin.Context){}) // Endpoint de POST /reportes/{folio}/evidencia 
	router.GET("/reportes/:folio/evidencias", func(c *gin.Context){}) // Enpoint de GET /reportes/ {folio}/evidencia 


	// Endpoints de territorios
	router.GET("/territorios", func(c * gin.Context){})
	router.GET("/territorios/:folio", func(c *gin.Context){})

	// Endpoint de Mapa
	router.GET("/mapa/reportes", func(c *gin.Context){})

	// Enpoints de seguimiento (la verdad no se si los necesitan me confunden las flechas)
	router.POST("/reportes/:folio/seguimientos", func(c *gin.Context){})
	router.GET("/reportes/:folio/seguimientos", func(c * gin.Context){})

	// Endpoints de Alimentadores
	router.GET("/alimentadores", func(c* gin.Context){})
	router.GET("/alimentadores/:id ", func(c* gin.Context){})
	router.POST("/alimentadores", func(c* gin.Context){})
	router.PATCH("/alimentadores/:id", func(* gin.Context){})
	
	// Este endpoints que esta en la libreta me imagino que es para la parte de las etiquetas 
	router.POST("reportes/:folio/etiquetas", func(c* gin.Context){})



}

func main(){
	router.Run(":8080")
}