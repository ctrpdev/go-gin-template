# Go API Template - Clean Architecture

Esta es una plantilla lista para producción (_production-ready_) diseñada para construir APIs REST en Go siguiendo principios de **Clean Architecture** y las mejores prácticas de la industria.

El proyecto está diseñado para funcionar con **PostgreSQL** para la persistencia de datos (usando la generación de código automática de SQLC) y **Redis** para el manejo de sesiones y caché. Todo el flujo de operaciones (DevOps, bases de datos y configuraciones) ha sido dockerizado para una experiencia de desarrollo libre de fricciones.

## Características Principales

*   **Clean Architecture:** Estructura de carpetas modular (`cmd`, `internal`, `domain`, `handler`, etc.) que separa las reglas de negocio de la infraestructura.
*   **Hot-Reloading en Desarrollo:** Integración con [Air](https://github.com/air-verse/air) para reconstrucción instantánea del código mientras desarrollas.
*   **Gestión de Bases de Datos:** 
    *   **PostgreSQL** usando `jackc/pgx/v5` para alto rendimiento.
    *   **SQLC** para generar código Go fuertemente tipado directamente desde las consultas SQL puras.
    *   Migraciones manejadas mediante `golang-migrate`.
*   **Redis:** Integrado nativamente para manejo de sesiones o almacenamiento temporal.
*   **Capa de Seguridad:** Middleware de Autenticación con JWT de fábrica.
*   **Configuración Atómica:** Manejo granular de variables de entorno mediante `viper`.
*   **Infraestructura Aislada:** Todos los manifiestos de contenedores y rutinas de CI/CD viven en la carpeta `/deployment`.

---

## Pre-requisitos (Lo que debes hacer antes de iniciar)

Antes de levantar el proyecto, asegúrate de seguir estos pasos:

1. **Instalar Docker y Docker Compose:** Todo el flujo del proyecto funciona principalmente a través de contenedores.
2. **(Opcional pero recomendado) Instalar Make:** Usamos un archivo `Makefile` para crear atajos a comandos largos.
3. **Configurar Variables de Entorno:**
   El proyecto incluye un archivo de base. Debes crear tu archivo `.env` configurando los valores de tus entornos.

   Renombra o copia el archivo de ejemplo:
   ```bash
   cp .env.example .env
   ```

   _Nota: Los valores por defecto de `.env.example` incluyen conectarse a puertos separados para la DB (`5433`) y Redis (`6380`) para evitar colisiones con bases de datos locales alojadas en tu máquina host._

---

## Despliegue del Proyecto

La infraestructura Docker está dividida en archivos separados que se encuentran en el directorio `deployment/`. Tienes dos modos principales para ejecutar el proyecto.

### 1. Entorno de Desarrollo (Hot-Reload)

Este entorno monta tus archivos locales directamente dentro del contenedor de Golang y usa `Air` para recompilar automáticamente la aplicación cuando haya cambios en el código (sin necesidad de reiniciar todo el contenedor de forma manual).

A través de **Make**:
```bash
make docker-dev
```

A través de **Docker puro** (desde la raíz):
```bash
docker compose -f deployment/docker-compose.yml -f deployment/docker-compose.dev.yml up -d --build
```

### 2. Entorno de Producción

Este entorno compila la aplicación usando un _multi-stage build_ optimizado, descartando el código fuente y generando una imagen super liviana basada en Alpine. Es segura, inmutable y apta para entornos Cloud.

A través de **Make**:
```bash
make docker-prod
```

A través de **Docker puro** (desde la raíz):
```bash
docker compose -f deployment/docker-compose.yml up -d --build
```

### Detener los servicios

Para dar de baja los contenedores sin borrar los volúmenes de las bases de datos:

```bash
make docker-down
# o: docker compose -f deployment/docker-compose.yml down
```

---

## Estructura del Proyecto

*   `cmd/api/main.go`: El punto de entrada principal (_entrypoint_) de la aplicación.
*   `deployment/`: Contiene los archivos de infraestructura y CI/CD (`Dockerfile`, `docker-compose.yml`, `docker-compose.dev.yml`).
*   `internal/`: Lógica privada de la aplicación. Se subdivide en la Arquitectura Limpia (`domain`, `handler`, `service`, `repository`, `config`).
*   `migrations/`: Archivos `.sql` up/down usados para definir el esquema de la base de datos PostgreSQL.
*   `sqlc/`: Consultas puras en SQL usadas por `sqlc` para autogenerar la lógica y los modelos en el repositorio (localizado en `internal/repository/postgres/db`).

## Herramientas de Desarrollo Comunes

Si tienes `Make` y las herramientas locales (`sqlc`, `golang-migrate`) instaladas en tu equipo host:

*   `make sqlc`: Genera de nuevo el código Go a partir de tus consultas de la carpeta `sqlc/`.
*   `make lint`: Aplica formato estándar a tu código e instala limpiamente los módulos en caso de ser necesario.
*   `make migrate-up`: Ejecuta las migraciones más recientes sobre tu base de datos (PostgreSQL).
*   `make migrate-down`: Revierte las migraciones en caso de error o ajuste.
*   `make run`: (Sólo local) Compila localmente a la carpeta `build/` y ejecuta el programa directamente en tu SO Windows/Linux/Mac ignorando Docker.
