package main

import (
	"net/http"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/auth-API/internal/auth"
	"github.com/gabrielmc05/sipinna-backend/auth-API/internal/db"
	"github.com/gin-gonic/gin"
)

// Referencias Bibliograficas
// https://gin-gonic.com/en/docs/routing/api-design/

// Estas estructuras describen los datos que la API espera recibir en el cuerpo
// (body) de las solicitudes. Los textos dentro de `json:"..."` son los nombres
// de las propiedades que deben venir en el JSON enviado por el cliente.
type LoginInput struct {
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

// Un endpoint es una ruta HTTP que permite al cliente pedir una acción al servidor.
// En Gin, la función que recibe *gin.Context puede leer la solicitud y escribir la respuesta.

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

func registro(c *gin.Context) {
	var input RegistroInput
	// ShouldBindJSON lee el JSON del body e intenta convertirlo a RegistroInput.
	// `&input` pasa la dirección de la variable para que Gin pueda llenarla.
	if err := c.ShouldBindJSON(&input); err != nil {
		// Si el JSON falta o no tiene el formato esperado, respondemos con HTTP 400.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro invalidos"})
		return
	}

	// Nunca guardamos la contraseña original: HashPassword genera un hash seguro
	// que sí se puede guardar. El error se comprueba antes de continuar.
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo procesar la contraseña"})
		return
	}

	// DB.Exec ejecuta una instrucción SQL que no devuelve filas. Los signos ?
	// son parámetros: el driver inserta los valores por separado para evitar
	// construir la consulta concatenando texto recibido del usuario.
	_, err = db.DB.Exec(
		"INSERT INTO usuario (correo, contrasenna_encriptada, nombre, apellidos) VALUES (?, ?, ?, ?)",
		input.Correo, hashedPassword, input.Nombre, input.Apellidos,
	)
	if err != nil {
		// MariaDB/MySQL suele indicar así que el correo ya existe (por una restricción
		// UNIQUE). Otros errores de base de datos reciben una respuesta genérica.
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "Ese correo ya esta registrado, por favor intenta con otro"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	// 201 Created significa que el servidor creó correctamente el recurso.
	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Usuario creado de manera exitosa",
	})
}

// login recibe y valida el formato del JSON de inicio de sesión.
// Por ahora solo responde "Login exitoso": todavía no busca al usuario,
// no verifica su contraseña y tampoco genera un token.
func login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// 200 OK indica que la solicitud fue atendida. En este punto es una respuesta
	// provisional, porque aún no se comprobó que las credenciales sean correctas.
	c.JSON(http.StatusOK, gin.H{
		// Regresa un mensaje
		"mensaje": "Login exitoso",
		// Regresa el token
	})
}

func recuperarPassword(c *gin.Context) {
	var input RecuperarPasswordInput
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

	// Cada línea registra un método HTTP y su ruta, y la asocia con la función
	// que atenderá esa solicitud. POST suele usarse para enviar datos o ejecutar
	// acciones; GET suele usarse para consultar información.
	router.POST("/registro", registro)
	router.POST("/login", login)
	router.POST("/recuperar-password", recuperarPassword)
	router.GET("/verificar-sesion", verificarSesion)
	router.POST("/logout", logout)

	// Inicia el servidor HTTP. Sin argumentos, Gin escucha en el puerto 8080.
	// Esta llamada bloquea aquí mientras el servidor está funcionando.
	router.Run()
}
