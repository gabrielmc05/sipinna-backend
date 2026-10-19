package dboper

import (
	"database/sql"
	"fmt"
	"time"
)

// UsuarioApp representa una cuenta registrada desde la aplicación.
//
// Las cuentas institucionales (Administrador y Alimentador)
// quedan excluidas de esta categoría.
type UsuarioApp struct {
	IDUsuario     uint32    `json:"id_usuario"`
	Correo        string    `json:"correo"`
	Nombre        string    `json:"nombre"`
	Apellidos     string    `json:"apellidos"`
	Telefono      *string   `json:"telefono"`
	FechaRegistro time.Time `json:"fecha_registro"`
	EsActivo      bool      `json:"es_activo"`
}

// ListarUsuariosApp obtiene los usuarios no institucionales.
//
// La autorización del administrador principal se valida
// en el controlador HTTP antes de llamar a esta función.
//
// Retorna la lista de cuentas y un posible error SQL.
func ListarUsuariosApp() ([]UsuarioApp, error) {
	const consulta = `
		SELECT
			u.id_usuario,
			u.correo,
			u.nombre,
			u.apellidos,
			u.telefono,
			u.fecha_registro,
			u.es_activo
		FROM usuario AS u
		WHERE NOT EXISTS (
			SELECT 1
			FROM administrador AS a
			WHERE a.id_usuario = u.id_usuario
		)
		ORDER BY u.fecha_registro DESC, u.id_usuario DESC
	`

	rows, err := db.Query(consulta)
	if err != nil {
		return nil, fmt.Errorf(
			"error al consultar usuarios de la aplicación: %w",
			err,
		)
	}
	defer rows.Close()

	usuarios := make([]UsuarioApp, 0)

	for rows.Next() {
		var usuario UsuarioApp
		var telefono sql.NullString

		err := rows.Scan(
			&usuario.IDUsuario,
			&usuario.Correo,
			&usuario.Nombre,
			&usuario.Apellidos,
			&telefono,
			&usuario.FechaRegistro,
			&usuario.EsActivo,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"error al leer usuario de la aplicación: %w",
				err,
			)
		}

		if telefono.Valid {
			usuario.Telefono = &telefono.String
		}

		usuarios = append(usuarios, usuario)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error al recorrer usuarios: %w",
			err,
		)
	}

	return usuarios, nil
}

// EditarUsuarioApp modifica nombre, apellidos, correo y teléfono.
//
// Solo permite modificar cuentas que no tengan una fila
// en la tabla administrador.
//
// Parámetros:
//   - id: identificador del usuario.
//   - nombre: nombre actualizado.
//   - apellidos: apellidos actualizados.
//   - correo: nueva dirección de correo.
//   - telefono: teléfono opcional.
//
// Retorna true si alguna fila fue modificada.
func EditarUsuarioApp(
	id uint32,
	nombre, apellidos, correo string,
	telefono *string,
) (bool, error) {
	resultado, err := db.Exec(
		`UPDATE usuario
		 SET nombre=?, apellidos=?, correo=?, telefono=?
		 WHERE id_usuario=?
		 AND NOT EXISTS (
			 SELECT 1
			 FROM administrador a
			 WHERE a.id_usuario=usuario.id_usuario
		 )`,
		nombre, apellidos, correo, telefono, id,
	)
	if err != nil {
		return false, err
	}

	n, err := resultado.RowsAffected()
	return n > 0, err
}

// CambiarEstadoUsuarioApp activa o desactiva una cuenta.
//
// Ejecuta ambas operaciones dentro de una transacción:
//   - Actualiza usuario.es_activo.
//   - Invalida la clave de firma de sus JWT.
//
// Bloquea la fila durante la operación para evitar
// modificaciones concurrentes del estado.
//
// Solo afecta usuarios no institucionales.
func CambiarEstadoUsuarioApp(
	id uint32,
	activo bool,
) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var existe uint32

	err = tx.QueryRow(
		`SELECT u.id_usuario
		 FROM usuario u
		 WHERE u.id_usuario=?
		 AND NOT EXISTS (
			 SELECT 1
			 FROM administrador a
			 WHERE a.id_usuario=u.id_usuario
		 )
		 FOR UPDATE`,
		id,
	).Scan(&existe)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if _, err = tx.Exec(
		`UPDATE usuario
		 SET es_activo=?
		 WHERE id_usuario=?`,
		activo, id,
	); err != nil {
		return false, err
	}

	if _, err = tx.Exec(
		`UPDATE token
		 SET llave_firma=NULL
		 WHERE id_usuario=?`,
		id,
	); err != nil {
		return false, err
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
