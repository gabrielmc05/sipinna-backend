package authbackend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func baseURL() string {
	u := os.Getenv("AUTH_BACKEND_URL")
	if u == "" {
		u = "http://localhost:8081"
	}
	return u
}

type ChallengeResponse struct {
	UserID    uint32 `json:"user_id"`
	Challenge string `json:"challenge"`
	Mensaje   string `json:"mensaje"`
}

// RequestChallenge envía al servicio el correo del usuario y recibe un challenge
// que el cliente debe resolver para demostrar que conoce la contraseña.
func RequestChallenge(correo string) (*ChallengeResponse, error) {

	resp, err := http.Get(fmt.Sprintf("%s/auth/%s", baseURL(), url.PathEscape(correo)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// 202 Accepted, si no, se lee el cuerpo para incluir la explicación recibida en el error.
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth-backend respondio %d: %s", resp.StatusCode, body)
	}

	// El decoder transforma JSON de respuesta -> estructura de Go.
	var parsed ChallengeResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

// VerifyChallenge envía al servicio la respuesta que ya calculó el cliente.
// Si es válida devuelve JWT
func VerifyChallenge(userID uint32, responseBase64 string) (token string, err error) {
	// Marshal convierte valores de Go -> JSON para enviarlos en el cuerpo HTTP
	payload, err := json.Marshal(struct {
		UserID   uint32 `json:"user_id"`
		Response string `json:"response"`
	}{UserID: userID, Response: responseBase64})
	if err != nil {
		return "", err
	}

	// bytes.NewReader para enviar JSON como cuerpo
	resp, err := http.Post(baseURL()+"/auth/", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 201 Created -> respuesta fue verificada
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("credenciales invalidas: %s", body)
	}

	// var parsed struct {
	// 	Token string `json:"token"`
	// }
	// // Decodr JSON -> Go
	// if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
	// 	return "", err
	// }
	tokenbytes, _:= io.ReadAll(resp.Body)
	token = string(tokenbytes)	
	return token, nil
}

// Register reenvía los datos de registro, auth-backend es quien cifra y guarda la contraseña
func Register(nombre, apellido, correo, password, telefono string) (userID uint32, err error) {

	payload, err := json.Marshal(struct {
		Nombre   string `json:"nombre"`
		Apellido string `json:"apellido"`
		Correo   string `json:"correo"`
		Password string `json:"contraseña"`
		Telefono string `json:"numero_de_telefono"`
	}{nombre, apellido, correo, password, telefono})
	if err != nil {
		return 0, err
	}

	// POST envía los datos de registro
	// El segundo argumento especifica que el cuerpo de la solicitud está como JSON
	resp, err := http.Post(baseURL()+"/Registro/", "application/json", bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("auth-backend rechazo el registro :( %s", body)
	}

	// Se leen del JSON de respuesta el ID y el mensaje. El mensaje se decodifica aunque esta función solo necesite devolver el ID.
	var parsed struct {
		UserID  uint32 `json:"user_id"`
		Mensaje string `json:"mensaje"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, err
	}

	return parsed.UserID, nil
}

// ChangePassword reenvía el cambio de contraseña incluyendo el JWT del cliente para que auth-backend verifique que puede hacer el cambio
func ChangePassword(correo, nuevaContrasena, bearerToken string) error {

	payload, err := json.Marshal(struct {
		NuevaContrasena string `json:"nueva_contraseña"`
	}{nuevaContrasena})
	if err != nil {
		return err
	}

	// Además del cuerpo hay que agregar el encabezado Authorization con el token del cliente
	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/Contrasenna/%s", baseURL(), url.PathEscape(correo)),
		bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken)

	resp, err := http.DefaultClient.Do(req) //envía la solicitud ya configurada y devuelve la respuesta
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth-backend rechazo el cambio: %s", body)
	}

	return nil
}
