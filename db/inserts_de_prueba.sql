-- Territorio
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Atizapán de Zaragoza');
SET @territorio_id = LAST_INSERT_ID();

-- Institución
INSERT INTO institucion (id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES (@territorio_id, 'SIPINNA Atizapán de Zaragoza', 'SIPINNA',
        'contacto@sipinna-atizapan.gob.mx', '5588223399', 'Av. Ejemplo 12345');
SET @institucion_id = LAST_INSERT_ID();

-- Instituciones adicionales para pruebas
INSERT INTO institucion (id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES (@territorio_id, 'Procuraduría de Protección de Atizapán', 'Procuraduría',
        'contacto.procuraduria.prueba@example.com', '5552223344', 'Calle Ficticia 200');
SET @institucion_procuraduria_id = LAST_INSERT_ID();

INSERT INTO institucion (id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES (@territorio_id, 'DIF Municipal de Atizapán', 'DIF',
        'contacto.dif.prueba@example.com', '5553334455', 'Av. Inventada 300');
SET @institucion_dif_id = LAST_INSERT_ID();

-- Usuario que será Administrador
INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
VALUES ('admin.prueba@sipinna.gob.mx',
        CONVERT('loooooool' USING BINARY),
        'Ana', 'Pérez López', '5551112222');
SET @usuario_admin_id = LAST_INSERT_ID();

INSERT INTO administrador (id_usuario, id_institucion, rol)
VALUES (@usuario_admin_id, @institucion_id, 'Administrador');

-- Alimentadores ficticios vinculados a las instituciones de prueba
INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
VALUES ('alimentador.procuraduria.prueba@example.com',
        CONVERT('holahola' USING BINARY),
        'María', 'López García', '5554445566');
SET @usuario_alimentador_procuraduria_id = LAST_INSERT_ID();

INSERT INTO administrador (id_usuario, id_institucion, rol)
VALUES (@usuario_alimentador_procuraduria_id, @institucion_procuraduria_id, 'Alimentador');

INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
VALUES ('alimentador.dif.prueba@example.com',
        CONVERT('adiosadios' USING BINARY),
        'Carlos', 'Hernández Ruiz', '5556667788');
SET @usuario_alimentador_dif_id = LAST_INSERT_ID();

INSERT INTO administrador (id_usuario, id_institucion, rol)
VALUES (@usuario_alimentador_dif_id, @institucion_dif_id, 'Alimentador');

-- Usuario que será el Reportante
INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
VALUES ('reportante.prueba@example.com',
        CONVERT('Papupro38' USING BINARY),
        'Juan', 'Ramírez Soto', NULL);
SET @usuario_reportante_id = LAST_INSERT_ID();

-- Reporte
INSERT INTO reporte (folio, id_usuario, correo_contacto, codigo_seguridad_hash,
                      descripcion, horario_aproximado, es_anonimo,
                      desea_seguimiento, hay_adulto_involucrado)
VALUES ('RPT0000001', @usuario_reportante_id, NULL,
        '$argon2id$v=19$m=65536,t=3,p=4$placeholder_hash_codigo',
        'Se observó a un menor vendiendo dulces en un cruce vehicular durante horario escolar.',
        '10:30:00', FALSE, TRUE, TRUE);

-- Ubicación
INSERT INTO ubicacion (folio, latitud, longitud, direccion_texto, id_territorio)
VALUES ('RPT0000001', 19.573800, -99.270500,
        'Cerca de la estación de camiones, Av. Central', @territorio_id);

-- Menor involucrado
INSERT INTO menor (folio, contador, edad_aproximada, genero, descripcion_fisica)
VALUES ('RPT0000001', 1, 10, 'Masculino', 'Playera roja, mochila azul');

-- Evidencia
INSERT INTO evidencia (folio, tipo_archivo, ruta)
VALUES ('RPT0000001', 'Imagen', 'evidencias/rpt0000001/foto1.jpg');

-- Etiquetas
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Trabajo infantil', 'Categoria');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Alta', 'Prioridad');

INSERT INTO reporteetiqueta (folio, nombre_etiqueta, tipo)
VALUES ('RPT0000001', 'Trabajo infantil', 'Categoria');
INSERT INTO reporteetiqueta (folio, nombre_etiqueta, tipo)
VALUES ('RPT0000001', 'Alta', 'Prioridad');

-- Expediente
INSERT INTO expediente (folio, id_institucion, id_usuario, observaciones)
VALUES ('RPT0000001', @institucion_id, @usuario_admin_id,
        'Caso recién asignado, pendiente de primera visita.');

-- Seguimiento
INSERT INTO seguimiento (folio, id_usuario, descripcion, fecha, fecha_aprox_seguimiento)
VALUES ('RPT0000001', @usuario_admin_id,
        'Se recibió el reporte y se asignó a revisión inicial.', CURRENT_DATE, NULL);
 
-- Token
INSERT INTO token (id_usuario, llave_firma)
VALUES (3, CONVERT('12345' USING BINARY));

-- Trusted proxies
INSERT INTO TRUSTED_PROXIES (IP_1, IP_2, IP_3, IP_4)
VALUES (127, 0, 0, 1);
