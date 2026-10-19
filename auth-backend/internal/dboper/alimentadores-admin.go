package dboper

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrAlimentadorNoEncontrado indica que la cuenta solicitada
// no existe o no tiene el rol Alimentador.
var ErrAlimentadorNoEncontrado = errors.New(
	"alimentador no encontrado",
)

// ErrInstitucionNoDisponible indica que la institución de destino
// no existe o está desactivada.
var ErrInstitucionNoDisponible = errors.New(
	"institución no disponible",
)

// DatosAlimentadorEditado contiene los datos editables
// de una cuenta institucional de tipo Alimentador.
type DatosAlimentadorEditado struct {
	Nombre        string
	Apellidos     string
	Correo        string
	Telefono      *string
	InstitucionID uint32
}

// EditarAlimentador actualiza los datos personales y la institución
// en una sola transacción.
//
// Solo permite editar cuentas con rol Alimentador.
// La institución de destino debe existir y estar activa.
//
// No modifica la contraseña ni el rol institucional.
func EditarAlimentador(
	usuarioID uint32,
	datos DatosAlimentadorEditado,
) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar edición: %w", err)
	}
	defer tx.Rollback()

	var encontrado uint32

	err = tx.QueryRow(
		`SELECT u.id_usuario
		 FROM usuario u
		 INNER JOIN administrador a
		   ON a.id_usuario = u.id_usuario
		 WHERE u.id_usuario = ?
		   AND a.rol = 'Alimentador'
		 FOR UPDATE`,
		usuarioID,
	).Scan(&encontrado)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAlimentadorNoEncontrado
	}
	if err != nil {
		return fmt.Errorf("consultar alimentador: %w", err)
	}

	var institucionDisponible uint32
	err = tx.QueryRow(
		`SELECT id_institucion
		 FROM institucion
		 WHERE id_institucion = ?
		   AND es_activo = 1`,
		datos.InstitucionID,
	).Scan(&institucionDisponible)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrInstitucionNoDisponible
	}
	if err != nil {
		return fmt.Errorf("consultar institución: %w", err)
	}

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
		usuarioID,
	)
	if err != nil {
		return fmt.Errorf("actualizar datos personales: %w", err)
	}

	_, err = tx.Exec(
		`UPDATE administrador
		 SET id_institucion = ?
		 WHERE id_usuario = ?
		   AND rol = 'Alimentador'`,
		datos.InstitucionID,
		usuarioID,
	)
	if err != nil {
		return fmt.Errorf("actualizar institución: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("confirmar edición: %w", err)
	}

	return nil
}

// CambiarEstadoAlimentador activa o desactiva una cuenta.
//
// Verifica que la cuenta tenga rol Alimentador.
// La actualización del estado y la revocación de JWT
// se ejecutan en una misma transacción.
//
// Al reactivar, la clave JWT queda invalidada igualmente;
// el alimentador deberá iniciar sesión otra vez.
func CambiarEstadoAlimentador(
	usuarioID uint32,
	activo bool,
) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar cambio de estado: %w", err)
	}
	defer tx.Rollback()

	var encontrado uint32

	err = tx.QueryRow(
		`SELECT u.id_usuario
		 FROM usuario u
		 INNER JOIN administrador a
		   ON a.id_usuario = u.id_usuario
		 WHERE u.id_usuario = ?
		   AND a.rol = 'Alimentador'
		 FOR UPDATE`,
		usuarioID,
	).Scan(&encontrado)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAlimentadorNoEncontrado
	}
	if err != nil {
		return fmt.Errorf("consultar alimentador: %w", err)
	}

	_, err = tx.Exec(
		`UPDATE usuario
		 SET es_activo = ?
		 WHERE id_usuario = ?`,
		activo,
		usuarioID,
	)
	if err != nil {
		return fmt.Errorf("actualizar estado: %w", err)
	}

	_, err = tx.Exec(
		`UPDATE token
		 SET llave_firma = NULL
		 WHERE id_usuario = ?`,
		usuarioID,
	)
	if err != nil {
		return fmt.Errorf("revocar JWT: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("confirmar estado: %w", err)
	}

	return nil
}
