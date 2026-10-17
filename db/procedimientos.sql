-- DROPS
DROP PROCEDURE IF EXISTS RegistrarAdministrador;
DROP PROCEDURE IF EXISTS RegistrarAlimentador;
DROP PROCEDURE IF EXISTS CrearReporte;
DROP PROCEDURE IF EXISTS AgregarEvidencias;
DROP PROCEDURE IF EXISTS ObtenerReportes;
DROP PROCEDURE IF EXISTS AsignarReporte;
DROP PROCEDURE IF EXISTS ObtenerEvidencias;
DROP PROCEDURE IF EXISTS ObtenerEstados;
DROP PROCEDURE IF EXISTS ObtenerTerritorios;
DROP PROCEDURE IF EXISTS ObtenerTerritorioId;
DROP PROCEDURE IF EXISTS ObtenerAlimentadores;
DROP PROCEDURE IF EXISTS ObtenerPuntosMapa;
DROP PROCEDURE IF EXISTS ConsultarReporte;
DROP PROCEDURE IF EXISTS ConsultarReportePersonal;
DROP PROCEDURE IF EXISTS ListarReportesUsuario;
DROP PROCEDURE IF EXISTS ActualizarEstadoReporte;
DROP FUNCTION IF EXISTS EsAdministrador;
DROP FUNCTION IF EXISTS ExisteUsuario;
DROP PROCEDURE IF EXISTS ConsultarFolio; 
DROP FUNCTION IF EXISTS AsignarEtiquetaNuevoReporte;

-- Procedimientos y funciones
DELIMITER //
CREATE PROCEDURE RegistrarAdministrador (
	p_correo VARCHAR(254),
	p_hash_contrasenna VARBINARY(255),
	p_nombre VARCHAR(30),
	p_apellidos VARCHAR(100),
	p_telefono VARCHAR(30),
	p_id_institucion INT(11)
)
BEGIN

DECLARE p_id_usuario INT(11);
DECLARE EXIT HANDLER FOR SQLEXCEPTION
BEGIN
	ROLLBACK;
	RESIGNAL;
END;	

	IF NOT EXISTS (
		SELECT 1
		FROM institucion AS i
		WHERE i.id_institucion = p_id_institucion)
	THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Institución no encontrada.';
	END IF;
	
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
START TRANSACTION;

	INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
	VALUES (p_correo, UNHEX(p_hash_contrasenna), p_nombre, p_apellidos, p_telefono);
	SET p_id_usuario = LAST_INSERT_ID();
	
	INSERT INTO administrador (id_usuario, id_institucion, rol)
	VALUES (p_id_usuario, p_id_institucion, 'Administrador');

COMMIT;

SELECT id_usuario, nombre_completo FROM usuario WHERE id_usuario = p_id_usuario;
	
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE RegistrarAlimentador (
	p_correo VARCHAR(254),
	p_hash_contrasenna VARBINARY(255),
	p_nombre VARCHAR(30),
	p_apellidos VARCHAR(100),
	p_telefono VARCHAR(30),
	p_id_institucion INT(11)
)
BEGIN

DECLARE p_id_usuario INT(11);
DECLARE EXIT HANDLER FOR SQLEXCEPTION
BEGIN
	ROLLBACK;
	RESIGNAL;
END;	

	IF NOT EXISTS (
		SELECT 1
		FROM institucion AS i
		WHERE i.id_institucion = p_id_institucion)
	THEN
			SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Institución no encontrada.';
	END IF;
	
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
START TRANSACTION;

	INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
	VALUES (p_correo, UNHEX(p_hash_contrasenna), p_nombre, p_apellidos, p_telefono);
	SET p_id_usuario = LAST_INSERT_ID();
	
	INSERT INTO administrador (id_usuario, id_institucion, rol)
	VALUES (p_id_usuario, p_id_institucion, 'Alimentador');

COMMIT;

SELECT id_usuario, nombre_completo FROM usuario WHERE id_usuario = p_id_usuario;
	
END
//
DELIMITER ;

DELIMITER //
CREATE FUNCTION EsAdministrador(f_id_usuario INT(11))
RETURNS BOOL
BEGIN
    RETURN EXISTS ( 
	       SELECT 1 FROM administrador WHERE id_usuario = f_id_usuario);
END; //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE CrearReporte (
	IN p_folio CHAR(15),
	IN p_id_usuario INT,
	IN p_correo_contacto VARCHAR(254),
	IN p_codigo_seguridad_hash CHAR(64),
	IN p_descripcion TEXT,
	IN p_horario_aproximado TIME,
	IN p_es_anonimo BOOLEAN,
	IN p_desea_seguimiento BOOLEAN,
	IN p_hay_adulto_involucrado BOOLEAN,
	IN p_peligro_inmediato VARCHAR(10),
	IN p_latitud DECIMAL(9,6),
	IN p_longitud DECIMAL(9,6),
	IN p_direccion_texto TINYTEXT,
	IN p_id_territorio INT,
	IN p_menores LONGTEXT
)
BEGIN
	DECLARE v_id_institucion INT;
	DECLARE v_indice INT DEFAULT 0;
	DECLARE v_total INT DEFAULT 0;
	DECLARE v_edad TINYINT UNSIGNED;
	DECLARE v_genero VARCHAR(30);
	DECLARE v_descripcion_fisica TINYTEXT;
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		RESIGNAL;
	END;

	IF p_menores IS NULL OR JSON_VALID(p_menores) = 0 OR JSON_TYPE(p_menores) <> 'ARRAY' THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La lista de menores es inválida.';
	END IF;
	IF NOT EXISTS (SELECT 1 FROM territorio WHERE id_territorio = p_id_territorio) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Territorio no encontrado.';
	END IF;
	IF p_peligro_inmediato IS NULL OR p_peligro_inmediato NOT IN ('Sí', 'No', 'No sé') THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Valor de peligro inmediato inválido.';
	END IF;

	SELECT id_institucion INTO v_id_institucion
	FROM institucion
	WHERE id_territorio = p_id_territorio AND es_activo = 1
	ORDER BY id_institucion
	LIMIT 1;
	IF v_id_institucion IS NULL THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No hay una institución activa para el territorio.';
	END IF;

	SET v_total = JSON_LENGTH(p_menores);
	IF v_total < 1 OR v_total > 255 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La cantidad de menores está fuera de rango.';
	END IF;

	START TRANSACTION;
	INSERT INTO reporte (
		folio, id_usuario, correo_contacto, codigo_seguridad_hash, descripcion,
		horario_aproximado, es_anonimo, desea_seguimiento, hay_adulto_involucrado,
		peligro_inmediato
	) VALUES (
		p_folio, p_id_usuario, p_correo_contacto, UNHEX(p_codigo_seguridad_hash), p_descripcion,
		p_horario_aproximado, p_es_anonimo, p_desea_seguimiento, p_hay_adulto_involucrado,
		p_peligro_inmediato
	);
	INSERT INTO ubicacion (folio, latitud, longitud, direccion_texto, id_territorio)
	VALUES (p_folio, p_latitud, p_longitud, p_direccion_texto, p_id_territorio);
	INSERT INTO expediente (folio, id_institucion)
	VALUES (p_folio, v_id_institucion);

	WHILE v_indice < v_total DO
		SET v_edad = NULL;
		SET v_genero = JSON_UNQUOTE(JSON_EXTRACT(p_menores, CONCAT('$[', v_indice, '].genero')));
		SET v_descripcion_fisica = JSON_UNQUOTE(JSON_EXTRACT(p_menores, CONCAT('$[', v_indice, '].descripcion_fisica')));
		IF JSON_TYPE(JSON_EXTRACT(p_menores, CONCAT('$[', v_indice, '].edad_aproximada'))) <> 'NULL' THEN
			SET v_edad = CAST(JSON_UNQUOTE(JSON_EXTRACT(p_menores, CONCAT('$[', v_indice, '].edad_aproximada'))) AS UNSIGNED);
		END IF;
		IF v_genero NOT IN ('Masculino', 'Femenino', 'No identificado') THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Género inválido para un menor.';
		END IF;
		INSERT INTO menor (folio, contador, edad_aproximada, genero, descripcion_fisica)
		VALUES (p_folio, v_indice + 1, v_edad, v_genero, v_descripcion_fisica);
		SET v_indice = v_indice + 1;
	END WHILE;
	COMMIT;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE AgregarEvidencias (
	IN p_folio CHAR(15),
	IN p_evidencias LONGTEXT
)
BEGIN
	DECLARE v_indice INT DEFAULT 0;
	DECLARE v_total INT DEFAULT 0;
	DECLARE v_tipo_archivo VARCHAR(10);
	DECLARE v_ruta VARCHAR(255);
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		RESIGNAL;
	END;

	IF p_evidencias IS NULL OR JSON_VALID(p_evidencias) = 0 OR JSON_TYPE(p_evidencias) <> 'ARRAY' THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'La lista de evidencias es inválida.';
	END IF;
	IF NOT EXISTS (SELECT 1 FROM reporte WHERE folio = p_folio) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Reporte no encontrado.';
	END IF;

	SET v_total = JSON_LENGTH(p_evidencias);
	IF v_total < 1 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Se requiere al menos una evidencia.';
	END IF;
	START TRANSACTION;
	WHILE v_indice < v_total DO
		SET v_tipo_archivo = JSON_UNQUOTE(JSON_EXTRACT(p_evidencias, CONCAT('$[', v_indice, '].tipo_archivo')));
		SET v_ruta = JSON_UNQUOTE(JSON_EXTRACT(p_evidencias, CONCAT('$[', v_indice, '].ruta')));
		IF v_tipo_archivo NOT IN ('Imagen', 'Video') OR v_ruta IS NULL OR CHAR_LENGTH(v_ruta) > 255 THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Datos de evidencia inválidos.';
		END IF;
		INSERT INTO evidencia (folio, tipo_archivo, ruta)
		VALUES (p_folio, v_tipo_archivo, v_ruta);
		SET v_indice = v_indice + 1;
	END WHILE;
	COMMIT;
END
//
DELIMITER ;

DELIMITER // 
CREATE PROCEDURE ConsultarFolio (
	IN p_folio CHAR(15),
	IN p_id_usuario INT(11)
)
BEGIN
	IF NOT EXISTS (SELECT 1 FROM reporte WHERE folio = p_folio)
	THEN
		SIGNAL SQLSTATE '45000'
		SET MESSAGE_TEXT = 'No existen reportes con ese folio';
	END IF;
	
	IF NOT EXISTS (
		SELECT 1 FROM reporte
		WHERE folio = p_folio
		AND id_usuario = p_id_usuario)
	AND NOT EXISTS (
		SELECT 1 FROM administrador 
			WHERE id_usuario = p_id_usuario
			AND rol = 'Administrador')
	AND NOT EXISTS (
		SELECT 1 FROM expediente
			WHERE folio = p_folio
			AND id_usuario = p_id_usuario)
	THEN
		SIGNAL SQLSTATE '45000'
		SET MESSAGE_TEXT = 'Usuario no autorizado';
	END IF;
	
	SELECT r.folio, e.estado, r.fecha_reporte, r.descripcion, COALESCE(u.direccion_texto, ''),
		r.peligro_inmediato
	FROM reporte AS r
	INNER JOIN expediente AS e ON e.folio = r.folio
	INNER JOIN ubicacion AS u ON u.folio = r.folio
	WHERE r.folio = p_folio;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ConsultarReportePersonal (
	IN p_folio CHAR(15),
	IN p_id_usuario INT(11)
)
BEGIN
	IF NOT EXISTS (SELECT 1 FROM reporte WHERE folio = p_folio) THEN
		SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'No existen reportes con ese folio';
	END IF;

	IF NOT EXISTS (
		SELECT 1 FROM reporte
		WHERE folio = p_folio AND id_usuario = p_id_usuario
	) AND NOT EXISTS (
		SELECT 1 FROM administrador
		WHERE id_usuario = p_id_usuario AND rol = 'Administrador'
	) AND NOT EXISTS (
		SELECT 1 FROM expediente
		WHERE folio = p_folio AND id_usuario = p_id_usuario
	) THEN
		SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'Usuario no autorizado';
	END IF;

	SELECT
		r.folio,
		e.estado,
		r.fecha_reporte,
		r.descripcion,
		COALESCE(u.direccion_texto, ''),
		r.peligro_inmediato
	FROM reporte AS r
	INNER JOIN expediente AS e ON e.folio = r.folio
	INNER JOIN ubicacion AS u ON u.folio = r.folio
	WHERE r.folio = p_folio;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ConsultarReporte (
	IN p_folio CHAR(15),
	IN p_codigo_seguridad_hash CHAR(64)
)
BEGIN
	SELECT r.folio, e.estado, r.fecha_reporte, r.descripcion, COALESCE(u.direccion_texto, ''),
		r.peligro_inmediato
	FROM reporte AS r
	INNER JOIN expediente AS e ON e.folio = r.folio
	INNER JOIN ubicacion AS u ON u.folio = r.folio
	WHERE r.folio = p_folio
		AND r.codigo_seguridad_hash = p_codigo_seguridad_hash;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ListarReportesUsuario (
	IN p_id_usuario INT
)
BEGIN
	SELECT r.folio, e.estado, r.fecha_reporte, r.descripcion, COALESCE(u.direccion_texto, ''),
		r.peligro_inmediato
	FROM reporte AS r
	INNER JOIN expediente AS e ON e.folio = r.folio
	INNER JOIN ubicacion AS u ON u.folio = r.folio
	WHERE r.id_usuario = p_id_usuario
	ORDER BY r.fecha_reporte DESC;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ActualizarEstadoReporte (
	IN p_id_usuario INT,
	IN p_folio CHAR(15),
	IN p_estado VARCHAR(40),
	IN p_observaciones TEXT
)
BEGIN
	DECLARE v_reporte_encontrado BOOLEAN DEFAULT FALSE;
	DECLARE v_usuario_autorizado BOOLEAN DEFAULT FALSE;

	IF p_estado NOT IN (
		'Nuevo', 'En revisión', 'Asignado', 'En investigación',
		'Canalizado a Procuraduría', 'Cerrado', 'Descartado'
	) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Estado de expediente inválido.';
	END IF;
	SET v_reporte_encontrado = EXISTS (
		SELECT 1 FROM expediente WHERE folio = p_folio
	);
	SET v_usuario_autorizado = EXISTS (
		SELECT 1
		FROM expediente AS e
		INNER JOIN administrador AS a ON a.id_institucion = e.id_institucion
		WHERE e.folio = p_folio AND a.id_usuario = p_id_usuario
	);
	IF v_reporte_encontrado AND v_usuario_autorizado THEN
		UPDATE expediente
		SET estado = p_estado, observaciones = p_observaciones
		WHERE folio = p_folio;
	END IF;
	SELECT v_reporte_encontrado, v_usuario_autorizado;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerReportes (
	IN p_id_usuario INT,
	IN p_estado VARCHAR(40),
	IN p_id_territorio INT,
	IN p_orden_fecha VARCHAR(4)
)
BEGIN
	DECLARE v_rol VARCHAR(20) DEFAULT NULL;

	SET v_rol = (
		SELECT rol
		FROM administrador
		WHERE id_usuario = p_id_usuario
		LIMIT 1
	);

	IF v_rol IS NULL OR v_rol NOT IN ('Administrador', 'Alimentador') THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Usuario sin rol de personal autorizado.';
	END IF;
	IF p_orden_fecha NOT IN ('asc', 'desc') THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Orden de fecha inválido.';
	END IF;

	SELECT
		r.folio,
		x.estado,
		r.fecha_reporte,
		r.descripcion,
		COALESCE(u.direccion_texto, ''),
		u.id_territorio,
		x.id_usuario,
		r.peligro_inmediato
	FROM reporte AS r
	INNER JOIN expediente AS x ON x.folio = r.folio
	INNER JOIN ubicacion AS u ON u.folio = r.folio
	WHERE (p_estado IS NULL OR x.estado = p_estado)
		AND (p_id_territorio IS NULL OR u.id_territorio = p_id_territorio)
		AND (v_rol = 'Administrador' OR (v_rol = 'Alimentador' AND x.id_usuario = p_id_usuario))
	ORDER BY
		CASE WHEN p_orden_fecha = 'asc' THEN r.fecha_reporte END ASC,
		CASE WHEN p_orden_fecha = 'desc' THEN r.fecha_reporte END DESC,
		r.folio ASC;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE AsignarReporte (
	IN p_id_administrador INT,
	IN p_folio CHAR(15),
	IN p_id_alimentador INT
)
BEGIN
	DECLARE v_reporte_encontrado BOOLEAN DEFAULT FALSE;
	DECLARE v_es_administrador BOOLEAN DEFAULT FALSE;
	DECLARE v_es_alimentador BOOLEAN DEFAULT FALSE;

	SET v_reporte_encontrado = EXISTS (
		SELECT 1 FROM expediente WHERE folio = p_folio
	);
	SET v_es_administrador = EXISTS (
		SELECT 1
		FROM administrador
		WHERE id_usuario = p_id_administrador AND rol = 'Administrador'
	);
	SET v_es_alimentador = EXISTS (
		SELECT 1
		FROM administrador
		WHERE id_usuario = p_id_alimentador AND rol = 'Alimentador'
	);

	IF v_reporte_encontrado AND v_es_administrador AND v_es_alimentador THEN
		UPDATE expediente
		SET id_usuario = p_id_alimentador
		WHERE folio = p_folio;
	END IF;

	SELECT v_reporte_encontrado, v_es_administrador, v_es_alimentador;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerEvidencias (
	IN p_folio CHAR(15),
	IN p_id_usuario INT
)
BEGIN
	DECLARE v_reporte_encontrado BOOLEAN DEFAULT FALSE;
	DECLARE v_usuario_autorizado BOOLEAN DEFAULT FALSE;

	SET v_reporte_encontrado = EXISTS (
		SELECT 1 FROM reporte WHERE folio = p_folio
	);
	SET v_usuario_autorizado = EXISTS (
		SELECT 1
		FROM administrador
		WHERE id_usuario = p_id_usuario
			AND rol = 'Administrador'
	) OR EXISTS (
		SELECT 1
		FROM expediente AS x
		INNER JOIN administrador AS a ON a.id_usuario = x.id_usuario
		WHERE x.folio = p_folio
			AND x.id_usuario = p_id_usuario
			AND a.rol = 'Alimentador'
	);

	SELECT
		v_reporte_encontrado,
		v_usuario_autorizado,
		e.id_evidencia,
		e.tipo_archivo,
		e.ruta
	FROM (SELECT 1) AS solicitud
	LEFT JOIN evidencia AS e ON e.folio = p_folio;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerEstados ()
BEGIN
	SELECT DISTINCT estado
	FROM territorio
	WHERE estado IS NOT NULL
	ORDER BY estado ASC;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerTerritorios (
	IN p_estado VARCHAR(50),
	IN p_municipio VARCHAR(50)
)
BEGIN
	SELECT id_territorio, COALESCE(estado, ''), nombre
	FROM territorio
	WHERE (p_estado IS NULL OR estado = p_estado)
		AND (p_municipio IS NULL OR nombre LIKE CONCAT('%', p_municipio, '%'))
	ORDER BY estado ASC, nombre ASC, id_territorio ASC;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerTerritorioId (
	IN p_id_territorio INT
)
BEGIN
	SELECT id_territorio, COALESCE(estado, ''), nombre
	FROM territorio
	WHERE id_territorio = p_id_territorio;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerAlimentadores ()
BEGIN
	SELECT
		u.id_usuario,
		u.nombre,
		u.apellidos,
		u.correo,
		u.telefono,
		i.id_institucion,
		i.nombre
	FROM administrador AS a
	INNER JOIN usuario AS u ON u.id_usuario = a.id_usuario
	INNER JOIN institucion AS i ON i.id_institucion = a.id_institucion
	WHERE a.rol = 'Alimentador'
		AND u.es_activo = 1
		AND i.es_activo = 1
	ORDER BY u.apellidos ASC, u.nombre ASC, u.id_usuario ASC;
END
//
DELIMITER ;

DELIMITER //
CREATE PROCEDURE ObtenerPuntosMapa (
	IN p_id_usuario INT
)
BEGIN
	DECLARE v_rol VARCHAR(20) DEFAULT NULL;

	SET v_rol = (
		SELECT rol
		FROM administrador
		WHERE id_usuario = p_id_usuario
		LIMIT 1
	);

	IF v_rol IS NULL OR v_rol NOT IN ('Administrador', 'Alimentador') THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Usuario sin rol de personal autorizado.';
	END IF;

	SELECT
		r.folio,
		u.latitud,
		u.longitud,
		x.estado,
		r.fecha_reporte,
		t.id_territorio,
		COALESCE(t.estado, ''),
		t.nombre,
		r.peligro_inmediato
	FROM reporte AS r
	INNER JOIN expediente AS x ON x.folio = r.folio
	INNER JOIN ubicacion AS u ON u.folio = r.folio
	INNER JOIN territorio AS t ON t.id_territorio = u.id_territorio
	WHERE v_rol = 'Administrador'
		OR (v_rol = 'Alimentador' AND x.id_usuario = p_id_usuario)
	ORDER BY r.fecha_reporte DESC, r.folio ASC;
END
//
DELIMITER ;

DELIMITER //
CREATE FUNCTION ExisteUsuario(f_id_usuario INT(11))
RETURNS BOOL
BEGIN
    RETURN EXISTS ( 
		SELECT 1 FROM usuario WHERE id_usuario = f_id_usuario);
END; //
DELIMITER ;

-- Este solo se usará al asignar etiquetas a nuevos reportes
DELIMITER //
CREATE PROCEDURE AsignarEtiquetaNuevoReporte(
	IN p_folio CHAR(15),
	IN p_nombre_etiqueta VARCHAR(50),
	IN p_tipo VARCHAR(20)
)
BEGIN
	DECLARE v_etiqueta_encontrada BOOLEAN DEFAULT FALSE;
	DECLARE v_reporte_encontrado BOOLEAN DEFAULT FALSE;

	IF p_tipo NOT IN ('Categoria','Prioridad','Riesgo','Observacion'
	) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Tipo de etiqueta inválido.';
	END IF;
	
	SET v_etiqueta_encontrada = EXISTS (
		SELECT 1 FROM etiqueta
			WHERE nombre_etiqueta = p_nombre_etiqueta
			AND tipo = p_tipo
			AND es_activo = TRUE
	);
	
	SET v_reporte_encontrado = EXISTS (
		SELECT 1 FROM reporte WHERE folio = p_folio 
	);
	
	IF v_etiqueta_encontrada AND v_reporte_encontrado THEN
		INSERT INTO reporteetiqueta (folio, nombre_etiqueta, tipo)
		VALUES (p_folio, p_nombre_etiqueta, p_tipo);
	END IF;
	
	SELECT folio, nombre_etiqueta, tipo FROM reporteetiqueta; 
END //
DELIMITER ;
