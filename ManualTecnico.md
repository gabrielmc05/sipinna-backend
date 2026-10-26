# Manual técnico de creación, configuración e instalación de los servicios de SIPINNA-app en AWS (v1.0, octubre 2026)

---

Autores:
- Ian Morgado Villaseñor A01754495
- Ana Karen Abrego Flores A01753979
- Juan Pablo Rodriguez Carrillo A01802063
- Rodrigo Castillo Francisco A01799191
- Gabriel de Jesús Manzo Cuevas A01803514

**¡Atención!** Por motivos de convención, los comandos, secciones y opciones cuyo uso se implique para los servicios necesarios, se redactarán en su versión inglesa.

---

## Requerimientos de hardware y software
### Hardware y capacidad
- **Amazon EC2:** instancia Linux para ejecutar `auth-backend`. La familia _t3_ utiliza arquitectura x86-64; si se desea usar ARM64, elegir una familia compatible como _t4g_. El tamaño debe definirse con pruebas de carga y el número esperado de usuarios; _t3.small_ puede servir como punto de partida de pruebas, no como capacidad garantizada para producción.
- **Amazon RDS:** instancia MariaDB con almacenamiento suficiente para la base, índices y crecimiento esperado. Seleccionar capacidad y escalamiento según el volumen de reportes y evidencias.
- **AWS Lambda:** funciones `auth-API` y `comms-api` con arquitectura ARM64, compatible con los paquetes compilados en este manual. Ajustar memoria, tiempo de espera y concurrencia de acuerdo con mediciones de carga.
- **Equipo de administración/despliegue:** conexión a Internet y almacenamiento seguro para la llave SSH y los archivos de configuración sensibles.

### Software y herramientas
- Amazon Linux 2023 para la instancia EC2.
- Go `1.27.1` o una versión posterior compatible con los módulos (`go.mod` de cada servicio declara `go 1.27.1`).
- Git para obtener el código fuente.
- `zip` para empaquetar las funciones Lambda.
- Cliente MariaDB/MySQL para ejecutar y verificar los scripts SQL.
- OpenSSL para generar la clave simétrica que requiere `auth-backend`.
- Docker y Docker Compose solo para levantar la base localmente con `docker-compose.yml`.
- Navegador web para administrar AWS; cliente SSH para conectarse a EC2 si no se usa la consola de AWS.

Verificar las herramientas antes de compilar:

```bash
go version
git --version
zip -v
```

En Amazon Linux, se pueden instalar Git y las herramientas disponibles en los repositorios mediante `dnf`. Verificar que la versión de Go instalada satisfaga `go.mod`; no editar `go.mod` únicamente para evitar instalar la versión requerida.

---

## Dependencias del proyecto
Los tres componentes están escritos en Go y administran sus dependencias mediante sus archivos `go.mod` y `go.sum`:

| Componente | Dependencias directas principales | Función |
|---|---|---|
| `auth-backend` | Gin `1.12.0`, `go-sql-driver/mysql` `1.10.1`, `godotenv` `1.5.1` | API de autenticación y acceso a MariaDB. |
| `auth-API` | AWS Lambda Go `1.55.1`, adaptador Gin para Lambda `0.16.2`, Gin `1.12.0`, `godotenv` `1.5.1` | API de autenticación que delega operaciones a `auth-backend`. |
| `comms-api` | AWS Lambda Go `1.55.1`, adaptador Gin para Lambda `0.16.2`, Gin `1.12.0`, `go-sql-driver/mysql` `1.10.1`, JWT `5.3.1`, `godotenv` `1.5.1` | API de reportes y conexión a MariaDB. |

Instalar o descargar las dependencias desde el directorio de cada componente con `go mod download`. No modificar manualmente `go.sum`; se actualiza al cambiar dependencias con las herramientas de Go.

---

## Cuentas, accesos y credenciales
- La administración necesita una cuenta de AWS con permisos para crear y configurar EC2, Lambda, RDS, grupos de seguridad e IAM. S3 no se integra actualmente en el código. Para la operación cotidiana, usar roles con permisos mínimos y evitar compartir la cuenta raíz de AWS.
- El equipo de despliegue requiere acceso SSH autorizado a EC2 y su llave privada correspondiente. Guardar la llave de forma segura y revocar o reemplazar el acceso cuando cambie el personal autorizado.
- Crear un rol de ejecución independiente para cada Lambda y asignar únicamente los permisos necesarios. `comms-api` requiere acceso de red a RDS.
- Crear en MariaDB un usuario de aplicación llamado `sipinna_admin` con permisos solo sobre `sipinna_db`; no utilizar el usuario maestro de RDS desde la aplicación. El nombre de variable `MARIADB_ROOT_PASSWORD` se conserva por compatibilidad con el código, pero su valor debe ser la contraseña del usuario de aplicación.
- `auth-backend` requiere además `MARIADB_DATABASE_SYMMETRIC_KEY` y `ROOT_ADMIN_INITIAL_PASSWORD`. Generar la clave simétrica para cada instalación y establecer una contraseña inicial segura para el administrador.
- Guardar las contraseñas, la clave simétrica y llaves privadas en un gestor de secretos autorizado. No incluir sus valores en este manual, el repositorio, capturas de pantalla, paquetes de despliegue ni registros compartidos. Usar una cuenta de prueba para las comprobaciones.

---

## Licencias
Estos programas son _software_ libre, puede distribuirlo y/o modificarlos bajo los términos de la **GNU Affero General Public License** 3 o superior, como es publicada por la **Free Software Foundation**.

Las dependencias de Go tienen licencias libres propias, que deben revisarse antes de publicar o distribuir el sistema. Los archivos `go.mod` y `go.sum` registran módulos y versiones. Deberán verificarse por separado las condiciones de uso de AWS y de las imágenes de contenedor.

---

## Levantar servicios de AWS
**Prerrequisitos:**
- Contar con los servicios de AWS habilitados en su cuenta administrativa (EC2, Lambda y RDS). S3 no se integra actualmente en el código; consultar la sección de limitaciones antes de asumir que las evidencias se almacenan allí.

**Notas:**
- Dado que no se mencione la alteración o selección de una configuración adicional, queda a disposición de las necesidades de la administración. Sin embargo, se recomienda dejarlos con su valor por defecto.

### Instancia EC2
1. En la consola de **EC2** desplegar el panel izquierdo, y en la sección **Instances**, hacer click en **Instances**. 
2. Hacer click en el botón **Launch Instances**.
3. En **Name and tags**, ingresar un nombre adecuado (se sugiere: _sipinna-auth-backend_).
4. Desplegar **Application and OS Images (Amazon Machine Image)**  y en abrir la subsección **Quick Start**.
    - Seleccionar _Amazon Linux_.
    - En **Amazon Machine Image** seleccionar _Amazon Linux 2023 kernel-6.18 AMI_.
    - En **Architecture** seleccionar _64-bit (x86)_, para usar la familia _t3_ indicada abajo. Si se elige ARM64, seleccionar una familia compatible con esa arquitectura, como _t4g_, y compilar el binario para ARM64.
5. Desplegar **Instance type**.
    - Seleccionar un tipo de la familia _t3_ compatible con x86-64 y con la escala esperada. _t3.small_, _t3.medium_ o _t3.large_ son puntos de partida, no mínimos garantizados; validar la capacidad con pruebas de carga.
6. Desplegar **Key pair (login)**.
    - Para despliegue inicial (o si se desea crear una nueva llave):
        1. En **Key pair name** hacer click en **⟳ Create new key pair**.
        2. En la ventana emergente, y en la sección **Key pair name**, ingresar un nombre (se sugiere: _auth-backend-key_).
        3. En **Key pair type** seleccionar _RSA_.
        4. En **Private key file format** seleccionar _.pem_.
        5. Hacer click en **Create key pair**.
        6. Comenzará una descarga de la nueva llave creada, esta deberá ser almacenada en un lugar seguro por la administración. Se usará para las conexiones posteriores.
    - Para despliegues posteriores (o si no se desea crear una nueva llave):
        1. En **Key pair name** desplegar las opciones.
        2. Seleccionar una llave previamente creada y cuyo acceso se conozca
7. Desplegar **Network settings**.
    - Para despliegue inicial (o si se desea crear un nuevo grupo de seguridad):
       1. En la sección **Firewall (security groups)** seleccionar **Create security group**.
       2. Marcar _Allow SSH traffic from_, desplegar las opciones y seleccionar o ingresar las direcciones permitidas por la administración para la conexión segura a la instancia (idealmente, seleccionar la opción _Custom_).
       3. No exponer directamente el puerto de la aplicación. Para acceso externo, usar un balanceador o proxy con HTTPS y permitir el tráfico al puerto de la aplicación solo desde dicho componente. La aplicación actual escucha HTTP en el puerto `8080`; abrir el puerto HTTPS en el grupo de seguridad no configura TLS por sí mismo.
    - Para despliegues posteriores (o si no se desea crear un nuevo grupo de seguridad):
        1. En la sección **Firewall (security groups)** seleccionar **Select existing security group**.
        2. En **Common security groups** desplegar **Select security groups** y seleccionar el grupo de seguridad previamente creado.
8. Desplegar el menú derecho **Summary**.
    - Hacer click en **Launch instance**. 

### Instancia Lambda (auth-API)
1. En la consola de **Lambda**, abrir **Functions**.
2. Hacer click en **Create function**.
3. En **Basic information**:
    - En **Function name** ingresar un nombre adecuado (se sugiere: _sipinna-auth-api_).
    - En **Runtime**, seleccionar el runtime personalizado basado en Amazon Linux 2023 (`provided.al2023`).
4. En **Custom settings**, desplegar **Additional settings**.
    - En la sección **General**, marcar _ARM64 architecture_.
    - En la sección **General**, marcar _Custom execution role_. Se desplegará el menú lateral **Configure custom execution role**.
        - Para despliegue inicial (o si se desea crear un nuevo rol):
            1. Hacer click en **Create new role**.
            2. Asignar un nombre adecuado.
            3. En **Additional policy**, asignar una política de ejecución mínima que permita escribir registros en CloudWatch; no asignar permisos administrativos generales.
            4. Hacer click en **Create role**.
            5. Hacer click en **Save**.
        - Para despliegues posteriores (o si no se desea crear un nuevo rol):
            1. Desplegar **Choose an existing role**.
            2. Seleccionar un rol existente (para pruebas iniciales, se sugiere el rol _LabRole_ si es que existe).
            3. Hacer click en **Save**.
    - En la sección **Networking**, marcar _Function URL_. Se desplegará el menú lateral **Configure function url**.
        1. En **Auth type**, seleccionar la opción que corresponda al diseño de autenticación aprobado. Una URL con autenticación _NONE_ queda accesible públicamente; usarla solo si la API aplica controles adecuados y se acepta esa exposición.
        2. En **Invoke mode** marcar _Configure cross-origin resource sharing (CORS)_.
        3. Hacer click en **Save**.
5. Hacer click en **Create function**. Configurar el nombre del handler como `bootstrap` al publicar el paquete.

### Instancia Lambda (comms-API)
1. En la consola de **Lambda**, abrir **Functions**.
2. Hacer click en **Create function**.
3. En **Basic information**:
    - En **Function name** ingresar un nombre adecuado (se sugiere: _sipinna-comms-api_).
    - En **Runtime**, seleccionar el runtime personalizado basado en Amazon Linux 2023 (`provided.al2023`).
4. En **Custom settings**, desplegar **Additional settings**.
    - En la sección **General**, marcar _ARM64 architecture_.
    - En la sección **General**, marcar _Custom execution role_. Se desplegará el menú lateral **Configure custom execution role**.
        - Para despliegue inicial (o si se desea crear un nuevo rol):
            1. Hacer click en **Create new role**.
            2. Asignar un nombre adecuado.
            3. En **Additional policy**, asignar una política de ejecución mínima que permita escribir registros en CloudWatch y el acceso de red de Lambda a la VPC; no asignar permisos administrativos generales.
            4. Hacer click en **Create role**.
            5. Hacer click en **Save**.
        - Para despliegues posteriores (o si no se desea crear un nuevo rol):
            1. Desplegar **Choose an existing role**.
            2. Seleccionar un rol existente (para pruebas iniciales, se sugiere el rol _LabRole_ si es que existe).
            3. Hacer click en **Save**.
    - En **Networking**, configurar el acceso a la VPC y subredes privadas donde se encuentra RDS. Asignar un grupo de seguridad para Lambda y permitir su salida por TCP `3306` hacia RDS.
    - Si se requiere un acceso HTTP público a la API, configurar también su mecanismo de invocación. La conectividad con RDS no se habilita creando una Function URL.
5. Hacer click en **Create function**. Configurar el nombre del handler como `bootstrap` al publicar el paquete.

### Instancia RDS
1. En la consola de **RDS**, desplegar el panel izquierdo y hacer click en **Databases**.
2. Hacer click en **Create database**.
3. En **Choose a database creation method**, seleccionar **Standard create**.
4. En **Engine options**, seleccionar **MariaDB** y una versión compatible con el esquema del proyecto (la base usa características y collations propias de MariaDB).
5. En **Templates**, seleccionar la opción que corresponda al entorno y presupuesto de la administración.
6. En **Settings**:
    - En **DB instance identifier**, ingresar un nombre adecuado (se sugiere: _sipinna-db_).
    - Configurar las credenciales administrativas y almacenarlas en un lugar seguro. No reutilizar contraseñas de otros servicios.
7. En **Instance configuration** y **Storage**, seleccionar la capacidad adecuada para el uso esperado. Para una primera instalación municipal, puede conservarse la configuración recomendada por AWS y ajustarse posteriormente.
8. En **Connectivity**:
    - Seleccionar la misma **VPC** que usa la instancia EC2 que ejecutará los servicios.
    - En **Public access**, seleccionar **No**.
    - En **VPC security group (firewall)**, crear o seleccionar un grupo de seguridad para RDS.
    - En las reglas de entrada del grupo de seguridad, permitir **MYSQL/Aurora (TCP 3306)** desde el grupo de seguridad de EC2. Después de configurar `comms-api` Lambda en la VPC, permitir también el tráfico desde el grupo de seguridad de esa función. No abrir el puerto a todas las direcciones de Internet.
9. Desplegar **Additional configuration** y, en **Initial database name**, ingresar _sipinna_db_.
10. Revisar el resumen y hacer click en **Create database**. Esperar a que el estado de la instancia sea **Available**.

---

## Instalación y configuración de los programas

### Configuración de RDS y comprobación de la conexión
1. En la consola de **RDS**, abrir la instancia creada y, en **Connectivity & security**, copiar el **Endpoint** y el puerto. El puerto predeterminado de MariaDB en RDS es `3306`.
2. Conectarse al endpoint desde la instancia EC2 o desde una máquina autorizada temporalmente por el grupo de seguridad, y cargar los archivos SQL del proyecto, en este orden:
    1. `db/schema.sql`
    2. `db/inserts_iniciales.sql`
    3. `db/procedimientos.sql`
    
    **Importante:** `schema.sql` elimina y vuelve a crear `sipinna_db`; ejecutarlo solo durante la instalación inicial o después de contar con un respaldo, ya que elimina los datos existentes. El usuario que lo ejecute debe tener permisos para crear y eliminar esa base.
    
    Si la carga se realiza desde una computadora administrativa, permitir temporalmente el puerto `3306` solo desde su dirección IP y retirar esa regla al terminar.
3. Crear en la base de datos el usuario de aplicación `sipinna_admin` y otorgarle los permisos necesarios únicamente sobre `sipinna_db`. Los servicios de `auth-backend` y `comms-api` usan este nombre de usuario.
4. Descargar el certificado CA vigente para RDS desde la documentación oficial de AWS y guardar el bundle como `global-bundle.pem` en el directorio de ejecución de cada servicio que se conecte a la base. Los servicios configuran TLS y esperan encontrar ese archivo.
5. Configurar en `auth-backend` y `comms-api` los datos de conexión. En Lambda, establecerlos en **Configuration > Environment variables** o mediante la solución de secretos aprobada; no incluir contraseñas en los paquetes ZIP:
    - `MARIADB_ADDRESS`: endpoint de RDS seguido de `:3306` (por ejemplo, `sipinna-db.xxxxxxxxxxxx.region.rds.amazonaws.com:3306`).
    - `MARIADB_DATABASE`: `sipinna_db`.
    - `MARIADB_ROOT_PASSWORD`: contraseña asignada al usuario de aplicación `sipinna_admin` (el código actual conserva este nombre de variable; no es la contraseña maestra de RDS).
    
    No almacenar contraseñas en el repositorio. Verificar la conexión desde EC2 antes de iniciar los servicios.

### Configuración de EC2
**Prerrequisitos:**
- Contar con la instancia de Amazon EC2 creada.
- Tener una máquina autorizada previamente con un cliente SSH (opcional).

**Notas:**
- Si la configuración inicial ha sido terminada, puede ir directamente a [Configuración de EC2 (continuación)](#configuración-de-ec2-continuación "Configuración de EC2 (continuación)").

#### Configuración de EC2 inicial
1. En la consola de **EC2** desplegar el panel izquierdo, y en la sección **Instances**, hacer click en **Instances**.
2. Asegurarse que la instancia de EC2 previamente creada esté en estado (columna **Instance state**) _☑ Running_.
3. Hacer click en el enlace de la columna **Instance ID** correspondiente a la instancia de EC2 previamente creada
4. En **Instance summary** hacer click a **Connect**.
5. Seguir las instrucciones de conexión que correspondan al método (_In web browser_ y _In SSH client_).
6. Actualizar e instalar las dependencias:

```bash
sudo dnf update
```

```bash
sudo dnf install git
```

Instalar Go `1.27.1` o posterior compatible con `go.mod` siguiendo el método aprobado para el entorno. Comprobar la versión con `go version`.

7. Instalar un editor de texto de terminal (opcional).

```bash
sudo dnf install vim
```

8. Clonar el repositorio de git:

```bash
git clone https://github.com/gabrielmc05/sipinna-backend.git
```

```bash
cd sipinna-backend/auth-backend
```

9. Crear archivos `.env` y `.env.keys` utilizando como plantilla los archivos `.env.sample` y `.env.keys.sample` respectivamente.

```bash
cp .env.sample .env && cp .env.keys.sample .env.keys
```

10. Modificar las variables de entorno del archivo `.env` con respecto a la [configuración de la instancia de RDS](#configuración-de-rds-y-comprobación-de-la-conexión "Configuración de RDS y comprobación de la conexión"), ejemplo:

> `MARIADB_ROOT_PASSWORD="<contraseña-del-usuario-sipinna_admin>"`
> 
> `MARIADB_DATABASE="sipinna_db"`
> 
> `MARIADB_ADDRESS="sipinna-db.xxxxxxxxxxxx.region.rds.amazonaws.com:3306"`

11. Generar una clave simétrica aleatoria de 32 bytes en formato hexadecimal, puede hacerse con `openssl`:

```bash
openssl rand -hex 32
```

12. Guardar en `.env.keys` la clave simétrica generada en el paso anterior y establecer una contraseña inicial segura para `ROOT_ADMIN_INITIAL_PASSWORD`; no copiar valores de ejemplo o reutilizar credenciales. El programa restablece la contraseña del administrador raíz al valor configurado cada vez que inicia y actualmente escribe ese valor en sus registros. Proteger esos registros y cambiar este comportamiento antes de considerar el servicio listo para producción.

El archivo `.env.keys` debe contener `MARIADB_DATABASE_SYMMETRIC_KEY` y `ROOT_ADMIN_INITIAL_PASSWORD`; mantener sus valores fuera del repositorio.

13. Copiar el archivo certificado CA `global-bundle.pem` a la carpeta actual, para más información consulte el punto 3 de [Configuración de RDS y comprobación de la conexión](#configuración-de-rds-y-comprobación-de-la-conexión "Configuración de RDS y comprobación de la conexión")

#### Configuración de EC2 (continuación)
1. Actualizar el repositorio (mantenerse en la rama de producción: `main`).

```bash
git pull
```

2. Compilar el programa con una versión de Go compatible con `go.mod`. No cambiar `go.mod` solo para compilar con una versión anterior.

```bash
go build
```

3. Ejecutar el programa.

```bash
./auth-backend
```

### Publicación segura de auth-backend
El servidor de `auth-backend` escucha HTTP en el puerto `8080`; configurar `TLSConfig` no habilita HTTPS mientras se use `ListenAndServe`. Antes de permitir tráfico desde Internet:
1. Configurar un balanceador o proxy con certificado TLS válido y reenvío HTTP al puerto `8080` de EC2.
2. Restringir el grupo de seguridad de EC2 para aceptar el puerto `8080` solo desde el grupo de seguridad del balanceador/proxy. Mantener SSH limitado a las IP autorizadas.
3. Usar la URL HTTPS del balanceador en `AUTH_BACKEND_URL`; no publicar la dirección de EC2 con el puerto `8080`.
4. Configurar y probar la comprobación de estado del balanceador. El servicio no define actualmente una ruta `/health`, por lo que debe elegirse una comprobación compatible o implementarse una ruta de salud antes de depender de ella.

### Configuración de Lambda (auth-API)
**Prerrequisitos:**
- Contar con la instancia de AWS Lambda creada y correspondiente a esta función.
- Haber terminado la [configuración de la instancia de EC2](#configuración-de-ec2 "Configuración de EC2").
- Contar con un dominio HTTPS de balanceador o proxy que reenvíe solicitudes al backend en EC2.
- Tener una máquina que cuente con las dependencias para la compilación y despliegue del programa (`go`, `zip` y `git`). El build indicado usa `CGO_ENABLED=0`, por lo que no requiere `gcc`.
- Haber clonado el repositorio del `git` en esa máquina:

```bash
git clone https://github.com/gabrielmc05/sipinna-backend.git
```

```bash
cd sipinna-backend/auth-API
```

**Notas:**
- Si la configuración inicial ha sido terminada, puede ir directamente a [Configuración de Lambda (auth-API) (continuación)](#configuración-de-lambda-auth-api-continuación "Configuración de Lambda (auth-API) (continuación)").

#### Configuración de Lambda (auth-API) inicial
1. Crear archivo `.env` utilizando como plantilla el archivo `.env.sample`.

```bash
cp .env.sample .env
```

2. Modificar `AUTH_BACKEND_URL` con la URL HTTPS del balanceador o proxy que expone `auth-backend`. La aplicación de EC2 escucha HTTP en el puerto `8080`; no usar `https` directamente contra ese puerto sin configurar TLS:

> `AUTH_BACKEND_URL="https://<dominio-del-balanceador>"`

#### Configuración de Lambda (auth-API) (continuación)
1. Actualizar el repositorio (mantenerse en la rama de producción: `main`).

```bash
git pull
```

2. Compilar el programa, el ejecutable para la instancia de Lambda quedará en un archivo con nombre `bootstrap`.

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc
```

3. Crear el paquete de despliegue en un archivo `.zip`. Este servicio necesita el archivo `.env` porque la inicialización actual lo carga obligatoriamente; su única variable es la URL no secreta de `auth-backend`. No agregar credenciales:

```bash
zip auth-api.zip bootstrap .env
```

4. Abrir la consola de **Lambda**.
5. En la consola de **Lambda**, abrir **Functions**.
6. Hacer click en el enlace de la instancia correspondiente (en la columna **Function name**).
7. En la sección **Code**, en la ventana **Code source**, desplegar **Update**.
8. Hacer click en la opción _Update from a .zip file_.
9. Hacer click en **Choose file**.
10. Navegar el sistema de archivos y entrar a la carpeta del repositorio `sipinna-backend/auth-API`.
11. Seleccionar el archivo `auth-api.zip`.
12. Hacer click en **Update**.

### Configuración de Lambda (comms-API)
**Prerrequisitos:**
- Contar con la instancia de AWS Lambda creada y correspondiente a esta función.
- Tener una máquina que cuente con las dependencias para la compilación y despliegue del programa (`go`, `zip` y `git`). El build indicado usa `CGO_ENABLED=0`, por lo que no requiere `gcc`.
- Haber clonado el repositorio del `git` en esa máquina:

```bash
git clone https://github.com/gabrielmc05/sipinna-backend.git
```

```bash
cd sipinna-backend/comms-api
```

**Notas:**
- Si la configuración inicial ha sido terminada, puede ir directamente a [Configuración de Lambda (comms-API) (continuación)](#configuración-de-lambda-comms-api-continuación "Configuración de Lambda (comms-API) (continuación)").

#### Configuración de Lambda (comms-API) inicial
1. Crear archivo `.env` utilizando como plantilla el archivo `.env.sample`.

```bash
cp .env.sample .env
```

2. Configurar estas variables de conexión en **Configuration > Environment variables** de Lambda. Si se trabaja localmente, agregarlas al archivo `.env`, que no debe compartirse. El código actual no obtiene secretos automáticamente de AWS Secrets Manager; esa integración requiere implementarse aparte:

> `MARIADB_ROOT_PASSWORD="<contraseña-del-usuario-sipinna_admin>"`
> 
> `MARIADB_DATABASE="sipinna_db"`
> 
>`MARIADB_ADDRESS="sipinna-db.xxxxxxxxxxxx.region.rds.amazonaws.com:3306"`

3. Copiar el archivo certificado CA `global-bundle.pem` a la carpeta actual, para más información consulte el punto 3 de [Configuración de RDS y comprobación de la conexión](#configuración-de-rds-y-comprobación-de-la-conexión "Configuración de RDS y comprobación de la conexión")

#### Configuración de Lambda (comms-API) (continuación)
1. Actualizar el repositorio (mantenerse en la rama de producción: `main`).

```bash
git pull
```

2. Compilar el programa, el ejecutable para la instancia de Lambda quedará en un archivo con nombre `bootstrap`.

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc
```

3. Crear el paquete de despliegue en un archivo `.zip`. Las variables de conexión se configuran en Lambda y no se incluyen en el paquete:

```bash
zip comms-api.zip bootstrap global-bundle.pem
```

4. Abrir la consola de **Lambda**.
5. En la consola de **Lambda**, abrir **Functions**.
6. Hacer click en el enlace de la instancia correspondiente (en la columna **Function name**).
7. En la sección **Code**, en la ventana **Code source**, desplegar **Update**.
8. Hacer click en la opción _Update from a .zip file_.
9. Hacer click en **Choose file**.
10. Navegar el sistema de archivos y entrar a la carpeta del repositorio `sipinna-backend/comms-api`.
11. Seleccionar el archivo `comms-api.zip`.
12. Hacer click en **Update**.

---

## Scripts y archivos de datos
| Archivo | Función | Precauciones |
|---|---|---|
| `db/schema.sql` | Crea la estructura de `sipinna_db`. | Elimina y vuelve a crear la base; ejecutar solo en una instalación inicial o después de un respaldo aprobado. |
| `db/inserts_iniciales.sql` | Inserta los datos iniciales requeridos, incluidos catálogos e instituciones. | Ejecutar una vez después del esquema; una segunda ejecución puede duplicar datos o fallar por restricciones únicas. |
| `db/procedimientos.sql` | Crea los procedimientos y funciones utilizados por los servicios; elimina las versiones existentes antes de recrearlas. | Aplicar después del esquema. Revisar cambios antes de ejecutarlo sobre una base con datos de producción. |
| `db/municipios.sql` | Inserta nombres de municipios en `territorio`. | El propio archivo indica que no hace falta añadirlos; no ejecutarlo en la carga inicial descrita aquí, pues puede duplicar registros. |
| `docker-compose.yml` | Levanta MariaDB local y ejecuta los tres scripts iniciales de la carpeta `db/` al crear por primera vez el volumen. | Está destinado al desarrollo local, no sustituye la configuración de RDS ni debe exponerse a Internet. |

Antes de ejecutar SQL de actualización en una base existente, generar y comprobar un respaldo. Conservar un registro de scripts aplicados y la versión del código asociada.

## Verificación de la instalación
Realizar las comprobaciones en un entorno de pruebas antes de producción. No ejecutar pruebas que inicialicen o modifiquen datos contra una base productiva.

### Base de datos
1. En la consola de RDS, comprobar que el estado de la instancia sea **Available** y que el endpoint y puerto coincidan con `MARIADB_ADDRESS`.
2. Desde EC2 o desde una ubicación temporalmente autorizada por el grupo de seguridad, abrir una sesión con el cliente MariaDB. El cliente pedirá la contraseña de forma interactiva:

```bash
mariadb --host=<endpoint-rds> --port=3306 --user=sipinna_admin --password sipinna_db
```

3. Ejecutar estas consultas y comprobar que existe la base, las tablas y los procedimientos esperados:

```sql
SELECT VERSION();
SHOW TABLES;
SHOW PROCEDURE STATUS WHERE Db = 'sipinna_db';
SHOW SESSION STATUS LIKE 'Ssl_cipher';
```

La fila de `Ssl_cipher` debe mostrar un valor no vacío. Al finalizar, retirar cualquier regla temporal de acceso desde una computadora administrativa.

### Compilación y servicios
1. En cada directorio (`auth-backend`, `auth-API` y `comms-api`), descargar dependencias y comprobar que el código compila:

```bash
go mod download
go build ./...
```

2. Ejecutar las pruebas disponibles (`go test ./...`) en un entorno de pruebas con las variables de entorno, certificados y base de datos configurados. `auth-backend` puede modificar la contraseña inicial del administrador al arrancar; no ejecutar sus pruebas contra producción.
3. Para `auth-backend`, comprobar desde el balanceador HTTPS que una solicitud de autenticación de prueba llega al servicio y recibe una respuesta esperada. No existe actualmente una ruta dedicada `/health`.
4. Para `auth-API`, invocar `GET /Login/<correo-de-prueba>` en su URL de Lambda; para un usuario de prueba válido se espera `202 Accepted` y un challenge.
5. Para `comms-api`, invocar `GET /api/v1/territorios` en su URL de API; se espera `200 OK` y una respuesta JSON.
6. Revisar los registros de EC2 y CloudWatch Logs de Lambda para confirmar que no haya errores de arranque, conexión a RDS, TLS o permisos. No compartir registros que puedan contener datos personales o secretos.

### Limitación de almacenamiento de evidencias
`comms-api` guarda actualmente los archivos de evidencia en el sistema de archivos local; no hay integración con S3 en el código revisado. En Lambda, el directorio del paquete no es un almacenamiento persistente y solo se dispone del espacio temporal de ejecución. Configurar `UPLOAD_DIR` en `/tmp` puede permitir escrituras temporales, pero no garantiza que el archivo esté disponible en invocaciones posteriores. Antes de habilitar la carga de evidencias en producción con Lambda, integrar almacenamiento persistente (por ejemplo, S3) y verificar también la lectura de esos archivos.

## Mantenimiento
- **Respaldos:** habilitar respaldos automáticos y snapshots de RDS conforme a la política de retención institucional. Probar periódicamente la restauración en una instancia separada; un respaldo no verificado no confirma que la recuperación funcione.
- **Cambios de esquema:** respaldar antes de modificar datos o rutinas; revisar los scripts SQL y aplicarlos primero en pruebas. Nunca volver a ejecutar `schema.sql` como una actualización.
- **Actualizaciones:** registrar la versión desplegada y aplicar cambios de Go, dependencias, Amazon Linux, MariaDB y funciones Lambda primero en pruebas. Ejecutar compilación, pruebas y comprobaciones de API antes de promoverlos.
- **Monitoreo:** revisar CloudWatch Logs, métricas de Lambda, estado/capacidad de RDS y consumo de EC2. Configurar alertas de indisponibilidad, errores, espacio y costos.
- **Credenciales y accesos:** revisar periódicamente usuarios IAM, roles, reglas de grupos de seguridad y llaves SSH. Rotar credenciales de base y secretos de aplicación con un procedimiento que actualice todos los servicios dependientes.
- **Clave simétrica:** proteger y respaldar `MARIADB_DATABASE_SYMMETRIC_KEY` de forma separada y con acceso restringido. No reemplazarla sin un plan de migración y respaldo, ya que puede afectar los datos cifrados con ella.
- **Certificados TLS:** mantener actualizado `global-bundle.pem` según las indicaciones de AWS y verificar la conexión TLS después de renovarlo.
- **Operación de EC2:** `./auth-backend` inicia el proceso en primer plano y no constituye por sí mismo un servicio persistente. Para producción, configurarlo con un gestor de servicios (por ejemplo, `systemd`) y una política de reinicio, acceso limitado y registros protegidos.
- **Evidencias:** hasta integrar almacenamiento persistente, verificar respaldos de los archivos locales junto con la base y evitar asumir que los datos de Lambda se conservan entre invocaciones.