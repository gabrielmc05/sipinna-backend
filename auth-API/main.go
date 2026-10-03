package main

import (
	"net/http"
	"database/sql"

	"github.com/gabrielmc05/sipinna-backend/auth-API/internal/authbackend"
	"github.com/gabrielmc05/sipinna-backend/auth-API/internal/db"
	"github.com/gin-gonic/gin"
)

// net/http aporta los códigos de estado HTTP y database/sql permite reconocer
// errores específicos de consultas. Los otros imports son paquetes del proyecto
// y Gin, el framework que recibe las solicitudes y construye las respuestas.

// Referencias Bibliograficas
// https://gin-gonic.com/en/docs/routing/api-design/

// Estas estructuras describen los datos que la API espera recibir en el cuerpo
// (body) de las solicitudes. Los textos dentro de `json:"..."` son los nombres
// de las propiedades que deben venir en el JSON enviado por el cliente.
type LoginInput struct {
	// Los nombres después de json: son las claves esperadas en el JSON.
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type RegistroInput struct {
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Correo    string `json:"correo"`
	Password  string `json:"password"`
}

type RecuperarPasswordInput struct {
	Correo string `json:"correo"`
}


// func registro(c *gin.Context) {
// 	var input RegistroInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro invalidos"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		// mensaja de usuario creado de manera exitosa
// 		"mensaje": "Usuario creado de manera exitosa",
// 	})
// }


func login(c *gin.Context) {
	var input LoginInput
	// ShouldBindJSON convierte el cuerpo JSON de la solicitud a la estructura
	// LoginInput. &input pasa su dirección para que Gin pueda llenarla.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	} 

	// el valor se envía separado de la consulta y no se concatena al texto SQL.
	// Scan copia la columna obtenida dentro de idUsuario.
	var idUsuario uint32
	err := db.DB.QueryRow("SELECT id_usuario FROM usuario WHERE correo = ?", input.Correo).Scan(&idUsuario)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró un usuario con ese correo"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar el usuario"})
		return
	} 

	// Se pide al auth-backend un reto aleatorio para este usuario.
	challenge, err := authbackend.RequestChallenge(idUsuario)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales invalidas"})
		return
	}

	// Se envían la contraseña y el reto al cliente interno, que calcula la prueba
	// criptográfica y pide al servicio que la valide. Si es correcta, devuelve un token.
	token, err := authbackend.VerifyChallenge(idUsuario, input.Password, challenge)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales invalidas"})
		return
	}

	// c.JSON serializa este mapa como JSON. 200 OK indica que el login terminó bien.
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Login exitoso",
		"token":   token,
	})
}

func recuperarPassword(c *gin.Context) {
	var input RecuperarPasswordInput
	// Lee y valida que el cuerpo tenga la forma esperada; input.Correo es el correo recibido.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un correo"})
		return
	}
	// Esta respuesta también es provisional: aquí todavía no se envía un correo
	// ni se crea un proceso real de recuperación.
	c.JSON(http.StatusOK, gin.H{
		// mensaje de correo de recuperacion
		"mensaje": "Correo de recuperacion enviado",
	})
}

// Esta ruta podría usarse al abrir la app para consultar si la sesión sigue activa.
// Actualmente siempre contesta que sí y asigna un rol fijo; no valida una sesión real.
func verificarSesion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"valido": true,
		"rol":    "Ciudadano",
	})
}

func logout(c *gin.Context) {
	// Por ahora solo confirma la acción. No invalida una sesión o token.
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Sesion cerrada de manera exitosa",
	})
}

func main() {
	// Primero intenta conectar a la base de datos. Si falla, panic detiene el proceso
	// y muestra el error; así el servidor no arranca creyendo que puede usar la DB.
	if err := db.Connect(); err != nil {
		panic("No se pudo conectar a la base de datos: " + err.Error())
	}

	router := gin.Default()

	// Cada registro conecta un método HTTP y una ruta con su manejador.
	// Por ejemplo, POST /login ejecuta login y GET /verificar-sesion ejecuta
	// verificarSesion. POST suele enviar datos/acciones; GET suele consultar.
	// router.POST("/registro", registro)
	router.POST("/login", login)
	router.POST("/recuperar-password", recuperarPassword)
	router.GET("/verificar-sesion", verificarSesion)
	router.POST("/logout", logout)

	// Inicia el servidor HTTP. Sin argumentos, Gin escucha en el puerto 8080.
	// Esta llamada bloquea aquí mientras el servidor está funcionando.
	router.Run()
}
