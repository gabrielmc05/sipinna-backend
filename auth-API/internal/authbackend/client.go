package authbackend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func baseURL() string {
	u := os.Getenv("AUTH_BACKEND_URL")
	if u == "" {
		u = "http://localhost:8081"
	}
	return u
}

// para que auth-API no se quede esperando para siempre si auth-backend no repsonde
var httpClient = &http.Client{Timeout: 10 * time.Second}

// guardar HTTP que devolvió auth-backend. Así los handlers deciden qué contestarle al cliente sin reenviar mensajes internos.
type BackendError struct {
	Status int
}
func (e *BackendError) Error() string {
	return fmt.Sprintf("auth-backend respondio con codigo %d", e.Status)
}

// doJSON hace una petición a auth-backend. body puede ser nil. authHeader es el
// valor del header Authorization tal como lo mandó el cliente (puede ir vacío).
// Si el código HTTP no es 2xx, regresa un *BackendError.
func doJSON(method, path string, body any, authHeader string) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(payload)
	}
 
	req, err := http.NewRequest(method, baseURL()+path, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
 
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
 
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, &BackendError{Status: resp.StatusCode}
	}
	return respBody, nil
}

type ChallengeResponse struct {
	UserID    uint32 `json:"user_id"`
	Challenge string `json:"challenge"`
	Mensaje   string `json:"mensaje"`
}

// RequestChallenge envía al servicio el correo del usuario y recibe un challenge
// que el cliente debe resolver para demostrar que conoce la contraseña.
func RequestChallenge(correo string) (*ChallengeResponse, error) {
 
	resp, err := httpClient.Get(fmt.Sprintf("%s/auth/%s", baseURL(), url.PathEscape(correo)))
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
// Si es válida devuelve el JWT.
func VerifyChallenge(userID uint32, responseBase64 string) (token string, err error) {
	// Marshal convierte valores de Go -> JSON para enviarlos en el cuerpo HTTP
	payload, err := json.Marshal(struct {
		UserID   uint32 `json:"user_id"`
		Response string `json:"response"`
	}{UserID: userID, Response: responseBase64})
	if err != nil {
		return "", err
	}
 
	// bytes.NewReader para enviar JSON como cuerpo.
	// La ruta es "/auth" sin diagonal final, igual que la registró auth-backend.
	resp, err := httpClient.Post(baseURL()+"/auth", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
 
	// 201 Created -> respuesta fue verificada
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("credenciales invalidas: %s", body)
	}
 
	// auth-backend regresa el JWT como texto plano (no como JSON), así que se lee tal cual.
	tokenBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

// Register reenvía los datos de registro, auth-backend es quien cifra y guarda la contraseña
func Register(nombre, apellido, correo, password string, telefono *string) (userID uint32, err error) {
 
	payload, err := json.Marshal(struct {
		Nombre   string  `json:"name"`
		Apellido string  `json:"surname"`
		Correo   string  `json:"email"`
		Password string  `json:"password"`
		Telefono *string `json:"phone_number,omitempty"`
	}{nombre, apellido, correo, password, telefono})
	if err != nil {
		return 0, err
	}
 
	// POST envía los datos de registro
	// El segundo argumento especifica que el cuerpo de la solicitud está como JSON
	resp, err := httpClient.Post(baseURL()+"/register", "application/json", bytes.NewReader(payload))
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
		Mensaje string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, err
	}
 
	return parsed.UserID, nil
}

// ChangePassword reenvía el cambio de contraseña incluyendo el JWT del cliente para que auth-backend verifique que puede hacer el cambio
func ChangePassword(correo, nuevaContrasena, authHeader string) error {
	_, err := doJSON(http.MethodPatch, "/password/"+url.PathEscape(correo),
		struct {
			NewPassword string `json:"new_password"`
		}{nuevaContrasena}, authHeader)
	return err
}

// Logout: auth-backend borra la llave de firma del usuario, y con eso sus tokens dejan de servir.
func Logout(userID uint32, authHeader string) error {
	_, err := doJSON(http.MethodDelete, fmt.Sprintf("/auth/%d", userID), nil, authHeader)
	return err
}
 
type UserData struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber any    `json:"phone_number"` // any: pasa tal cual lo que mande auth-backend (texto o null)
	UserID      uint32 `json:"user_id"`
}
 
// GetUser pide los datos del usuario con ese correo. Requiere el JWT del propio usuario.
func GetUser(correo, authHeader string) (*UserData, error) {
	body, err := doJSON(http.MethodGet, "/user/"+url.PathEscape(correo), nil, authHeader)
	if err != nil {
		return nil, err
	}
	var data UserData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
 
// Los punteros + omitempty hacen que los campos que no mandes no viajen en el JSON.
type EditUserFields struct {
	Name        *string `json:"name,omitempty"`
	Surname     *string `json:"surname,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}
 
func EditUser(correo string, fields EditUserFields, authHeader string) error {
	_, err := doJSON(http.MethodPatch, "/user/"+url.PathEscape(correo), fields, authHeader)
	return err
}
 
func ChangeEmail(userID uint32, newEmail, authHeader string) error {
	_, err := doJSON(http.MethodPatch, "/email", struct {
		UserID   uint32 `json:"user_id"`
		NewEmail string `json:"new_email"`
	}{userID, newEmail}, authHeader)
	return err
}