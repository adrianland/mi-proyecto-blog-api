# Blog API - Microservicio REST en Go

API REST desarrollada en Go para la gestión de Autores y Artículos.

El objetivo de este README es que cualquier persona pueda:

* Levantar el proyecto rápidamente
* Probar los endpoints
* Ejecutar los tests

---

# 🚀 Cómo ejecutar el proyecto

## Requisitos

* Docker
* Docker Compose
* (Opcional) Go 1.21+ si se quiere ejecutar sin Docker

---

## ▶️ Ejecutar con Docker (recomendado)

### 1. Construir las imágenes

```bash
docker-compose build --no-cache
```

### 2. Levantar los contenedores

```bash
docker-compose up
```

La API quedará disponible en:

```
http://localhost:8080
```

### 3. Verificar que el servicio está funcionando

```bash
curl http://localhost:8080/health
```

---

# 🧪 Ejecutar los tests

Desde la raíz del proyecto:

## Tests unitarios

```bash
go test ./tests/unit -v
```

## Tests de integración

```bash
go test ./tests/integration -v
```

---

# 📌 Endpoints disponibles

## AUTORES

POST   /autores                    → Crear autor
GET    /autores/{id}/resumen       → Resumen del autor
GET    /autores/{id}/articulos     → Artículos del autor
GET    /autores/top?n=3            → Top N autores por score

---

## ARTÍCULOS

POST   /articulos                  → Crear artículo (BORRADOR)
GET    /articulos                  → Listar artículos publicados (paginado)
GET    /articulos/{id}             → Obtener artículo
PUT    /articulos/{id}             → Editar artículo (solo BORRADOR)
PUT    /articulos/{id}/publicar    → Publicar artículo (con validaciones)

---

## SISTEMA

GET    /health                     → Health check

---

# 📄 Variables de entorno

El proyecto utiliza un archivo `.env` incluido en el repositorio.

Variables principales:

```
SERVER_PORT=8080

DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=

RATE_LIMIT_REQUESTS=
RATE_LIMIT_WINDOW=
```

Si necesitas cambiar el puerto o credenciales de base de datos, modifica el archivo `.env` antes de ejecutar `docker-compose up`.

---

# 📌 Reglas de negocio importantes

Estados del artículo:

* BORRADOR
* PUBLICADO

Validaciones al publicar:

* Mínimo 120 palabras
* Máximo 35% de palabras repetidas
* El score se calcula automáticamente al publicar

---

# 🗂 Estructura básica del proyecto

```
cmd/                → Punto de entrada
internal/           → Lógica de negocio y dominio
interfaces/         → Handlers HTTP
tests/              → Tests unitarios e integración
Dockerfile
docker-compose.yml
```

---

# ✅ Resumen rápido para probar

```bash
docker-compose build --no-cache
docker-compose up
```

Luego probar endpoints en:

```
http://localhost:8080
```

Y ejecutar tests con:

```bash
go test ./tests/unit -v
go test ./tests/integration -v
```
