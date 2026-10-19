package dboper

import (
	"database/sql"
	"errors"
	"fmt"
)

// Errores de gestión de cuentas administradoras.
var (
	// ErrAdministradorNoEncontrado indica que la cuenta
	// no existe o no tiene rol Administrador.
	ErrAdministradorNoEncontrado = errors.New(
		"administrador no encontrado",
	)

	// ErrAdministradorProtegido impide desactivar
	// la cuenta administradora principal ID 3.
	ErrAdministradorProtegido = errors.New(
		"no se puede desactivar al administrador principal",
	)
)

// AdministradorCuenta representa una cuenta institucional
// cuyo rol en MariaDB es Administrador.
type AdministradorCuenta struct {
	ID                uint32  `json:"id_usuario"`
	Name              string  `json:"nombre"`
	Surname           string  `json:"apellidos"`
	Email             string  `json:"correo"`
	Phone             *string `json:"telefono"`
	InstitutionID     uint32  `json:"id_institucion"`
	Institution       string  `json:"institucion"`
	Active            bool    `json:"es_activo"`
	InstitutionActive bool    `json:"institucion_activa"`
}

// ListarAdministradores obtiene administradores activos
// e inactivos directamente desde MariaDB.
//
// Únicamente devuelve cuentas con rol Administrador.
// La autorización HTTP se comprueba en el handler.
func ListarAdministradores() ([]AdministradorCuenta, error) {
	const query = `
		SELECT
			u.id_usuario,
			u.nombre,
			u.apellidos,
			u.correo,
			u.telefono,
			i.id_institucion,
			i.nombre,
			u.es_activo,
			i.es_activo
		FROM administrador AS a
		INNER JOIN usuario AS u
			ON u.id_usuario = a.id_usuario
		INNER JOIN institucion AS i
			ON i.id_institucion = a.id_institucion
		WHERE a.rol = 'Administrador'
		ORDER BY
			u.apellidos ASC,
			u.nombre ASC,
			u.id_usuario ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(
			"listar administradores: %w",
			err,
		)
	}
	defer rows.Close()

	administradores := make([]AdministradorCuenta, 0)

	for rows.Next() {
		var administrador AdministradorCuenta

		err := rows.Scan(
			&administrador.ID,
			&administrador.Name,
			&administrador.Surname,
			&administrador.Email,
			&administrador.Phone,
			&administrador.InstitutionID,
			&administrador.Institution,
			&administrador.Active,
			&administrador.InstitutionActive,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"leer administrador: %w",
				err,
			)
		}

		administradores = append(
			administradores,
			administrador,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"recorrer administradores: %w",
			err,
		)
	}

	return administradores, nil
}

// EditarAdministradorActual actualiza la información
// personal y la institución de un administrador.
//
// Ambas operaciones se ejecutan en una transacción.
// Si alguna falla, se revierten todos los cambios.
//
// No permite modificar cuentas de tipo Alimentador.
func EditarAdministradorActual(
	id uint32,
	datos DatosAlimentadorEditado,
) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf(
			"iniciar edición: %w",
			err,
		)
	}
	defer tx.Rollback()

	var encontrado uint32

	err = tx.QueryRow(
		`SELECT u.id_usuario
		 FROM usuario AS u
		 INNER JOIN administrador AS a
			ON a.id_usuario = u.id_usuario
		 WHERE u.id_usuario = ?
		   AND a.rol = 'Administrador'
		 FOR UPDATE`,
		id,
	).Scan(&encontrado)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAdministradorNoEncontrado
	}

	if err != nil {
		return fmt.Errorf(
			"buscar administrador: %w",
			err,
		)
	}

	// Verifica que la institución exista y esté activa.
	var institucionID uint32

	err = tx.QueryRow(
		`SELECT id_institucion
		 FROM institucion
		 WHERE id_institucion = ?
		   AND es_activo = 1`,
		datos.InstitucionID,
	).Scan(&institucionID)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrInstitucionNoDisponible
	}

	if err != nil {
		return fmt.Errorf(
			"consultar institución: %w",
			err,
		)
	}

	// Actualiza los datos personales del usuario.
	_, err = tx.Exec(
		`UPDATE usuario
		 SET nombre = ?,
		     apellidos = ?,
		     correo = ?,
		     telefono = ?
		 WHERE id_usuario = ?`,
		datos.Nombre,
		datos.Apellidos,
		datos.Correo,
		datos.Telefono,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"editar datos personales: %w",
			err,
		)
	}

	// Actualiza la institución de la cuenta.
	_, err = tx.Exec(
		`UPDATE administrador
		 SET id_institucion = ?
		 WHERE id_usuario = ?
		   AND rol = 'Administrador'`,
		datos.InstitucionID,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"editar institución: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"confirmar edición: %w",
			err,
		)
	}

	return nil
}

// CambiarEstadoAdministradorActual activa o desactiva
// una cuenta administradora.
//
// Verifica que la cuenta tenga rol Administrador,
// modifica es_activo y revoca su clave de firma JWT.
//
// Protege la cuenta administradora principal ID 3
// contra la desactivación accidental.
func CambiarEstadoAdministradorActual(
	id uint32,
	activo bool,
) error {
	if id == 3 && !activo {
		return ErrAdministradorProtegido
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf(
			"iniciar cambio de estado: %w",
			err,
		)
	}
	defer tx.Rollback()

	var encontrado uint32

	err = tx.QueryRow(
		`SELECT u.id_usuario
		 FROM usuario AS u
		 INNER JOIN administrador AS a
			ON a.id_usuario = u.id_usuario
		 WHERE u.id_usuario = ?
		   AND a.rol = 'Administrador'
		 FOR UPDATE`,
		id,
	).Scan(&encontrado)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAdministradorNoEncontrado
	}

	if err != nil {
		return fmt.Errorf(
			"consultar administrador: %w",
			err,
		)
	}

	// Actualiza el estado de la cuenta.
	_, err = tx.Exec(
		`UPDATE usuario
		 SET es_activo = ?
		 WHERE id_usuario = ?`,
		activo,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"actualizar estado: %w",
			err,
		)
	}

	// Invalida los JWT previamente emitidos.
	_, err = tx.Exec(
		`UPDATE token
		 SET llave_firma = NULL
		 WHERE id_usuario = ?`,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"revocar sesión: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"confirmar estado: %w",
			err,
		)
	}

	return nil
}
