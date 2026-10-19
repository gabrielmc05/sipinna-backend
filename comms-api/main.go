package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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

func newRouter() *gin.Engine {
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

	// GET para obtener reportes geolocalizados visibles al usuario institucional autenticado.
	// devuelve folio, coordenadas, estado del expediente y datos del territorio.
	// Requiere JWT, los administradores ven todos los reportes y los alimentadores solo los asignados.
	api.GET("/mapa/reportes", handlers.GetMapReports)

	// GET para obtener territorios; permite filtrar por estado y municipio.
	api.GET("/territorios", handlers.GetTerritories)
	api.GET("/territorios/:id", handlers.GetTerritory)

	// GET para consultar las cuentas de alimentadores (solo administradores).
	api.GET("/usuario/alimentadores", handlers.GetFeeders)

	// GET para consultar reporte solo con folio pero con usuario autenticado
	// ConsultFolio maneja la consulta de un reporte específico utilizando el folio y
	// el token de autorización del usuario.
	// Valida la solicitud y devuelve la información del reporte si el usuario tiene acceso.
	// Si el reporte no se encuentra o el usuario no tiene acceso, devuelve un error adecuado.
	api.GET("/reportes/:folio", handlers.ConsultFolio)

	// GET para consultar el histórico global (administradores) o los reportes asignados (alimentadores).
	// acepta filtros estado, id_territorio y orden_fecha=asc|desc. El administrador ve todos los reportes,
	// el alimentador, solo los asignados a su ID.
	api.GET("/reportes/historico", handlers.ReportHistory)

	// POST para asignar un reporte a un alimentador.
	// recibe {"id_alimentador": 123}. Solo un administrador puede asignar, y el ID debe pertenecer
	// a un alimentador. No se restringe la asignación por institución
	api.POST("/reportes/asignar/:folio", handlers.AssignReport)

	// GetEvidence consulta las evidencias registradas de un expediente.
	// Valida la autenticación y los permisos del usuario.
	api.GET("/reportes/:folio/evidencias", handlers.GetEvidence)

	// DownloadEvidence permite visualizar o descargar una evidencia
	// únicamente después de validar la autorización del usuario.
	api.GET("/reportes/:folio/evidencias/:id/archivo", handlers.DownloadEvidence)

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

	// Consultar la prioridad de un expediente.
	// Administradores y alimentadores autorizados.
	api.GET("/reportes/:folio/prioridad", handlers.GetReportPriority)

	// Modificar la prioridad de un expediente.
	// Solo administradores autorizados.
	api.PATCH("/reportes/:folio/prioridad", handlers.UpdateReportPriority)

	// GET de lista de instituciones
	// Institutions maneja la consulta de todas las instituciones disponibles en la base de datos.
	// Devuelve la lista de instituciones en formato JSON. Si ocurre un error al consultar la base de datos,
	// devuelve un error adecuado.
	api.GET("/instituciones", handlers.Institutions)
	api.GET("/instituciones/:id", handlers.InstitutionByID)
	api.POST("/instituciones", handlers.CreateInstitution)
	api.PATCH("/instituciones/:id", handlers.UpdateInstitution)
	api.PATCH("/instituciones/:id/estado", handlers.UpdateInstitutionStatus)

	// Catálogo y asignación de etiquetas para administradores.
	api.GET("/etiquetas", handlers.ListTags)
	api.POST("/etiquetas", handlers.CreateTag)
	api.PATCH("/etiquetas", handlers.RenameTag)
	api.POST("/etiquetas/aplicar", handlers.ApplyTags)

	return router
}

var ginLambda *ginadapter.GinLambdaV2

func init() {
	if err := dboper.Setup(); err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	//defer func() {
	//	if err := dboper.CloseDB(); err != nil {
	//		log.Printf("close database: %v", err)
	//	}
	//}()
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

/**
 * main inicia comms-api según el entorno de ejecución.
 *
 * En AWS Lambda utiliza el adaptador de API Gateway V2.
 * En desarrollo local inicia el servidor HTTP de Gin.
 */
func main() {
	// AWS Lambda proporciona esta variable automáticamente.
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		ginLambda = ginadapter.NewV2(newRouter())
		log.Println("Iniciando comms-api en AWS Lambda")
		lambda.Start(handler)
		return
	}

	// Ejecución local.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Iniciando comms-api local en puerto %s", port)

	router := newRouter()

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("No fue posible iniciar comms-api: %v", err)
	}
	defer func() {
		if err := dboper.CloseDB(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()
}
