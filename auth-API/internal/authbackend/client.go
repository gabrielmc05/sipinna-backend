package authbackend

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)


func baseURL() string {
	url := os.Getenv("AUTH_BACKEND_URL")
	if url == "" {
		url = "http://localhost:8081"
	}
	return url
}

// Le pedimos a auth-backend un nuevo reto para el usuario indicado.
func RequestChallenge(idUsuario uint32) ([]byte, error) { 
	// Hace un GET (consulta) a /auth/{id}
	// resp contiene el código
	// HTTP, los encabezados y el cuerpo
	resp, err := http.Get(fmt.Sprintf("%s/auth/%d", baseURL(), idUsuario)) 
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// en caso de error
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("auth-backend respondio %d al pedir el challenge", resp.StatusCode)
	} 

	// Deberpa ser llamado con la clave "challenge".
	var parsed struct {
		Challenge string `json:"challenge"`
	}
	// Decodifica la respuesta JSON y extrae el reto de auth-backend.
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	// base64 -> secuencia de bytes
	return base64.StdEncoding.DecodeString(parsed.Challenge) 
} 

// En resumen, RequestChallenge hace lo siguiente:
// 1. Hace un GET a /auth/{idUsuario} en auth-backend.
// 2. Recibe un JSON con un reto aleatorio el cual es generado por auth-backend.
// 3. Decodifica el reto de base64 a []byte y lo devuelve.


// VerifyChallenge calcula la respuesta CHAP con la contraseña real del usuario
// y se la manda a auth-backend para que la verifique.
func VerifyChallenge(idUsuario uint32, password string, challenge []byte) (token string, err error) {
	// []byte(password) convierte el texto a bytes, que es lo que requiere SHA-256.
	// Sum devuelve un arreglo de 32 bytes con el hash de la contraseña.
	passwordHash := sha256.Sum256([]byte(password))
	// CHAP: se calcula un hash usando el hash de la contraseña y el reto recibido.
	response := sha256.Sum256(append(passwordHash[:], challenge...))

	// json.Marshal convierte el valor Go en bytes JSON para enviarlo por HTTP.
	// Los nombres json: definen cómo se llamarán las propiedades del objeto enviado.
	payload, err := json.Marshal(struct {
		UserID   uint32 `json:"user_id"`
		Response string `json:"response"`
	}{
		UserID:   idUsuario,
		Response: base64.StdEncoding.EncodeToString(response[:]),
	})
	if err != nil {
		return "", err
	}

	// POST envía el JSON al servicio. bytes.NewReader permite usar esos bytes
	// como cuerpo de la solicitud; el encabezado informa que el cuerpo es JSON.
	resp, err := http.Post(baseURL()+"/auth/", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Este endpoint debe responder 201 Created cuando la verificación fue exitosa.
	if resp.StatusCode != http.StatusCreated {
		// En caso de error se lee el cuerpo para incluir el detalle en el error.
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("credenciales invalidas: %s", respBody)
	}

	// Decodifica el JSON de respuesta y extrae el token emitido por auth-backend.
	var parsed struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	// Esta función tiene retornos con nombre; aquí se devuelve el token y nil,
	// que en Go indica que no ocurrió ningún error.
	return parsed.Token, nil
}