# SentinelOS - Diccionario de Errores API

Este documento contiene el registro oficial de los códigos de error devueltos por la API REST y el motor transaccional de SentinelOS.

## Nomenclatura

Todos los errores siguen el formato estructurado: **`ERR_[CATEGORÍA]_[TIPO][CORRELATIVO]`**

### Tipos de Falla (Primer Dígito)
* **`1`**: Errores de sintaxis, formato o petición (Bad Request).
* **`2`**: Errores de lógica de red, semántica o conflictos (Validation).
* **`3`**: Errores de estado, sesión o concurrencia (Transaction Engine).
* **`4`**: Errores de sistema operativo, hardware o dependencias (Kernel/OS).
* **`5`**: Errores de seguridad, autenticación o permisos (Security).

---

## Categoría: NET (Networking)
*Incluye: Interfaces, VLANs, Zonas, Rutas.*

| Código           | Mensaje Base (Inglés)                                    | Definición / Causa                                                                                                         | Ejemplo de `details`           |
|:-----------------|:---------------------------------------------------------|:---------------------------------------------------------------------------------------------------------------------------|:-------------------------------|
| **ERR_NET_1001** | `Invalid IP or CIDR format`                              | La dirección IP enviada en el JSON no cumple con el estándar de parseo de red.                                             | `192.168.1.999/24`             |
| **ERR_NET_1003** | `Resource not found`                                     | El recurso solicitado no existe en la config                                                                               | ``                             |
| **ERR_NET_1002** | `Missing required field`                                 | La petición no incluye un campo obligatorio para crear el recurso.                                                         | `zone name is required`        |
| **ERR_NET_2001** | `Interface acts as parent for VLANs and cannot have IP`  | Intento de asignar una IP de Capa 3 a una interfaz física que actualmente opera en Capa 2 como troncal de VLANs.           | `enp0s8`                       |
| **ERR_NET_2002** | `Parent interface cannot be down while child VLAN is up` | Conflicto de estado administrativo. Una interfaz no puede apagarse si sus subinterfaces están encendidas.                  | `parent: enp0s8, vlan: vlan10` |
| **ERR_NET_2003** | `Resource references unknown entity`                     | El recurso intenta vincularse a otro que no existe en la candidate config actual (Ej. VLAN apuntando a interfaz fantasma). | `unknown zone: LAN`            |
| **ERR_NET_1004** | `Invalid JSON payload`                                   | El cuerpo de la peticion tiene un formato JSON invalido o con tipos incorrectos.                                           |                                |
| **ERR_NET_2004** | `Duplicate IP address detected`                          | Se detecto la misma IP asignada a mas de una interfaz fisica o virtual.                                                    |                                |
| **ERR_NET_2005** | `Vlan references unknown parent interface`               | La vlan intenta anclarse a una interfaz fisica que no existe.                                                              |                                |
| **ERR_NET_1005** | `Invalid zone type` | El tipo de zona enviado no es reconocido por el sistema. | `Use l2 or l3` |
| **ERR_NET_2006** | `Resource already exists` | Intento de crear un recurso con un identificador (nombre o ID) que ya está en uso en la configuración candidata. | `zone LAN already exists` |
| **ERR_NET_1006** | `Invalid protocol` | El protocolo especificado no es válido. Solo se permite tcp o udp. | `protocol: icmp` |
| **ERR_NET_1007** | `Invalid subnet mask` | La máscara de red proporcionada no es válida para el formato IPv4 (ej. no es de 32 bits). | `invalid mask in address LAN_NET` |
| **ERR_NET_1008** | `Invalid port number` | El puerto especificado está fuera del rango válido (1-65535) o es <= 0. | `invalid port in service HTTP` |
---

## Categoría: SEC (Security & NAT)
*Incluye: Reglas NAT, Políticas de Firewall, Filtros.*

| Código | Mensaje Base (Inglés) | Definición / Causa | Ejemplo de `details` |
| :--- | :--- | :--- | :--- |
| **ERR_SEC_1001** | `Invalid NAT action` | La acción enviada no está soportada. Solo se permite masquerade, snat o dnat. | `action: drop` |
| **ERR_SEC_1002** | `Invalid policy action` | La acción de la política no es reconocida. Solo se permite allow o deny. | `action: bypass` |
| **ERR_SEC_5001** | `Missing authorization token` | La petición no incluye el header de Authorization con el token JWT. | - |
| **ERR_SEC_5002** | `Invalid token format` | El header de Authorization no cumple con el esquema 'Bearer <token>'. | `expected Bearer token` |
| **ERR_SEC_5003** | `Invalid or expired token` | El token JWT no superó la validación criptográfica o su tiempo de vida expiró. | `token is expired by 10m` |
| **ERR_SEC_5004** | `User account disabled` | Intento de login de un usuario que existe pero tiene la bandera Enabled en false. | - |
| **ERR_SEC_5005** | `Invalid credentials` | Las credenciales no coinciden o el usuario no existe en la base de datos local. | - |
| **ERR_SEC_2001** | `Invalid policy reference` | La política hace referencia a un objeto nulo o inválido en la memoria. | `policy has nil service` |

---

## Categoría: SYS (System & Transaction)
*Incluye: Motor de Commit, Lockings, Persistencia YAML.*

| Código           | Mensaje Base (Inglés)                   | Definición / Causa                                                                                 | Ejemplo de `details` |
|:-----------------|:----------------------------------------|:---------------------------------------------------------------------------------------------------|:---------------------|
| **ERR_SYS_3001** | `No active config session`              | Intento de editar recursos o hacer commit sin haber iniciado sesión con /api/config/begin primero. | -                    |
| **ERR_SYS_3002** | `Config already locked by another user` | Intento de iniciar una sesión de configuración cuando otro administrador ya tiene el candado.      | `owner: admin`       |
| **ERR_SYS_4001** | `Internal server error`                 | Fallo generico del sistema, lectura de kernel o bd.                                                |                      |
| **ERR_SYS_3002** | `Failed to apply runtime configuration` | En enigine fallo al intentar aplicar la configuracion al sistema operativo.                        |                      |
| **ERR_SYS_3002** | `Failed to create config backup`        | Error al clonar la configuracion para el rollback.                                                 |                      |
