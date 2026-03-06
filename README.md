# Comercio API

Servicio robusto para la gestión de comercios, transacciones, reportes financieros y auditoría, desarrollado en Go siguiendo principios de arquitectura limpia y patrones de diseño modernos.

## Descripción del Proyecto

Este sistema proporciona una API RESTful para administrar el ciclo de vida de comercios y sus operaciones financieras. Incluye un módulo de auditoría asíncrona para rastrear todos los cambios en el sistema y un motor de reportes para análisis de ingresos.

### Tecnologías Principales

- **Lenguaje**: Go 1.26+
- **Framework Web**: Gin Gonic
- **Base de Datos**: PostgreSQL 18
- **Documentación**: Swagger (Swaggo)
- **Contenedores**: Docker y Docker Compose
- **Logging**: Slog (JSON estructurado)

## Arquitectura del Sistema

El proyecto implementa el patrón **App Bootstrap** para la inicialización y gestión del ciclo de vida de la aplicación. La estructura de directorios sigue las convenciones de la comunidad Go:

- `cmd/`: Punto de entrada de la aplicación.
- `internal/`: Código privado del proyecto.
  - `app/`: Lógica de inicialización, inyección de dependencias y apagado gradual (graceful shutdown).
  - `domain/`: Entidades de negocio e interfaces de repositorios.
  - `usecases/`: Lógica de negocio y orquestación de servicios.
  - `infrastructure/`: Implementaciones concretas (Postgres, HTTP Handlers, Routers).
- `pkg/`: Utilidades compartidas y librerías auxiliares.
- `db/`: Scripts de base de datos y configuración de contenedores de persistencia.

## Requisitos Previos

- Docker y Docker Compose instalados.
- Go 1.26 o superior (para ejecución local).
- Make (opcional, recomendado para simplificar tareas).

## Configuración y Ejecución

### 1. Variables de Entorno

Copie el archivo de ejemplo o cree un archivo `.env` en la raíz del proyecto con la siguiente configuración:

```env
# Configuración de Base de Datos
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=commerce_db
DB_HOST=db
DB_PORT=5432

# Configuración de Aplicación
APP_PORT=8080

# Cadena de Conexión
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=disable
```

### 2. Ejecución con Docker Compose

La forma recomendada de desplegar el entorno completo es mediante Docker Compose. Esto levantará tanto la API como la base de datos con el esquema inicializado automáticamente.

```bash
docker compose up --build
```

### 3. Uso del Makefile

Se ha incluido un `Makefile` para automatizar las tareas comunes de desarrollo. Ejecute `make help` para ver todos los comandos disponibles.

| Comando | Descripción |
|---------|-------------|
| `make docker-up` | Levanta los contenedores en segundo plano. |
| `make docker-down` | Detiene y elimina contenedores y volúmenes. |
| `make build` | Compila el binario de la aplicación localmente. |
| `make test` | Ejecuta la suite de pruebas unitarias e integración. |
| `make swagger` | Regenera la documentación de API Swagger. |
| `make run` | Ejecuta la aplicación directamente desde el código fuente. |

## API y Documentación

Una vez iniciada la aplicación, la documentación interactiva de Swagger está disponible en:

`http://localhost:8080/swagger/index.html`

### Seguridad
El sistema utiliza el encabezado `X-User-Id` para identificar al actor de las operaciones y registrarlo en los logs de auditoría.

## Sistema de Auditoría (Bitácora)

El proyecto incluye un robusto sistema de auditoría transversal diseñado para registrar todas las operaciones críticas (Creación, Actualización, Eliminación) sin impactar el rendimiento de la API.

### Características Principales

1.  **Ejecución Asíncrona**: El registro de auditoría se realiza mediante goroutines. Esto permite que la respuesta al cliente sea inmediata, delegando el guardado en base de datos a un proceso en segundo plano.
2.  **Identificación del Actor**: La identidad del usuario (Actor) se extrae automáticamente del contexto de la petición, el cual es inyectado desde la capa de transporte (HTTP) mediante el encabezado `X-User-Id`.
3.  **Arquitectura Desacoplada**: El sistema utiliza una interfaz `AuditRepository` permitiendo que los casos de uso dependan de una abstracción. Esto facilita el cambio del motor de auditoría (ej. de DB a una cola de mensajes o log externo) sin modificar la lógica de negocio.
4.  **Trazabilidad Completa**: Cada registro de auditoría incluye:
    - `LogID`: Identificador único global del registro de auditoría.
    - `Action`: Descripción de la operación realizada (ej: "Crear Comercio").
    - `Actor`: Identificación del usuario que realizó la acción (extraído del encabezado).
    - `ResourceID`: Identificador de la entidad afectada (el ID del comercio, transacción, etc.).
    - `Timestamp`: Fecha y hora de la operación (formato RFC3339).

### Flujo de Trabajo
1. El middleware/handler extrae el `X-User-Id` y lo guarda en el `context.Context`.
2. El UseCase recibe el contexto y ejecuta la lógica de negocio.
3. Tras una operación exitosa, se dispara una goroutine que invoca al `AuditRepository` pasándole el contexto para recuperar al actor y persistir el log de auditoría.

## Logging Estructurado (Slog)

El sistema implementa logging estructurado en formato JSON utilizando la librería estándar `log/slog`. Esto facilita la integración con sistemas de observabilidad (Elasticsearch, CloudWatch, Datadog).

### Middleware de Logging

Un middleware personalizado en `pkg/util/logger` intercepta todas las peticiones HTTP y registra información técnica relevante:

- **Campos Registrados**: Método, ruta, parámetros de consulta (query), IP del cliente, latencia (en segundos), código de estado HTTP y User-Agent.
- **Niveles Automáticos**: 
  - Las peticiones con estado `>= 500` se registran con nivel `Error`.
  - El resto de las peticiones se registran con nivel `Info`.

### Ejemplo de Salida (JSON)
```json
{
  "time": "2026-03-06T03:23:44Z",
  "level": "INFO",
  "msg": "request completed",
  "request": {
    "method": "GET",
    "path": "/api/v1/comercios",
    "status": 200,
    "latency_seconds": 0.0012,
    "ip": "127.0.0.1"
  }
}
```

## Buenas Prácticas y Patrones

### Manejo de Valores Monetarios
Para garantizar la precisión financiera y evitar errores de redondeo inherentes a los tipos de punto flotante (`float64`), el proyecto utiliza una abstracción personalizada en `pkg/util/money`. 
- **Librería**: Basado en `github.com/cockroachdb/apd/v3` (Arbitrary Precision Decimals).
- **Redondeo**: Implementa el estándar IEEE 754 (Round Half Even).
- **Tipos**: Define tipos `Amount` y `Rate` para diferenciar entre montos monetarios y tasas porcentuales, protegiendo las operaciones aritméticas mediante un contexto matemático controlado.

### Estrategia de Identificadores (UID)
Inspirado por los patrones de diseño del repositorio oficial de **Kubernetes**, el proyecto utiliza un tipo especializado `types.UID` (`pkg/types/uid.go`).
- **Seguridad de Tipos**: En lugar de usar `string` nativos, se emplea un alias de tipo para representar identificadores únicos. Esto evita confusiones accidentales entre diferentes tipos de datos y mejora la legibilidad del código al indicar explícitamente que se espera un identificador único en las firmas de funciones.
- **Consistencia**: Centraliza la definición de lo que constituye un identificador en el sistema, facilitando futuras migraciones o cambios en el formato de los IDs sin afectar la lógica de negocio.

## Pruebas

El proyecto utiliza `testcontainers-go` para pruebas de integración reales con PostgreSQL. Asegúrese de tener Docker en ejecución antes de iniciar los tests.

```bash
make test
```
