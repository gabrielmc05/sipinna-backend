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

func RequestChallenge(correo string) (*ChallengeResponse, error) {

	resp, err := http.Get(fmt.Sprintf("%s/auth/%s", baseURL(), url.PathEscape(correo)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Se espera que el servicio responda 202 Accepted. Si no, se lee el
	// cuerpo para incluir la explicación recibida en el error.
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
// Si es válida, devuelve el token de acceso (JWT).
func VerifyChallenge(userID uint32, responseBase64 string) (token string, err error) {
	// Marshal convierte valores de Go -> JSON para enviarlos en el cuerpo HTTP.
	payload, err := json.Marshal(struct {
		UserID   uint32 `json:"user_id"`
		Response string `json:"response"`
	}{UserID: userID, Response: responseBase64})
	if err != nil {
		return "", err
	}

	// bytes.NewReader permite que el contenido JSON se envíe como cuerpo.
	resp, err := http.Post(baseURL()+"/auth/", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// En este endpoint, 201 Created indica que la respuesta fue verificada
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("credenciales invalidas: %s", body)
	}

	// Esta estructura local sirve para leer solo el campo token del JSON.
	var parsed struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	// Los nombres de retorno permiten expresar token y err directamente.
	// En este caso la operación terminó bien, así que el error es nil.
	return parsed.Token, nil
}

// Register reenvía los datos de registro; auth-backend es quien cifra
// y guarda la contraseña.
func Register(nombre, apellido, correo, password, telefono string) (userID uint32, err error) {
	// Se arma el JSON usando una estructura anónima: una estructura definida
	// aquí mismo, porque solo hace falta para esta solicitud.
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

	// POST envía los datos de registro. El segundo argumento especifica
	// que el cuerpo de la solicitud está codificado como JSON.
	resp, err := http.Post(baseURL()+"/Registro/", "application/json", bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 201 Created es la respuesta esperada cuando se crea la cuenta.
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("auth-backend rechazo el registro: %s", body)
	}

	// Se leen del JSON de respuesta el ID y el mensaje. El mensaje se decodifica
	// aunque esta función solo necesite devolver el ID.
	var parsed struct {
		UserID  uint32 `json:"user_id"`
		Mensaje string `json:"mensaje"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, err
	}
	// 0 es el valor usado para uint32 cuando ocurre un error.
	return parsed.UserID, nil
}

// ChangePassword reenvía el cambio de contraseña, incluyendo el JWT del
// cliente para que auth-backend verifique que puede hacer este cambio.
func ChangePassword(correo, nuevaContrasena, bearerToken string) error {
	// Convierte la nueva contraseña a JSON para el cuerpo de la solicitud.
	payload, err := json.Marshal(struct {
		NuevaContrasena string `json:"nueva_contraseña"`
	}{nuevaContrasena})
	if err != nil {
		return err
	}

	// Se construye una solicitud PUT manualmente porque, además del cuerpo,
	// hay que agregar el encabezado Authorization con el token del cliente.
	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/Contrasenna/%s", baseURL(), url.PathEscape(correo)),
		bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken)

	// DefaultClient.Do envía la solicitud ya configurada y devuelve la respuesta.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// En un PUT exitoso, este servicio responde 200 OK.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth-backend rechazo el cambio: %s", body)
	}
	// En Go, nil representa la ausencia de error.
	return nil
}
