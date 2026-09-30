package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	
	
) 
// Referencias Bibliograficas
// https://gin-gonic.com/en/docs/routing/api-design/


type LoginInput struct {
	Correo string `json:"correo"`
	Password string `json:"password"`
}

type RegistroInput struct {
	Nombre string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Correo string `json:"correo"`
	Password string `json:"password"`
	
}

type RecuperarPasswordInput struct{
	Correo string `json:"correo"`
}


// Endpoints de autenticacion

// Estructuras de datos 
func main(){
	// Correcion del codigo 
	router := gin.Default()
	
	// Funcion de login donde procesa un iunicio de sesion con contrasena y correo 
	
	func login(c *gin.Context ){
		var input LoginInput
		if err:= c.ShouldBindJSON(&input); err != nill{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			// Regresa un mensaje
			// Regresa el token
		})

		
	}
	
	// Este es para el registro 
	
	func registro(c *gin.Context){
		var input RegistroInput
		if err:= c.ShouldBindJSON(&input); err != nill{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro invalidos"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			// mensaja de usuario creado de manera exitosa 
		})				
	}
	
	
	func recuperarPassword(c *gin.Context){
		var input RecuperarPasswordInput
		if err:= c.ShouldBindJSON(&input); err != nill{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un correo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			// mensaje de correo de recuperacion 
		})
		
		
	}
	
	// Se me ocurre que es para la carga inicial de la App
	fun verificarSesion(c *gin.Context){
		c.JSON(http.StatusOK, gin.H){
			"valido":true, 
			"rol": "Ciudadano"
		}
	}

	
	
	func  logout(c *gin.Context){
		c.JSON(http.StatusOK, gin.H){
			"mensaje": "Sesion cerrada de manera exitosa"
		}
	}
}