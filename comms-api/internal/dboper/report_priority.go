package dboper

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetReportPriority consulta la prioridad de un expediente.
//
// Administradores: pueden consultar expedientes de su institución.
// Alimentadores: solo pueden consultar expedientes asignados.
//
// Si no existe una prioridad, devuelve nil.
func GetReportPriority(userID uint32, folio string) (*string, error) {
	var exists bool

	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM expediente WHERE folio = ?)",
		folio,
	).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("verificar expediente: %w", err)
	}
	if !exists {
		return nil, ErrReportNotFound
	}

	var priority sql.NullString

	err = db.QueryRow(`
		SELECT (
			SELECT re.nombre_etiqueta
			FROM reporteetiqueta AS re
			WHERE re.folio = e.folio
			  AND re.tipo = 'Prioridad'
			LIMIT 1
		)
		FROM expediente AS e
		INNER JOIN administrador AS a
		  ON a.id_usuario = ?
		WHERE e.folio = ?
		  AND (
			(a.rol = 'Administrador'
			 AND a.id_institucion = e.id_institucion)
			OR
			(a.rol = 'Alimentador'
			 AND a.id_usuario = e.id_usuario
			 AND a.id_institucion = e.id_institucion)
		  )
	`, userID, folio).Scan(&priority)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrReportAccessDenied
	}
	if err != nil {
		return nil, fmt.Errorf("consultar prioridad: %w", err)
	}

	if !priority.Valid {
		return nil, nil
	}

	return &priority.String, nil
}

// UpdateReportPriority reemplaza la prioridad anterior.
//
// Solo administradores de la institución del expediente
// pueden modificarla.
//
// La transacción garantiza que se elimina la prioridad
// anterior y se registra la nueva como una sola operación.
func UpdateReportPriority(
	userID uint32,
	folio string,
	priority string,
) error {
	if priority != "Alta" &&
		priority != "Media" &&
		priority != "Baja" {
		return fmt.Errorf("prioridad inválida")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	var institutionID uint32

	// Bloquea el expediente durante el cambio de prioridad.
	err = tx.QueryRow(`
		SELECT id_institucion
		FROM expediente
		WHERE folio = ?
		FOR UPDATE
	`, folio).Scan(&institutionID)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrReportNotFound
	}
	if err != nil {
		return fmt.Errorf("consultar expediente: %w", err)
	}

	var authorized bool

	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM administrador
			WHERE id_usuario = ?
			  AND id_institucion = ?
			  AND rol = 'Administrador'
		)
	`, userID, institutionID).Scan(&authorized)

	if err != nil {
		return fmt.Errorf("verificar autorización: %w", err)
	}
	if !authorized {
		return ErrReportAccessDenied
	}

	var active bool

	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM etiqueta
			WHERE nombre_etiqueta = ?
			  AND tipo = 'Prioridad'
			  AND es_activo = 1
		)
	`, priority).Scan(&active)

	if err != nil {
		return fmt.Errorf("validar prioridad: %w", err)
	}
	if !active {
		return fmt.Errorf("prioridad no disponible")
	}

	// Retira exclusivamente las etiquetas de prioridad.
	// No modifica las categorías ni los riesgos.
	_, err = tx.Exec(`
		DELETE FROM reporteetiqueta
		WHERE folio = ?
		  AND tipo = 'Prioridad'
	`, folio)

	if err != nil {
		return fmt.Errorf("eliminar prioridad anterior: %w", err)
	}

	// Registra la nueva prioridad.
	_, err = tx.Exec(`
		INSERT INTO reporteetiqueta (
			folio,
			nombre_etiqueta,
			tipo
		)
		VALUES (?, ?, 'Prioridad')
	`, folio, priority)

	if err != nil {
		return fmt.Errorf("registrar prioridad: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("confirmar prioridad: %w", err)
	}

	return nil
}
