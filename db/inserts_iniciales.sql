-- Territorios e instituciones asociados
-- 1. Acambay de Ruíz Castañeda
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Acambay de Ruíz Castañeda');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Acambay', 'SIPINNA', 'secretaria.ayuntamiento@acambay.gob.mx', '718 101 0000 EXT 304', 'Calle Plaza Hidalgo N° 1, colonia centro 50300 Acambay, México.');

-- 2. Acolman
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Acolman');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Acolman', 'SIPINNA', 'presidencia@acolman.gob.mx; sec.delayuntamiento.25@gmail.com', '55 29 03 93 14', 'Calz. de Los Agustinos No 1, Centro, 55870 Acolman de Nezahualcóyotl, Méx.');

-- 3. Aculco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Aculco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Aculco', 'SIPINNA', 'valeriacruzolvera@gmail.com; aculco2527.sec.ayuntamiento@gmail.com', '(712) 1 07 34 18', 'Plaza de la Constitución No. 1 Colonia Centro, Aculco, Estado de México, C.P. 50360');

-- 4. Almoloya de Alquisiras
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Almoloya de Alquisiras');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Almoloya de Alquisiras', 'SIPINNA', 'ayuntamientoadea25-27@hotmail.com', '(716) 1445196', 'Antiguo Camino a la Fundición Crisantemos 11, 51860 Almoloya de Alquisiras, Méx.');

-- 5. Almoloya de Juárez
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Almoloya de Juárez');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Almoloya de Juárez', 'SIPINNA', 'secretario.ayuntamiento@almolyadejuarez.gob.mx', '(725) 136 0367 ext. 110', 'José María Morelos s/n, Centro, 50900 Villa de Almoloya de Juárez, Méx.');

-- 6. Almoloya del Río
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Almoloya del Río');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Almoloya del Río', 'SIPINNA', 'secretaria_ayuntamiento@almoloyadelrio.gob.mx; selenealvirderomero@yahoo.com.mx', '(713) 1320004 ext 125', 'Plaza Chignahuapan s/n, Centro, Almoloya del Río, Méx.');

-- 7. Amanalco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Amanalco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Amanalco', 'SIPINNA', 'secretaria@amanalco.gob.mx', '(722) 822 1113', 'Manzana 004, 51260 Amanalco de Becerra, Méx.');

-- 8. Amatepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Amatepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Amatepec', 'SIPINNA', 'presidenciamunicipalamatepec@gmail.com', '7161447089', 'Calle Miguel Hidalgpo s/n Amatepec, Méx.');

-- 9. Amecameca
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Amecameca');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Amecameca', 'SIPINNA', 'secretariaayuntamiento614@gmail.com', '(597) 9780989 ext. 0', 'Av. Fray Martin Valencia 5, Amecameca, 56900 Amecameca de Juárez, Méx.');

-- 10. Apaxco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Apaxco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Apaxco', 'SIPINNA', 'secretariaayuntamiento@apaxco.gob.mx', '(599) 9982700', 'Av. Juárez S/N, 55660 Apaxco de Ocampo, Méx.');

-- 11. Atenco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Atenco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Atenco', 'SIPINNA', 'secretaria.ayto.atenco@gmail.com', '(595) 9532894', 'C.27 de Septiembre S/N, Colonia Centro, 56300 San Salvador Atenco, Méx.');

-- 12. Atizapán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Atizapán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Atizapán', 'SIPINNA', 'secretaria.ayuntamiento@atizapansantacruz.gob.mx', '713-131-8603 ext. 101', 'Av. Doctor Guztavo Baz Prada 102 colonia centro, Santa Cruz Atizapan, Atizapan edo de Mexico. Cp 52500');

-- 13. Atizapán de Zaragoza
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Atizapán de Zaragoza');
SET @territorio_id = LAST_INSERT_ID();
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Atizapán de Zaragoza', 'SIPINNA', 'secretaria.ayuntamiento@atizapansantacruz.gob.mx', '713-131-8603 ext. 101', 'Av. Doctor Guztavo Baz Prada 102 colonia centro, Santa Cruz Atizapan, Atizapan edo de Mexico. Cp 52500');
SET @institucion_id = LAST_INSERT_ID();

-- 14. Atlacomulco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Atlacomulco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Atlacomulco', 'SIPINNA', 'sec.ayuntamiento.atlaco25.27@gmail.com', '712 1220333 / 712 1220200', 'Calle Lic Isidro Fabela Alfaro s/n, Centro, 50450 Atlacomulco de Fabela, Méx.');

-- 15. Atlautla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Atlautla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Atlautla', 'SIPINNA', 'secretariaayuntamiento.atlautla@yahoo.com', '5510676175', 'Plaza de la Revolución Manzana 030, San Jacinto, 56970 Atlautla de Victoria, Méx.');

-- 16. Axapusco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Axapusco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Axapusco', 'SIPINNA', 'ayanamexz33@gmail.com; secretariadelayuntamientoaxapu@gmail.com', '5510153896', 'C. 16 Septiembre 10, 55940 Axapusco, Méx.');

-- 17. Ayapango
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ayapango');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ayapango', 'SIPINNA', 'ivanpalmarodriguez81@gmail.com', '5979824149', 'Calle Plaza De La Constitución, Colonia Centro, Código Postal 56760, Municipio Ayapango de Gabriel Ramos');

-- 18. Calimaya
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Calimaya');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Calimaya', 'SIPINNA', 'secretaria.ayuntamiento@calimaya.gob.mx; secretaria.ayuntamiento@calimaya.gob.com.mx', '(722) 1715202 Ext. 130 y 109', 'Manzana 014, San Pedro, 52200 Calimaya de Díaz González, Méx.');

-- 19. Capulhuac
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Capulhuac');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Capulhuac', 'SIPINNA', 'secretaria.ayuntamiento@capulhuac.gob.mx', '713 1336451 / 713 1352045', 'Plaza, Hombres Ilustres 13, Centro, 52700 Capulhuac de Mirafuentes, Méx.');

-- 20. Coacalco de Berriozábal
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Coacalco de Berriozábal');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Coacalco de Berriozábal', 'SIPINNA', 'angelesyasmin23@hotmail.com', '5540815314', 'And. Severiano Reyes Sn, Centro, 55700 Coacalco de Berriozabal, Méx.');

-- 21. Coatepec Harinas
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Coatepec Harinas');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Coatepec Harinas', 'SIPINNA', 'secretariadelayuntamientoch@gmail.com', '7124950503', 'José Ma. Morelos Manzana 011, 51700 Coatepec Harinas, Méx.');

-- 22. Cocotitlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Cocotitlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Cocotitlán', 'SIPINNA', 'secayuntamiento@cocotitlan.gob.mx; emanuelkaedee@gmail.com', '5559820603', 'Salto del Agua S/N, Techichilco, 56680 Cocotitlán, Méx.');

-- 23. Coyotepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Coyotepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Coyotepec', 'SIPINNA', 'secretaria@coyotepec.gob.mx', '(55) 92451190', 'Av Constitución Sn, Centro, 54660 Coyotepec, Méx.');

-- 24. Cuautitlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Cuautitlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Cuautitlán', 'SIPINNA', 'cuautitlanmx.humanista@gmail.com', '(55) 26207800 ext. 103 y 104', 'Alfonso Reyes S/N Fracc, Santa María, 54820 Cuautitlán, Méx.');

-- 25. Chalco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chalco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chalco', 'SIPINNA', 'secretariadelayuntamiento@chalco.gob.mx', '55 59 72 82 80 Ext. 2208', 'C. Reforma 4, Chalco Centro, 56600 Chalco de Díaz Covarrubias, Méx.');

-- 26. Chapa de Mota
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chapa de Mota');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chapa de Mota', 'SIPINNA', 'cdm.presidencia@oulook.com', '588 992 0209', 'Calle Jóefa Ortiz de Dominguez s/n col. Centro, 54350 Chapa de Mota, Mex.');

-- 27. Chapultepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chapultepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chapultepec', 'SIPINNA', 'secretaria.ayuntamiento@chapultepec.gob.mx', '722-263-0855 EXT 2', 'C. Independencia 300, Co. Centro, 52240 Chapultepec, Méx.');

-- 28. Chiautla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chiautla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chiautla', 'SIPINNA', 'ayuntamientochiautlasecretaria@gmail.com', '(595) 9520550 ext. 102', 'Av del Trabajo Manzana 008, Centro, 56030 Chiautla, Méx.');

-- 29. Chicoloapan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chicoloapan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chicoloapan', 'SIPINNA', 'secretariadelayuntamiento@chicoloapan.gob.mx', '5513226108', 'Manzana 020, Cabecera Municipal, 56370 Chicoloapan de Juárez, Méx.');

-- 30. Chiconcuac
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chiconcuac');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chiconcuac', 'SIPINNA', 'secretaria.a@chiconcuac.gob.mx', '5959530757', 'Matamoros S/N, San Pedro, 56270 Chiconcuac de Juárez, Méx.');

-- 31. Chimalhuacán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Chimalhuacán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Chimalhuacán', 'SIPINNA', 'shayuntamiento.chimalhuacan@gmail.com', '5592086282', 'Ignacio Zaragoza Manzana 001, Cabecera Municipal, 56330 Chimalhuacán, Méx.');

-- 32. Donato Guerra
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Donato Guerra');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Donato Guerra', 'SIPINNA', 'secretaria@donatoguerra.gob.mx', '7261101648', 'Nicolas Bravo Manzana 001, 51030 Villa Donato Guerra, Méx.');

-- 33. Ecatepec de Morelos
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ecatepec de Morelos');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ecatepec de Morelos', 'SIPINNA', 'faustinodelacruz@ecatepec.gob.mx', '(55) 5836-1500', 'Benito Juárez Sur Manzana 022, San Cristóbal Centro, 55000 Ecatepec de Morelos, Méx.');

-- 34. Ecatzingo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ecatzingo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ecatzingo', 'SIPINNA', 'ecat.secretaria0125@gmail.com', '5971091769', 'Manzana 010, Santa Catarina, 56990 Ecatzingo de Hidalgo, Méx.');

-- 35. Huehuetoca
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Huehuetoca');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Huehuetoca', 'SIPINNA', 'secretaria.ayuntamiento@huehuetoca.gob.mx', '5939182486 ext 212', 'Av. Benito Juárez S/N, Centro, 54680 Huehuetoca, Méx.');

-- 36. Hueypoxtla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Hueypoxtla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Hueypoxtla', 'SIPINNA', 'ayuntamiento@hueypoxtla.gob.mx', '5996880242', 'Calle Plaza Principal SN, Centro, 55670 Hueypoxtla, Méx.');

-- 37. Huixquilucan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Huixquilucan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Huixquilucan', 'SIPINNA', 'secretaria.del.ayuntamiento@huixquilucan.gob.mx', '5525957050 ext. 2403', 'Galeana S/N, Centro, 52760 Huixquilucan de Degollado, Méx.');

-- 38. Isidro Fabela
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Isidro Fabela');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Isidro Fabela', 'SIPINNA', 'secretaria.ayuntamiento@isidrofabela.gob.mx', '5592392137', 'Av. Constitución #1, Col. Laureles. Municipio de Isidro Fabela, 54480 Méx.');

-- 39. Ixtapaluca
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ixtapaluca');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ixtapaluca', 'SIPINNA', 'secretaria.ayuntamiento@ixtapaluca.gob.mx', '5559728500 ext. 2003', 'Calle Municipio Libre s/n, Centro, 56530 Ixtapaluca, Méx.');

-- 40. Ixtapan de la Sal
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ixtapan de la Sal');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ixtapan de la Sal', 'SIPINNA', 'secretaria.ayuntamiento@ixtapandelasal.gob.mx', '7211433231', 'Plaza de los Mártires Manzana 002, Centro, 51900 Ixtapan de la Sal, Méx.');

-- 41. Ixtapan del Oro
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ixtapan del Oro');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ixtapan del Oro', 'SIPINNA', 'presidenciaixtapan2025.2027@gmail.com', '726 258 2097', 'Manzana 002, 51070 Ixtapan del Oro, Méx.');

-- 42. Ixtlahuaca
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ixtlahuaca');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ixtlahuaca', 'SIPINNA', 'secretaria_ayuntamiento@ixtlahuaca.gob.mx', '712 122 99 02 Ext. 223', 'Palacio Municipal SN, Ixtlahuaca, 50740 Ixtlahuaca de Rayón, Méx.');

-- 43. Xalatlaco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Xalatlaco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Xalatlaco', 'SIPINNA', 'secayteoxal@gmail.com', '(713) 690 49 83', '16 de Septiembre 1, Centro, 52680 Xalatlaco, Méx.');

-- 44. Jaltenco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Jaltenco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Jaltenco', 'SIPINNA', 'secretariaaytojaltenco2025.2027@gmail.com', '55 49 11 55 15 ext. 107', 'Vicente Guerrero, San Andrés, 55780 Jaltenco, Méx.');

-- 45. Jilotepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Jilotepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Jilotepec', 'SIPINNA', 's.ayuntamiento@jilotepecmex.gob.mx; jubris_05@hotmail.com', '7617344172', 'Leona Vicario 101, Centro, 54240 Jilotepec, Méx.');

-- 46. Jilotzingo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Jilotzingo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Jilotzingo', 'SIPINNA', 'secjilo24.27@gmail.com', '(55) 89969100 ext. 251', 'Manzana 001, 54570 Santa Ana Jilotzingo, Méx.');

-- 47. Jiquipilco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Jiquipilco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Jiquipilco', 'SIPINNA', 'secretaria.ayuntamiento@jiquipilco.gob.mx', '7121114151', 'Ricardo Flores Magón 5, Centro, 50800 Jiquipilco, Méx.');

-- 48. Jocotitlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Jocotitlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Jocotitlán', 'SIPINNA', 'ayuntamiento@jocotitlan.gob.mx', '7121611706', 'Constituyentes 1, Centro, 50700 Cdad. de Jocotitlán, Méx.');

-- 49. Joquicingo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Joquicingo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Joquicingo', 'SIPINNA', 'presidencia@joquicingo.gob.mx; presidenciamunicipaljoquicingo@gmail.com', '714 1400533 / 714 1470631', 'C. Jose M. Morelos Manzana 009, 52370 Joquicingo de León Guzmán, Méx.');

-- 50. Juchitepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Juchitepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Juchitepec', 'SIPINNA', 'secretario@juchitepec.gob.mx', '597 977 0872', 'Palacio Municipal, Manzana 023 S/N, Centro, 56860 Juchitepec, Méx.');

-- 51. Lerma
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Lerma');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Lerma', 'SIPINNA', 'secretaria.ayuntamiento@lerma.gob.mx', '(728) 282 99 61 / (728) 282 99 03 ext. 1121', 'Palacio Municipal, Pl. Juárez No 1, Centro, 52000 Lerma de Villada, Méx.');

-- 52. Malinalco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Malinalco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Malinalco', 'SIPINNA', 'secretaria@malinalco.gob.mx', '712 422 8048', 'Manzana 005, San Juan, 52440 Malinalco, Méx.');

-- 53. Melchor Ocampo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Melchor Ocampo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Melchor Ocampo', 'SIPINNA', 'secretaria.ayuntamiento@melchor-ocampo.gob.mx', '(55) 58780111 Ext. 103', 'Av. Adolfo Lopez Mateos #72, Barrio Señor de Los Milagros, 54880 Melchor Ocampo, Méx.');

-- 54. Metepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Metepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Metepec', 'SIPINNA', 'secretariaayuntamiento@metepec.gob.mx', '722 2358222', 'José Vicente Villada #330, Barrio del Espíritu Santo, 52140 Metepec, Méx.');

-- 55. Mexicaltzingo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Mexicaltzingo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Mexicaltzingo', 'SIPINNA', 'sa.mexicaltzingo25.27@gmail.com', '(722) 9367947', 'Independencia, Centro, 52180 San Mateo Mexicaltzingo, Méx.');

-- 56. Morelos
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Morelos');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Morelos', 'SIPINNA', 'secretaria2527morelos@gmail.com', '712 1239194', 'C. 5 de Mayo Manzana 001, Acotzilapan, 50550 San Bartolo Morelos; domicilio provisional en el auditorio, calle Luisa Isabel Campos de Jiménez Cantú s/n, Barrio Primero, Cabecera Municipal, Morelos, Estado de México.');

-- 57. Naucalpan de Juárez
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Naucalpan de Juárez');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Naucalpan de Juárez', 'SIPINNA', 'presidencia@naucalpan.gob.mx', '5564206250', 'Av. Benito Juárez 39, Naucalpan Centro, 53000 Naucalpan de Juárez, Méx. FRACCIONAMIENTO EL MIRADOR');

-- 58. Nezahualcóyotl
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Nezahualcóyotl');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Nezahualcóyotl', 'SIPINNA', 'mdrendon@nezahulcoyotl.gob.mx', '5557169070 ext 1202', 'Palacio municipal de, C. caballo bayo 93, Benito Juárez, 57000 Cdad. Nezahualcóyotl, Méx.');

-- 59. Nextlalpan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Nextlalpan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Nextlalpan', 'SIPINNA', 'secretariadelayutamiento2527@gmail.com', '5549115151', 'Avenida Ayuntamiento S/N, Central, 55790 Santa Ana Nextlalpan, Méx.');

-- 60. Nicolás Romero
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Nicolás Romero');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Nicolás Romero', 'SIPINNA', 'presidenciamuncipal@nicolasromero.gob.mx; secretariadelayuntamientoNR@gmail.com', '55 4377 6628', 'C. Juárez S/N, Centro, 54400 Cdad. Nicolás Romero, Méx.');

-- 61. Nopaltepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Nopaltepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Nopaltepec', 'SIPINNA', 'secretariadelayuntamientonopal@gmail.com', '592 9245005', 'C. Zaragoza, 55970 Nopaltepec, Méx.');

-- 62. Ocoyoacac
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ocoyoacac');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ocoyoacac', 'SIPINNA', 'secretaria.ayuntamiento@ocoyoacac.gob.mx', '(728) 1066630 ext. 0', 'de los 1, Insurgentes, Ocoyoacac, Méx.');

-- 63. Ocuilan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ocuilan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ocuilan', 'SIPINNA', 'secretaria@ocuilan.gob.mx', '72223761690', 'Lázaro Cárdenas 1, Ocuilan de Artiaga, 52480 Ocuilan de Arteaga, Méx.');

-- 64. El Oro
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'El Oro');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA El Oro', 'SIPINNA', 'secretaria@eloromexico.gob.mx', '71116880751', 'Constitución 27C, El Oro de Hidalgo, 50600 El Oro de Hidalgo, Méx.');

-- 65. Otumba
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Otumba');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Otumba', 'SIPINNA', 'Secretaria_2527@ayuntamientotumba.gob.mx', '5929222344 ext. 105', 'De La Constitución 1, 55900 Otumba de Gómez Farías, Méx.');

-- 66. Otzoloapan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Otzoloapan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Otzoloapan', 'SIPINNA', 'presidencia@otzoloapan.gob.mx; secretaria@otzoloapan.gob.mx', '722 665 9102', 'Plaza De La Constitución S/n, 51130 Otzoloapan, Méx.');

-- 67. Otzolotepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Otzolotepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Otzolotepec', 'SIPINNA', 'secretaria.ayuntamiento@otzolotepec.gob.mx', '7192860314', 'Hidalgo 1, Villa Cuauhtémoc, 52080 Villa Cuauhtémoc, Méx.');

-- 68. Ozumba
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Ozumba');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Ozumba', 'SIPINNA', 'secretario@ozumba.gob.mx', '(597) 9760104', 'Calle Plaza de la Constitución, número exterior 1 piso 1, colonia centro, municipio de Ozumba');

-- 69. Papalotla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Papalotla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Papalotla', 'SIPINNA', 'secretariaayuntamiento81@gmail.com', '(595) 9538579', 'Plaza Morelos No. 1, Colonia Centro.');

-- 70. La Paz
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'La Paz');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA La Paz', 'SIPINNA', 'secretario@losreyeslapaz.gob.mx', '(55) 58550024 ext. 0', 'Manzana 008, Los Reyes, 56400 Los Reyes Acaquilpan, Méx.');

-- 71. Polotitlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Polotitlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Polotitlán', 'SIPINNA', 'ayuntamiento.polotitlan@gmail.com', '4272660683', 'Portal Hidalgo 7-MZ 002, 54200 Polotitlán de la Ilustración, Méx.');

-- 72. Rayón
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Rayón');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Rayón', 'SIPINNA', 'secre.ayto25@gmail.com', '(717) 1320045', 'Plaza 22 de Octubre, Ignacio López Rayón Sur Manzana 006, Santa María Rayón Centro, Guadalupe, 52360 Santa María Rayón, Méx.');

-- 73. San Antonio la Isla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'San Antonio la Isla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA San Antonio La Isla', 'SIPINNA', 'secretariadeayto.sali2527@gmail.com', '717 690 8391', 'C. Vicente Villada Manzana 027, San Antonio la Isla, 52280 San Antonio, Méx.');

-- 74. San Felipe del Progreso
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'San Felipe del Progreso');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA San Felipe del Progreso', 'SIPINNA', 'sayuntamiento@sanfelipedelprogreso.gob.mx', '712 123 7474', 'Plaza y Posadas Garduño S/N, Centro, 50640 San Felipe del Progreso, Méx.');

-- 75. San Martín de las Pirámides
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'San Martín de las Pirámides');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA San Martín de las Pirámides', 'SIPINNA', 'ruben.cabrera@sanmartindelaspiramides.gob.mx', '(594) 9582091 ext. 101', 'Benito Juárez S/N, 55850 San Martín de las Pirámides, Méx.');

-- 76. San Mateo Atenco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'San Mateo Atenco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA San Mateo Atenco', 'SIPINNA', 'secayuntamiento@sanmateoatenco.gob.mx', '7286901100', 'Av. Lic. Benito Juárez 302, San Miguel, 52104 San Mateo Atenco, Méx.');

-- 77. San Simón de Guerrero
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'San Simón de Guerrero');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA San Simón de Guerrero', 'SIPINNA', 'secr.mpal.sansimon.gro@gmail.com', '722 902 8857', 'Manzana 005, 51470 San Simón de Guerrero, Méx.');

-- 78. Santo Tomás
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Santo Tomás');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Santo Tomás', 'SIPINNA', 'secretaria.st.2025.2027@gmail.com', '7261100595', 'Manzana 008, 51100 Santo Tomás de los Plátanos, Méx.');

-- 79. Soyaniquilpan de Juárez
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Soyaniquilpan de Juárez');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Soyaniquilpan de Juárez', 'SIPINNA', 'presidencia2527@soyaniquilpan.gob.mx', '(761) 7353125', 'Manzana 017, 54280 San Francisco Soyaniquilpan, Méx.');

-- 80. Sultepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Sultepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Sultepec', 'SIPINNA', 'ayuntamientosultepec2025.2027@gmail.com', '7296310219', 'Francisco I. Madero 1, Calvario Sultepec, 51600 Sultepec de Pedro Ascencio de Alquisiras, Méx.');

-- 81. Tecámac
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tecámac');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tecámac', 'SIPINNA', 'presidenciamunicipal@tecamac.gob.mx', '5540563880', 'Plaza Principal s/n, Tecámac de Felipe Villanueva, 55740 Tecámac, Méx.');

-- 82. Tejupilco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tejupilco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tejupilco', 'SIPINNA', 'secretaria@tejupilco.gob.mx', '7242675530', 'Palacio Municipal s/n, Centro, 51400 Tejupilco de Hidalgo, Méx.');

-- 83. Temamatla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Temamatla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Temamatla', 'SIPINNA', 'secretariatemamatla2527@outlook.com', '5559429137', 'Pl. Hidalgo 4, 56650 Temamatla, Méx.');

-- 84. Temascalapa
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Temascalapa');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Temascalapa', 'SIPINNA', 'marlonhumberto2005@gmail.com', '55 4100 2690', 'Dr. Gustavo Baz, San Miguel, 55980 Temascalapa, Méx.');

-- 85. Temascalcingo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Temascalcingo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Temascalcingo', 'SIPINNA', 'secretario.ayuntamiento@temascalcingo.gob.mx', '7181260145', 'Portal del Recreo Manzana 003, Temascalcingo, 50400 Temascalcingo de José María Velasco, Méx.');

-- 86. Temascaltepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Temascaltepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Temascaltepec', 'SIPINNA', 'secretario.ayto2025@gmail.com', '7162665103', 'Plaza Juárez No. 1, Colonia Centro, Méx.');

-- 87. Temoaya
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Temoaya');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Temoaya', 'SIPINNA', 'sria.ayuntamiento.temoaya2527@gmail.com', '719 265 0067 ext.107', 'Portal Ayuntamiento 103, Portal José María Morelos y Pavón 102, Centro, 50850 Temoaya, Méx.');

-- 88. Tenancingo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tenancingo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tenancingo', 'SIPINNA', 'SecretariadelAyuntamiento@tenancingo.gob.mx', '(718) 1260005', 'Jardín Morelos 101, Centro, 52400 Tenancingo de Degollado, Méx.');

-- 89. Tenango del Aire
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tenango del Aire');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tenango del Aire', 'SIPINNA', 'tenangodelairepresidencia2527@gmail.com; gobierno2527tenando del aire@gmail.com', '5979825400', 'Pl. de la Constitución 1, Centro, 56780 Tenango del Aire, Méx.');

-- 90. Tenango del Valle
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tenango del Valle');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tenango del Valle', 'SIPINNA', 'secretariadelayutamiento@tenangodelvalle.gob.mx', '717690 0600 EXT.105', 'Constitucin Nte 101, Centro, 52300 Tenango de Arista, Méx.');

-- 91. Teoloyucan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Teoloyucan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Teoloyucan', 'SIPINNA', 'secretariadelayutamiento@teoloyucan.gob.mx', '5939141700', 'Av. Dolores S/N, Tlatilco, 54770 Teoloyucan, Méx.');

-- 92. Teotihuacán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Teotihuacán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Teotihuacán', 'SIPINNA', 'contador_alex10@hotmail.com', '5516970882', 'Juárez 1, Teotihuacan Centro, 55800 Teotihuacán de Arista, Méx.');

-- 93. Tepetlaoxtoc
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tepetlaoxtoc');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tepetlaoxtoc', 'SIPINNA', 'secretariah@tepetlaoxtoc.gob.mx', '595 9230010 / 595 923001', 'Plaza principal S/N col Centro Tlaple 56070');

-- 94. Tepetlixpa
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tepetlixpa');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tepetlixpa', 'SIPINNA', 'secretariadelayuntamientotp@gmail.com', '597 975 0168 ext 127, 5574501252', 'Av. Morelos 10, Tepetlixpa, 56880 Tepetlixpa, Méx.');

-- 95. Tepotzotlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tepotzotlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tepotzotlán', 'SIPINNA', 'sec.aytto@tepotzotlan.gob.mx', '55 5876 0202', 'Plaza Virreynal No 1, San Martin, 54600 Tepotzotlán, Méx.');

-- 96. Tequixquiac
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tequixquiac');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tequixquiac', 'SIPINNA', 'presidenciatx@tequixquiac.gob.mx', '5535924504', '55650, Pl. Cuauhtémoc 1, Centro, Tequixquiac, Méx.');

-- 97. Texcaltitlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Texcaltitlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Texcaltitlán', 'SIPINNA', 'ayuntamientotexcaltitlan@gmail.com; jesusmetexcaltitlan@gmail.com', '72 23 57 57 59', 'Av. Benito Juárez No. 1 Centro.');

-- 98. Texcalyacac
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Texcalyacac');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Texcalyacac', 'SIPINNA', 'SecretariaAyuntamiento.2025.2027@gmail.com', '7131316230', 'Calle Lic. Alfredo Campanella Manzana 002, Centro, 52570 San Mateo Texcalyacac, Méx.');

-- 99. Texcoco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Texcoco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Texcoco', 'SIPINNA', 'secretaria@texcoco.gob.mx', '5959520000 Ext. 2007/2128', 'Primera Cda. de Nezahualcóyotl 110, Centro, 56100 Texcoco de Mora, Méx.');

-- 100. Tezoyuca
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tezoyuca');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tezoyuca', 'SIPINNA', 'secretariaayuntamiento@tezoyuca.lgcg.org', '5949565736', 'C. Pascual Luna Manzana 011, Tezoyuca, 56000 Tezoyuca, Méx.');

-- 101. Tianguistenco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tianguistenco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tianguistenco', 'SIPINNA', 'secretarioayuntamiento2025@tianguistenco.gob.mx', '7131352151 / 7131352048', 'Plaza Libertad Manzana 021, Centro, 52600 Santiago Tianguistenco, Méx.');

-- 102. Timilpan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Timilpan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Timilpan', 'SIPINNA', 'secretariadelayuntamiento2022@hotmail.com', 'Particular: Fabiola Martínez Chimal 71 21 77 13 44', 'Plaza Hidalgo S/N Col. Centro Timilpan');

-- 103. Tlalmanalco
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tlalmanalco');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tlalmanalco', 'SIPINNA', 'secretario@tlalmanalco.gob.mx', '5979779867', 'Fray Martin de Valencia SN, Tlalmanalco, 56700 Tlalmanalco de Velázquez, Méx.');

-- 104. Tlalnepantla de Baz
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tlalnepantla de Baz');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tlalnepantla de Baz', 'SIPINNA', 'miguel.bravo@tlalnepantla.gob.mx', '5553 66 38 00 ext 3858', 'Plaza Gustavo Baz, Tlalnepantla Centro, 54000 Tlalnepantla, Méx.');

-- 105. Tlatlaya
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tlatlaya');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tlatlaya', 'SIPINNA', 'presidencia@tlatlaya.gob.mx', '7161718063', 'Palacio Municipal S/N, Centro, 51550 Tlatlaya, Méx.');

-- 106. Toluca
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Toluca');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Toluca', 'SIPINNA', 'secretaria.ayuntamiento@toluca.gob.mx', '7222761900 ext 214', 'Av. Independencia 207, Centro, 50000 Toluca de Lerdo, Méx.');

-- 107. Tonatico
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tonatico');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tonatico', 'SIPINNA', 'garcialeguizamorosamaria@gmail.com; tonaticosecretariadelayunt@gmail.com', '7211410413', 'Plaza de La Constitución 1, San Gaspar, 51950 Tonatico, Méx.');

-- 108. Tultepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tultepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tultepec', 'SIPINNA', 'secretaria@tultepec.gob.mx', '5558920125', 'Manzana 024, Centro, 54960 Tultepec, Méx.');

-- 109. Tultitlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tultitlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tultitlán', 'SIPINNA', 'presidencia@tultitlan.gob.mx; secretaria@tultitlan.gob.mx', '5526208900 ext 1104', 'Av San Antonio 22, San Bartolo, 54900 Tultitlán de Mariano Escobedo, Méx.; Plaza Hidalgo 1, Cabecera Municipal, Tultitlán, Estado de México C.P. 54900');

-- 110. Valle de Bravo
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Valle de Bravo');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Valle de Bravo', 'SIPINNA', 'secretariadelayuntamiento@valledebravo.gob.mx', '722 889 6287', '5 de Febrero 100, Centro, 51200 Valle de Bravo, Méx.');

-- 111. Villa de Allende
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Villa de Allende');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Villa de Allende', 'SIPINNA', 'secretariahayuntamiento.va@gmail.com; presidencia2527@villa de allende.gob.mx; secretario.h.ayuntamiento@villadeallende.gob.mx', '7262690003 ext 111', 'Manzana 005, 51000 San José Villa de Allende, Méx.');

-- 112. Villa del Carbón
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Villa del Carbón');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Villa del Carbón', 'SIPINNA', 'secretariaayuntamiento@villadelcarbon.gob.mx', '5889130204 ext 109', 'Hidalgo 2, Centro, 54300 Villa del Carbón, Méx.');

-- 113. Villa Guerrero
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Villa Guerrero');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Villa Guerrero', 'SIPINNA', 'presidencia@villaguerrero.gob.mx', '7141460006 ext.101', 'Plaza Morelos S/N, Centro, 51760 Villa Guerrero, Méx.');

-- 114. Villa Victoria
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Villa Victoria');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Villa Victoria', 'SIPINNA', 'ayuntamiento_vv@hotmail.com', '7262515238', 'Calle Lázaro Cardenaz S/N Col. Centro Villa Victoria');

-- 115. Xonacatlán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Xonacatlán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Xonacatlán', 'SIPINNA', 'secretaria.ayuntamiento@xonacatlan.gob.mx', '7192860852 Ext. 119', 'Plaza Juárez #1, Col. Centro, Xonacatlán, 52060 Xonacatlán, Méx.');

-- 116. Zacazonapan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Zacazonapan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Zacazonapan', 'SIPINNA', 'pedro_ponce00@hotmail.com', '7226409730', 'PALACIO MUNICIPAL TOLUCA, Manzana 047 S/N, Centro, 51160 Zacazonapan, Méx.');

-- 117. Zacualpan
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Zacualpan');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Zacualpan', 'SIPINNA', 'secretariazac25.27@gmail.com', '721 1476124', 'Plaza Juárez S/N Centro, Zacualpan, Mex');

-- 118. Zinacantepec
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Zinacantepec');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Zinacantepec', 'SIPINNA', 'secretaria.ayuntamiento@zinacantepec.gob.mx', '722 9177206 ext 129', 'Av. 16 de Septiembre 101, del calvario, 51350 San Miguel Zinacantepec, Méx.');

-- 119. Zumpahuacán
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Zumpahuacán');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Zumpahuacán', 'SIPINNA', 'secretaria@zumpahuacan.gob.mx', '714 146 9033', 'Benito Juárez 1, Centro, 51980 Zumpahuacán, Méx.');

-- 120. Zumpango
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Zumpango');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Zumpango', 'SIPINNA', 'secretaria.ayuntamiento@zumpango.gob.mx', '5634301946 / 5582900512', 'Juárez n, Centro, 55600 Zumpango de Ocampo, Méx.');

-- 121. Cuautitlán Izcalli
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Cuautitlán Izcalli');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Cuautitlán Izcalli', 'SIPINNA', 'secretaria.ayuntamiento@cizcalli.gob.mx', '(55) 58642500 ext 5200', 'Av. Primero de Mayo 100, Col. Centro Urbano, 54740 Cuautitlán Izcalli, Méx.');

-- 122. Valle de Chalco Solidaridad
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Valle de Chalco Solidaridad');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Valle de Chalco Solidaridad', 'SIPINNA', 'secretaria@valledechalco.gob.mx; sipinna@valledechalco.gob.mx', '5549580728', 'Av. del Mazo s/n, Alfredo Baranda, 56600 Valle de Chalco Solidaridad, Méx.');

-- 123. Luvianos
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Luvianos');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Luvianos', 'SIPINNA', 'secretaria.ayuntamiento@luvianos.gob.mx', '7224266086', 'Plaza Venustiano Carranza 1, Centro, 51440 Villa Luvianos, Méx.');

-- 124. San José del Rincón
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'San José del Rincón');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA San José del Rincón', 'SIPINNA', 'secretariadelayuntamiento@sanjosedelrincon.gob.mx', '7121242173', 'Calle Guadalupe Victoria 12, 50660 San José del Rincón Centro, Méx.');

-- 125. Tonanitla
INSERT INTO territorio (estado, nombre)
VALUES ('Estado de México', 'Tonanitla');
INSERT INTO institucion
(id_territorio, nombre, tipo, correo, telefono, direccion)
VALUES
(LAST_INSERT_ID(), 'SIPINNA Tonanitla', 'SIPINNA', 'secretariaayuntamiento100@gmail.com', '5533321115', 'Bicentenario 1, Villas de Santa María, 55785 Santa María Tonanitla, Méx.');

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
VALUES ('admin.inicial@sipinna.gob.mx',
        CONVERT('administrador12345' USING BINARY),
        'Administrador', 'Raíz', '5551112222');
SET @usuario_admin_id = LAST_INSERT_ID();

INSERT INTO administrador (id_usuario, id_institucion, rol)
VALUES (@usuario_admin_id, @institucion_id, 'Administrador');

-- Alimentadores ficticios vinculados a las instituciones de prueba 
--INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
--VALUES ('alimentador.procuraduria.prueba@example.com',
--        CONVERT('holahola' USING BINARY),
--        'María', 'López García', '5554445566');
--SET @usuario_alimentador_procuraduria_id = LAST_INSERT_ID();

--INSERT INTO administrador (id_usuario, id_institucion, rol)
--VALUES (@usuario_alimentador_procuraduria_id, @institucion_procuraduria_id, 'Alimentador');

--INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
--VALUES ('alimentador.dif.prueba@example.com',
--        CONVERT('adiosadios' USING BINARY),
--        'Carlos', 'Hernández Ruiz', '5556667788');
--SET @usuario_alimentador_dif_id = LAST_INSERT_ID();

--INSERT INTO administrador (id_usuario, id_institucion, rol)
--VALUES (@usuario_alimentador_dif_id, @institucion_dif_id, 'Alimentador');

-- Usuario que será el Reportante
--INSERT INTO usuario (correo, hash_contrasenna, nombre, apellidos, telefono)
--VALUES ('reportante.prueba@example.com',
--        CONVERT('Papupro38' USING BINARY),
--        'Juan', 'Ramírez Soto', NULL);
--SET @usuario_reportante_id = LAST_INSERT_ID();

-- Reporte
INSERT INTO reporte (folio, id_usuario, correo_contacto, codigo_seguridad_hash,
                      descripcion, horario_aproximado, es_anonimo,
                      desea_seguimiento, hay_adulto_involucrado)
VALUES ('RPT-INICIAL', @usuario_reportante_id, NULL,
        '',
        'Reporte de prueba.',
        '10:30:00', FALSE, TRUE, TRUE);

-- Ubicación
INSERT INTO ubicacion (folio, latitud, longitud, direccion_texto, id_territorio)
VALUES ('RPT-INICIAL', 19.573800, -99.270500,
        'Cerca de la estación de camiones, Av. Central', @territorio_id);

-- Menor involucrado
INSERT INTO menor (folio, contador, edad_aproximada, genero, descripcion_fisica)
VALUES ('RPT-INICIAL', 1, 10, 'Masculino', 'Playera roja, mochila azul');

-- Evidencia
INSERT INTO evidencia (folio, tipo_archivo, ruta)
VALUES ('RPT-INICIAL', 'Imagen', 'evidencias/rpt0000001/foto1.jpg');

-- Etiquetas
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Trabajo infantil', 'Categoria');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Explotación sexual', 'Categoria');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Utilización para actividades ilícitas', 'Categoria');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Trabajo peligroso', 'Categoria');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Mendicidad forzada', 'Categoria');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Otra situación de explotación y/o vulneración', 'Categoria');

INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Alta', 'Prioridad');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Media', 'Prioridad');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Baja', 'Prioridad');

INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Alto', 'Riesgo');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Medio', 'Riesgo');
INSERT INTO etiqueta (nombre_etiqueta, tipo) VALUES ('Bajo', 'Riesgo');

INSERT INTO reporteetiqueta (folio, nombre_etiqueta, tipo)
VALUES ('RPT-INICIAL', 'Trabajo infantil', 'Categoria');
INSERT INTO reporteetiqueta (folio, nombre_etiqueta, tipo)
VALUES ('RPT-INICIAL', 'Alta', 'Prioridad');

-- Expediente
INSERT INTO expediente (folio, id_institucion, id_usuario, observaciones)
VALUES ('RPT-INICIAL', @institucion_id, @usuario_admin_id,
        'Caso recién asignado, pendiente de primera visita.');

-- Seguimiento
INSERT INTO seguimiento (folio, id_usuario, descripcion, fecha, fecha_aprox_seguimiento)
VALUES ('RPT-INICIAL', @usuario_admin_id,
        'Se recibió el reporte y se asignó a revisión inicial.', CURRENT_DATE, NULL);
 
-- Token
--INSERT INTO token (id_usuario, llave_firma)
--VALUES (3, CONVERT('12345' USING BINARY));

-- Trusted proxies
INSERT INTO TRUSTED_PROXIES (IP_1, IP_2, IP_3, IP_4)
VALUES (127, 0, 0, 1);
