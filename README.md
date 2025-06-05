# Prueba Técnica Go Guarapo

API REST para autenticación y gestión de tareas, desarrollada en Go con Gin, GORM y SQLite. Incluye autenticación por token, documentación Swagger y despliegue con Docker.

---

## Características

- **CRUD de tareas** por usuario autenticado
- **Autenticación** con token (header `Authorization: Bearer <token>`)
- **Persistencia** con SQLite (usando GORM)
- **Documentación interactiva** con Swagger (OpenAPI)
- **Tests unitarios y de integración** con mocks y base en memoria
- **Despliegue fácil** con Docker y Docker Compose
- **Control de Versiones** con gitflow y git

---

## Requisitos

- Go 1.20+ (para desarrollo local)
- Docker y Docker Compose (para despliegue)

---

## Instalación y uso local

1. **Clona el repositorio:**

   ```sh
   git clone <url-del-repo>
   cd prueba-tecnica-go-guarapo
   ```

2. **En caso de ser necesario Instala dependencias:**

   ```sh
   go mod download
   ```

3. **En caso de ser necesario Genera la documentación Swagger:**

   ```sh
   swag init --generalInfo api/server/server.go --output api/docs
   ```

4. **Crea el archivo `.env` si lo necesitas (por defecto esta el 8080 asi que puedes iniciarlo asi):**

   ```env
   PORT=8080
   ```

5. **Ejecuta la API:**

   ```sh
   go run ./api/cmd/main.go
   ```

6. **Accede a Swagger UI:**

   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Uso con Docker

1. **Construye y levanta los servicios:**

   ```sh
   docker-compose up --build
   ```

2. **La API estará disponible en:**

   [http://localhost:8080](http://localhost:8080)

3. **Swagger UI:**

   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Autenticación

- Todos los endpoints de `/api/tasks` requieren autenticación.
- Usa el endpoint `/api/login` para obtener un token:
  ```json
  POST /api/login
  {
    "username": "usuario"
  }
  ```
- Usa el token en el header:
  ```
  Authorization: Bearer <token>
  ```

---

## Endpoints principales

- `POST   /api/login` — Login de usuario (devuelve token)
- `GET    /api/tasks` — Listar tareas del usuario autenticado
- `POST   /api/tasks` — Crear tarea
- `GET    /api/tasks/{id}` — Obtener tarea por ID
- `PUT    /api/tasks/{id}` — Actualizar tarea
- `DELETE /api/tasks/{id}` — Eliminar tarea

Consulta la [documentación Swagger](http://localhost:8080/swagger/index.html) para detalles y ejemplos.

---

## Pruebas

- Ejecuta todos los tests:
  ```sh
  go test ./...
  ```
  - Ejecuta todos los tests con el porcentaje de cover:
  ```sh
  go test ./... -cover
  ```

- Los tests incluyen:
  - Unitarios con mocks (handlers y servicios)
  - Integración con SQLite en memoria

---

## Estructura del proyecto

```
api/
  cmd/           # main.go (entrypoint)
  docs/          # Documentación Swagger generada
  handlers/      # Handlers de Gin (auth, task)
  models/        # Modelos de datos
  server/        # Inicialización del servidor y rutas
  services/      # Lógica de negocio (auth, task)
Dockerfile
docker-compose.yml
.env
README.md
```

---

## Notas

- El token debe enviarse como: `Authorization: Bearer <token>` en el swagger es necesario que coloques Bearer <pegas el token>
- La base de datos SQLite se crea automáticamente en el contenedor.
- Para desarrollo, puedes borrar el archivo `tasks.db` y se recreará vacío.

---
