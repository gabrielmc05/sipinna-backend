-- =========================================================
-- 1. territorio
-- =========================================================

CREATE TABLE `territorio` (
  `id_territorio` int(11) NOT NULL AUTO_INCREMENT,
  `estado` enum('Aguascalientes','Baja California','Baja California Sur','Campeche','Chiapas','Chihuahua','Ciudad de México','Coahuila','Colima','Durango','Estado de México','Guanajuato','Guerrero','Hidalgo','Jalisco','Michoacán','Morelos','Nayarit','Nuevo León','Oaxaca','Puebla','Querétaro','Quintana Roo','San Luis Potosí','Sinaloa','Sonora','Tabasco','Tamaulipas','Tlaxcala','Veracruz','Yucatán','Zacatecas') DEFAULT NULL,
  `nombre` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  PRIMARY KEY (`id_territorio`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 2. usuario 
-- =========================================================

CREATE TABLE `usuario` (
  `id_usuario` int(11) NOT NULL AUTO_INCREMENT,
  `correo` varchar(254) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `hash_contrasenna` varbinary(255) NOT NULL,
  `nombre` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `apellidos` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `telefono` varchar(30) DEFAULT NULL,
  `fecha_registro` timestamp NOT NULL DEFAULT current_timestamp(),
  `es_activo` tinyint(1) NOT NULL DEFAULT 1,
  `nombre_completo` varchar(130) AS ( CONCAT(`nombre`, ' ', `apellidos`) ) VIRTUAL,
  PRIMARY KEY (`id_usuario`),
  UNIQUE KEY `correo` (`correo`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 3. etiqueta
-- =========================================================

CREATE TABLE `etiqueta` (
  `nombre_etiqueta` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `tipo` enum('Categoria','Prioridad','Riesgo') NOT NULL,
  `es_activo` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`nombre_etiqueta`,`tipo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 4. institucion
-- Depende de: territorio
-- =========================================================

CREATE TABLE `institucion` (
  `id_institucion` int(11) NOT NULL AUTO_INCREMENT,
  `id_territorio` int(11) NOT NULL,
  `nombre` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `tipo` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `correo` varchar(254) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
  `telefono` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
  `direccion` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
  `es_activo` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`id_institucion`),
  KEY `fk_institucion_territorio` (`id_territorio`),
  CONSTRAINT `fk_institucion_territorio`
    FOREIGN KEY (`id_territorio`)
    REFERENCES `territorio` (`id_territorio`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 5. reporte
-- Depende de: usuario
-- =========================================================

CREATE TABLE `reporte` (
  `folio` char(15) NOT NULL,
  `id_usuario` int(11) DEFAULT NULL,
  `correo_contacto` varchar(254) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci DEFAULT NULL,
  `codigo_seguridad_hash` varchar(255) NOT NULL,
  `descripcion` text NOT NULL,
  `fecha_reporte` timestamp NOT NULL DEFAULT current_timestamp(),
  `horario_aproximado` time DEFAULT NULL,
  `es_anonimo` tinyint(1) NOT NULL,
  `desea_seguimiento` tinyint(1) NOT NULL,
  `hay_adulto_involucrado` tinyint(1) NOT NULL,
  PRIMARY KEY (`folio`),
  KEY `fk_reporte_usuario` (`id_usuario`),
  CONSTRAINT `fk_reporte_usuario`
    FOREIGN KEY (`id_usuario`)
    REFERENCES `usuario` (`id_usuario`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 6. administrador
-- Depende de: usuario, institucion
-- =========================================================

CREATE TABLE `administrador` (
  `id_usuario` int(11) NOT NULL,
  `id_institucion` int(11) NOT NULL,
  `rol` enum('Administrador','Alimentador') NOT NULL,
  PRIMARY KEY (`id_usuario`),
  KEY `fk_administrador_institucion` (`id_institucion`),
  CONSTRAINT `fk_administrador_institucion`
    FOREIGN KEY (`id_institucion`)
    REFERENCES `institucion` (`id_institucion`),
  CONSTRAINT `fk_administrador_usuario`
    FOREIGN KEY (`id_usuario`)
    REFERENCES `usuario` (`id_usuario`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 7. reporteetiqueta
-- Depende de: reporte, etiqueta
-- =========================================================

CREATE TABLE `reporteetiqueta` (
  `folio` char(15) NOT NULL,
  `nombre_etiqueta` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci NOT NULL,
  `tipo` enum('Categoria','Prioridad','Riesgo') NOT NULL,
  PRIMARY KEY (`folio`,`nombre_etiqueta`,`tipo`),
  KEY `fk_reporteetiqueta_etiqueta` (`nombre_etiqueta`,`tipo`),
  CONSTRAINT `fk_reporteetiqueta_etiqueta`
    FOREIGN KEY (`nombre_etiqueta`, `tipo`)
    REFERENCES `etiqueta` (`nombre_etiqueta`, `tipo`),
  CONSTRAINT `fk_reporteetiqueta_reporte`
    FOREIGN KEY (`folio`)
    REFERENCES `reporte` (`folio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 8. ubicacion
-- Depende de: reporte, territorio
-- =========================================================

CREATE TABLE `ubicacion` (
  `folio` char(15) NOT NULL,
  `latitud` decimal(9,6) NOT NULL,
  `longitud` decimal(9,6) NOT NULL,
  `direccion_texto` tinytext DEFAULT NULL,
  `id_territorio` int(11) NOT NULL,
  PRIMARY KEY (`folio`),
  KEY `fk_ubicacion_territorio` (`id_territorio`),
  CONSTRAINT `fk_ubicacion_reporte`
    FOREIGN KEY (`folio`)
    REFERENCES `reporte` (`folio`),
  CONSTRAINT `fk_ubicacion_territorio`
    FOREIGN KEY (`id_territorio`)
    REFERENCES `territorio` (`id_territorio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 9. evidencia
-- Depende de: reporte
-- =========================================================

CREATE TABLE `evidencia` (
  `id_evidencia` int(11) NOT NULL AUTO_INCREMENT,
  `folio` char(15) NOT NULL,
  `tipo_archivo` enum('Imagen','Video') NOT NULL,
  `ruta` varchar(255) NOT NULL,
  PRIMARY KEY (`id_evidencia`),
  KEY `fk_evidencia_reporte` (`folio`),
  CONSTRAINT `fk_evidencia_reporte`
    FOREIGN KEY (`folio`)
    REFERENCES `reporte` (`folio`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 10. menor
-- Depende de: reporte
-- =========================================================

CREATE TABLE `menor` (
  `folio` char(15) NOT NULL,
  `contador` tinyint(3) unsigned NOT NULL,
  `edad_aproximada` tinyint(3) unsigned DEFAULT NULL,
  `genero` enum('Masculino','Femenino','No identificado') NOT NULL DEFAULT 'No identificado',
  `descripcion_fisica` tinytext DEFAULT NULL,
  PRIMARY KEY (`folio`,`contador`),
  CONSTRAINT `fk_menor_reporte`
    FOREIGN KEY (`folio`)
    REFERENCES `reporte` (`folio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 11. expediente
-- Depende de: reporte, institucion, administrador
-- =========================================================

CREATE TABLE `expediente` (
  `folio` char(15) NOT NULL,
  `id_institucion` int(11) NOT NULL,
  `id_usuario` int(11) DEFAULT NULL,
  `estado` enum('Nuevo','En revisión','Asignado','En investigación','Canalizado a Procuraduría','Cerrado','Descartado') NOT NULL DEFAULT 'Nuevo',
  `fecha_inicio` timestamp NOT NULL DEFAULT current_timestamp(),
  `fecha_ultima_actualizacion` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `observaciones` text DEFAULT NULL,
  PRIMARY KEY (`folio`),
  KEY `fk_expediente_institucion` (`id_institucion`),
  KEY `fk_expediente_administrador` (`id_usuario`),
  CONSTRAINT `fk_expediente_administrador`
    FOREIGN KEY (`id_usuario`)
    REFERENCES `administrador` (`id_usuario`),
  CONSTRAINT `fk_expediente_institucion`
    FOREIGN KEY (`id_institucion`)
    REFERENCES `institucion` (`id_institucion`),
  CONSTRAINT `fk_expediente_reporte`
    FOREIGN KEY (`folio`)
    REFERENCES `reporte` (`folio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;


-- =========================================================
-- 12. seguimiento
-- Depende de: expediente, administrador
-- =========================================================

CREATE TABLE `seguimiento` (
  `contador` int(11) NOT NULL AUTO_INCREMENT,
  `folio` char(15) NOT NULL,
  `id_usuario` int(11) DEFAULT NULL,
  `descripcion` text NOT NULL,
  `fecha` date NOT NULL DEFAULT curdate(),
  `fecha_aprox_seguimiento` date DEFAULT NULL,
  PRIMARY KEY (`contador`),
  KEY `fk_seguimiento_expediente` (`folio`),
  KEY `fk_seguimiento_administrador` (`id_usuario`),
  CONSTRAINT `fk_seguimiento_administrador`
    FOREIGN KEY (`id_usuario`)
    REFERENCES `administrador` (`id_usuario`),
  CONSTRAINT `fk_seguimiento_expediente`
    FOREIGN KEY (`folio`)
    REFERENCES `expediente` (`folio`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

-- =========================================================
-- 13. token
-- Depende de: usuario
-- =========================================================

CREATE TABLE `token` (
  `id_usuario` int(11) NOT NULL,
  `llave_firma` binary(32) DEFAULT NULL,
  PRIMARY KEY (`id_usuario`),
  CONSTRAINT `id_usuario`
    FOREIGN KEY (`id_usuario`)
    REFERENCES `usuario` (`id_usuario`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;

-- =========================================================
-- 14. AUTH-BACKEND TRUSTED PROXIES
-- =========================================================

CREATE TABLE `TRUSTED_PROXIES` (
  `IP_1` TINYINT NOT NULL,
  `IP_2` TINYINT NOT NULL,
  `IP_3` TINYINT NOT NULL,
  `IP_4` TINYINT NOT NULL,
  `IP` VARCHAR(15) AS ( CONCAT(
  CAST(IP_1 AS VARCHAR(3)), '.',
  CAST(IP_2 AS VARCHAR(3)), '.',
  CAST(IP_3 AS VARCHAR(3)), '.',
  CAST(IP_4 AS VARCHAR(3)))) VIRTUAL,
  CONSTRAINT `PK_IP`
    PRIMARY KEY (IP_1, IP_2, IP_3, IP_4)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
