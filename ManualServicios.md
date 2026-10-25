# Manual de creación, configuración e instalación de los servicios de SIPINNA-app en AWS (v1.0, octubre 2026)

---

**¡Atención!** Por motivos de convención, los comandos, secciones y opciones cuyo uso se implique para los servicios necesarios, se redactarán en su versión inglesa.

## Levantar servicios de AWS
**Prerrequisitos:**
- Contar con los servicios de AWS habilitados en su cuenta administrativa (EC2, Lambda, S3 y RDS).

**Notas:**
- Dado que no se mencione la alteración o selección de una configuración adicional, queda a disposición de las necesidades de la administración. Sin embargo, se recomienda dejarlos con su valor por defecto.

### Instancia EC2
1. En la consola de EC2 desplegar el panel izquierdo, y en la sección **Instances**, hacer click en **Instances**. 
2. Hacer click en el botón **Launch Instances**.
3. En **Name and tags**, ingresar un nombre adecuado (se sugiere: _auth-backend_).
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
        2. Desplegar el menú **Common security groups** y seleccionar el grupo de seguridad previamente creado.
8. Desplegar el menú derecho **Summary**.
    - Hacer click en **Launch instance**. 

### Instancia Lambda (auth-API)

### Instancia Lambda (comms-API)

### Instancia RDS

---

## Instalación y configuración de los programas
