-- DROPS
DROP PROCEDURE IF EXISTS RegistrarAdministrador;
DROP PROCEDURE IF EXISTS RegistrarAlimentador;
DROP FUNCTION IF EXISTS EsAdministrador;
DROP FUNCTION IF EXISTS ExisteUsuario;

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
	VALUES (p_correo, p_hash_contrasenna, p_nombre, p_apellidos, p_telefono);
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
	VALUES (p_correo, p_hash_contrasenna, p_nombre, p_apellidos, p_telefono);
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
CREATE FUNCTION ExisteUsuario(f_id_usuario INT(11))
RETURNS BOOL
BEGIN
    RETURN EXISTS ( 
		SELECT 1 FROM usuario WHERE id_usuario = f_id_usuario);
END; //
DELIMITER ;
