# Manual técnico de creación, configuración e instalación de los servicios de SIPINNA-app en AWS (v1.0, octubre 2026)

---

Autores:
- Ian Morgado Villaseñor A01754495
- Ana Karen Abrego Flores A01753979
- Juan Pablo Rodriguez Carrillo A01802063
- Rodrigo Castillo Francisco A01799191
- Gabriel de Jesús Manzo Cuevas A01803514

**¡Atención!** Por motivos de convención, los comandos, secciones y opciones cuyo uso se implique para los servicios necesarios, se redactarán en su versión inglesa.

## Levantar servicios de AWS
**Prerrequisitos:**
- Contar con los servicios de AWS habilitados en su cuenta administrativa (EC2, Lambda, S3 y RDS).

**Notas:**
- Dado que no se mencione la alteración o selección de una configuración adicional, queda a disposición de las necesidades de la administración. Sin embargo, se recomienda dejarlos con su valor por defecto.

### Instancia EC2
1. En la consola de **EC2** desplegar el panel izquierdo, y en la sección **Instances**, hacer click en **Instances**. 
2. Hacer click en el botón **Launch Instances**.
3. En **Name and tags**, ingresar un nombre adecuado (se sugiere: _sipinna-auth-backend_).
4. Desplegar **Application and OS Images (Amazon Machine Image)**  y en abrir la subsección **Quick Start**.
    - Seleccionar _Amazon Linux_.
    - En **Amazon Machine Image** seleccionar _Amazon Linux 2023 kernel-6.18 AMI_.
    - En **Architecture** seleccionar _64-bit (Arm)_.
5. Desplegar **Instance type**.
    - Seleccionar cualquier tipo listado de la familia _t3_ según la escala esperada. Se recomienda _t3.small_, _t3.medium_ o _t3.large_ para una primera instalación de uso municipal.
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
       3. Marcar _Allow HTTPS traffic from the internet_.
    - Para despliegues posteriores (o si no se desea crear un nuevo grupo de seguridad):
        1. En la sección **Firewall (security groups)** seleccionar **Select existing security group**.
        2. En **Common security groups** desplegar **Select security groups** y seleccionar el grupo de seguridad previamente creado.
8. Desplegar el menú derecho **Summary**.
    - Hacer click en **Launch instance**. 

### Instancia Lambda (auth-API)
1. En la consola de **Lambda** desplegar el panel izquierdo, y hacer click en **Event functions**. 
2. Hacer click en el botón **Create Event Function**.
3. En **Basic information**:
    - En **Function name** ingresar un nombre adecuado (se sugiere: _sipinna-auth-api_).
    - En **Runtime** desplegar las opciones y seleccionar _Amazon Linux 2023_.
4. En **Custom settings**, desplegar **Additional settings**.
    - En la sección **General**, marcar _ARM64 architecture_.
    - En la sección **General**, marcar _Custom execution role_. Se desplegará el menú lateral **Configure custom execution role**.
        - Para despliegue inicial (o si se desea crear un nuevo rol):
            1. Hacer click en **Create new role**.
            2. Asignar un nombre adecuado.
            3. En **Additional policy**, asignar la política deseada.
            4. Hacer click en **Create role**.
            5. Hacer click en **Save**.
        - Para despliegues posteriores (o si no se desea crear un nuevo rol):
            1. Desplegar **Choose an existing role**.
            2. Seleccionar un rol existente (para pruebas iniciales, se sugiere el rol _LabRole_ si es que existe).
            3. Hacer click en **Save**.
    - En la sección **Networking**, marcar _Function URL_. Se desplegará el menú lateral **Configure function url**.
        1. En **Auth type** marcar _NONE_.
        2. En **Invoke mode** marcar _Configure cross-origin resource sharing (CORS)_.
        3. Hacer click en **Save**.
5. Hacer click en **Create Event Function**.

### Instancia Lambda (comms-API)
1. En la consola de **Lambda** desplegar el panel izquierdo, y hacer click en **Event functions**. 
2. Hacer click en el botón **Create Event Function**.
3. En **Basic information**:
    - En **Function name** ingresar un nombre adecuado (se sugiere: _sipinna-comms-api_).
    - En **Runtime** desplegar las opciones y seleccionar _Amazon Linux 2023_.
4. En **Custom settings**, desplegar **Additional settings**.
    - En la sección **General**, marcar _ARM64 architecture_.
    - En la sección **General**, marcar _Custom execution role_. Se desplegará el menú lateral **Configure custom execution role**.
        - Para despliegue inicial (o si se desea crear un nuevo rol):
            1. Hacer click en **Create new role**.
            2. Asignar un nombre adecuado.
            3. En **Additional policy**, asignar la política deseada.
            4. Hacer click en **Create role**.
            5. Hacer click en **Save**.
        - Para despliegues posteriores (o si no se desea crear un nuevo rol):
            1. Desplegar **Choose an existing role**.
            2. Seleccionar un rol existente (para pruebas iniciales, se sugiere el rol _LabRole_ si es que existe).
            3. Hacer click en **Save**.
    - En la sección **Networking**, marcar _Function URL_. Se desplegará el menú lateral **Configure function url**.
        1. En **Auth type** marcar _NONE_.
        2. Hacer click en **Save**.
5. Hacer click en **Create Event Function**.

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
    - En las reglas de entrada del grupo de seguridad, permitir **MYSQL/Aurora (TCP 3306)** únicamente desde el grupo de seguridad de la instancia EC2. No abrir el puerto a todas las direcciones de Internet.
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
    
    **Importante:** `schema.sql` elimina y vuelve a crear `sipinna_db`; ejecutarlo solo durante la instalación inicial o después de contar con un respaldo, ya que elimina los datos existentes.
    
    Si la carga se realiza desde una computadora administrativa, permitir temporalmente el puerto `3306` solo desde su dirección IP y retirar esa regla al terminar.
3. Crear en la base de datos el usuario de aplicación `sipinna_admin` y otorgarle los permisos necesarios únicamente sobre `sipinna_db`. Los servicios de `auth-backend` y `comms-api` usan este nombre de usuario.
4. Descargar el certificado CA vigente para RDS desde la documentación oficial de AWS y guardar el bundle como `global-bundle.pem` en el directorio de ejecución de cada servicio que se conecte a la base. Los servicios configuran TLS y esperan encontrar ese archivo.
5. Configurar en el entorno de `auth-backend` y `comms-api` los datos de conexión:
    - `MARIADB_ADDRESS`: endpoint de RDS seguido de `:3306` (por ejemplo, `sipinna-db.xxxxxxxxxxxx.region.rds.amazonaws.com:3306`).
    - `MARIADB_DATABASE`: `sipinna_db`.
    - `MARIADB_ROOT_PASSWORD`: contraseña asignada al usuario de aplicación `sipinna_admin` (el código actual conserva este nombre de variable).
    
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
sudo dnf install go git
```

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

> `MARIADB_ROOT_PASSWORD="contrasena-sipinna-db"`
> `MARIADB_DATABASE="sipinna_db"`
> `MARIADB_ADDRESS="sipinna-db.xxxxxxxxxxxx.region.rds.amazonaws.com:3308"`

11. Generar una clave simétrica aleatoria de 32 bytes en formato hexadecimal, puede hacerse con `openssl`:

```bash
openssl rand -hex 32
```

12. Modificar las variables de entorno del archivo `.env.keys` con respecto a la [configuración de la instancia de RDS](#configuración-de-rds-y-comprobación-de-la-conexión "Configuración de RDS y comprobación de la conexión"), **precuación**: la variable `ROOT_ADMIN_INITIAL_PASSWORD` se utiliza para el correcto funcionamiento del programa, tras cada ejecución esta contraseña deberá ser cambiada a una contraseña fuerte y almacenada de forma segura. Ejemplo:

> `MARIADB_DATABASE_SYMMETRIC_KEY="c6e88cff37844d1243a794011440993297287d36292b0e4299c60e9083bc2429"`
> `ROOT_ADMIN_INITIAL_PASSWORD="administrador12345"`

13. Copiar el archivo certificado CA `global-bundle.pem` a la carpeta actual, para más información consulte el punto 3 de [Configuración de RDS y comprobación de la conexión](#configuración-de-rds-y-comprobación-de-la-conexión "Configuración de RDS y comprobación de la conexión")

#### Configuración de EC2 (continuación)
1. Actualizar el repositorio (mantenerse en la rama de producción: `main`).

```bash
git pull
```

2. Compilar el programa, editar el archivo `go.mod` si la versión disponible de `go` es distinta a la versión del módulo.

```bash
go build
```

3. Ejecutar el programa.

```bash
./auth-backend
```

### Configuración de Lambda (auth-API)
**Prerrequisitos:**
- Contar con la instancia de AWS Lambda creada y correspondiente a esta función.
- Haber terminado la [configuración de la instancia de EC2](#configuración-de-ec2 "Configuración de EC2").
- Tener copiado el DNS público de la instancia de EC2 (se encuentra en su panel de configuración), ejemplo: `ec2-xxx-xxx-xxx-xxx.compute-1.amazonaws.com`.
- Tener una máquina que cuente con las dependencias para la compilación y despliegue del programa (`go`, `gcc`, `zip` y `git`).
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

2. Modificar las variables de entorno del archivo `.env` con respecto a la [configuración de la instancia de EC2](#configuración-de-ec2 "Configuración de EC2") y a su puerto correspondiente, ejemplo:

> `AUTH_BACKEND_URL="https://ec2-xxx-xxx-xxx-xxx.compute-1.amazonaws.com:8080"`

#### Configuración de Lambda (auth-API) (continuación)
1. Actualizar el repositorio (mantenerse en la rama de producción: `main`).

```bash
git pull
```

2. Compilar el programa, el ejecutable para la instancia de Lambda quedará en un archivo con nombre `bootstrap`.

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc
```

3. Crear el paquete de despliegue en un archivo `.zip`.

```bash
zip auth-api.zip bootstrap .env
```

4. Abrir la consola de **Lambda**.
5. En la consola de **Lambda** desplegar el panel izquierdo, y hacer click en **Event Functions**.
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
- Tener una máquina que cuente con las dependencias para la compilación y despliegue del programa (`go`, `gcc`, `zip` y `git`).
- Haber clonado el repositorio del `git` en esa máquina:

```bash
git clone https://github.com/gabrielmc05/sipinna-backend.git
```

```bash
cd sipinna-backend/comms-API
```

**Notas:**
- Si la configuración inicial ha sido terminada, puede ir directamente a [Configuración de Lambda (comms-API) (continuación)](#configuración-de-lambda-comms-api-continuación "Configuración de Lambda (comms-API) (continuación)").

#### Configuración de Lambda (comms-API) inicial
1. Crear archivo `.env` utilizando como plantilla el archivo `.env.sample`.

```bash
cp .env.sample .env
```

2. Modificar las variables de entorno del archivo `.env` con respecto a la [configuración de la instancia de RDS](#configuración-de-rds-y-comprobación-de-la-conexión "Configuración de RDS y comprobación de la conexión"), ejemplo:

> `MARIADB_ROOT_PASSWORD="contrasena-sipinna-db"`
> `MARIADB_DATABASE="sipinna_db"`
> `MARIADB_ADDRESS="sipinna-db.xxxxxxxxxxxx.region.rds.amazonaws.com:3308"`

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

3. Crear el paquete de despliegue en un archivo `.zip`.

```bash
zip comms-api.zip bootstrap .env global-bundle.pem
```

4. Abrir la consola de **Lambda**.
5. En la consola de **Lambda** desplegar el panel izquierdo, y hacer click en **Event Functions**.
6. Hacer click en el enlace de la instancia correspondiente (en la columna **Function name**).
7. En la sección **Code**, en la ventana **Code source**, desplegar **Update**.
8. Hacer click en la opción _Update from a .zip file_.
9. Hacer click en **Choose file**.
10. Navegar el sistema de archivos y entrar a la carpeta del repositorio `sipinna-backend/auth-API`.
11. Seleccionar el archivo `comms-api.zip`.
12. Hacer click en **Update**.