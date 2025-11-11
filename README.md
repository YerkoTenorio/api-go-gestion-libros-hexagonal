# API Go Gestión de Libros - Arquitectura Hexagonal

Una API RESTful para la gestión de libros implementada en Go siguiendo el patrón de arquitectura hexagonal (Ports and Adapters). Este proyecto demuestra las mejores prácticas de desarrollo de software con una clara separación de responsabilidades.

## 🏗️ Arquitectura

El proyecto sigue la arquitectura hexagonal con los siguientes componentes:

```
api-go-gestion-libros-hexagonal/
├── cmd/server/           # Punto de entrada de la aplicación
├── modules/
│   └── book/
│       ├── domain/       # Lógica de negocio y entidades
│       ├── application/  # Casos de uso y servicios
│       ├── infrastructure/ # Implementación de repositorios
│       └── presentation/ # Handlers HTTP y DTOs
├── shared/
│   ├── config/          # Configuración de la aplicación
│   └── database/        # Conexión a base de datos
└── docs/               # Documentación adicional
```

### Capas de la Arquitectura

- **Domain**: Contiene las entidades de negocio y reglas de validación
- **Application**: Implementa los casos de uso y lógica de aplicación
- **Infrastructure**: Proporciona implementaciones concretas (base de datos)
- **Presentation**: Maneja las peticiones HTTP y respuestas API

## 🚀 Características

- ✅ **CRUD completo** para gestión de libros
- ✅ **Validación de ISBN-10 e ISBN-13** con verificación de checksum
- ✅ **Búsqueda avanzada** por título, autor, año y género
- ✅ **Arquitectura hexagonal** para mejor mantenibilidad
- ✅ **Testing unitario** con mocks
- ✅ **Middleware** para logging, recuperación y CORS
- ✅ **Base de datos Turso** (SQLite en la nube)
- ✅ **Validación de entrada** con struct tags

## 📋 Prerrequisitos

- Go 1.25.0 o superior
- Base de datos Turso (o SQLite compatible)
- Git

## 🛠️ Instalación

1. **Clonar el repositorio**
   ```bash
   git clone <repository-url>
   cd api-go-gestion-libros-hexagonal
   ```

2. **Instalar dependencias**
   ```bash
   go mod download
   ```

3. **Configurar variables de entorno**
   ```bash
   # Crear archivo .env
   cp .env.example .env
   ```

   Editar el archivo `.env` con tus credenciales:
   ```env
   TURSO_DATABASE_URL=tu-database-url.turso.io
   TURSO_AUTH_TOKEN=tu-auth-token
   PORT=8080
   ```

4. **Ejecutar la aplicación**
   ```bash
   go run cmd/server/main.go
   ```

   O construir y ejecutar:
   ```bash
   go build -o api cmd/server/main.go
   ./api
   ```

## 📡 Endpoints de la API

### Base URL
```
http://localhost:8080/api/v1/books
```

### Endpoints Disponibles

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Health check de la API |
| POST | `/books` | Crear un nuevo libro |
| GET | `/books` | Obtener todos los libros |
| GET | `/books/search` | Buscar libros por filtros |
| GET | `/books/:id` | Obtener libro por ID |
| GET | `/books/isbn/:isbn` | Obtener libro por ISBN |
| PUT | `/books/:id` | Actualizar libro existente |
| DELETE | `/books/:id` | Eliminar libro |

### Ejemplos de Uso

#### Crear un Libro
```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Cien Años de Soledad",
    "author": "Gabriel García Márquez",
    "year": 1967,
    "genre": "Realismo Mágico",
    "isbn": "978-0-06-088328-7"
  }'
```

#### Buscar Libros
```bash
curl -X POST http://localhost:8080/api/v1/books/search \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Gabriel García Márquez",
    "year": 1967
  }'
```

#### Actualizar un Libro
```bash
curl -X PUT http://localhost:8080/api/v1/books/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Cien Años de Soledad (Edición Especial)",
    "genre": "Realismo Mágico, Clásico"
  }'
```

## 📊 Modelo de Datos

### Book
```json
{
  "id": 1,
  "title": "Cien Años de Soledad",
  "author": "Gabriel García Márquez",
  "year": 1967,
  "genre": "Realismo Mágico",
  "isbn": "9780060883287",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Validaciones
- **Title**: Requerido, no puede estar vacío
- **Author**: Requerido, no puede estar vacío
- **Year**: Requerido, debe estar entre 1450 y el año actual
- **ISBN**: Requerido, debe ser ISBN-10 o ISBN-13 válido
- **Genre**: Opcional

## 🧪 Testing

### Ejecutar Tests
```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Ejecutar tests en un módulo específico
go test ./modules/book/application/test/
```

### Generar Mocks
```bash
# Generar mocks para interfaces
go generate ./...
```

## 🔧 Tecnologías Utilizadas

- **Go 1.25.0**: Lenguaje de programación principal
- **Fiber v2**: Framework web rápido y minimalista
- **Turso**: Base de datos SQLite en la nube
- **Validator**: Validación de structs
- **Testify**: Framework de testing
- **Go Mock**: Generación de mocks

## 📁 Estructura del Proyecto

```
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada
├── modules/
│   └── book/
│       ├── domain/
│       │   ├── book.go          # Entidad Book y lógica de negocio
│       │   └── repository.go    # Interfaces de repositorio
│       ├── application/
│       │   ├── service.go       # Servicios de aplicación
│       │   ├── interfaces.go    # Interfaces de servicio
│       │   ├── mocks/           # Mocks generados
│       │   └── test/            # Tests unitarios
│       ├── infrastructure/
│       │   └── sql_repository.go # Implementación SQL
│       └── presentation/
│           ├── handlers.go      # HTTP handlers
│           ├── routes.go        # Definición de rutas
│           └── dtos.go          # Data Transfer Objects
├── shared/
│   ├── config/
│   │   └── config.go           # Manejo de configuración
│   └── database/
│       └── turso.go            # Conexión a Turso
├── go.mod                       # Módulos de Go
├── go.sum                       # Checksum de dependencias
├── .env                         # Variables de entorno
└── README.md                    # Este archivo
```

## 🎯 Principios de Diseño

- **Separación de responsabilidades**: Cada capa tiene una responsabilidad clara
- **Inversión de dependencias**: La capa de dominio no depende de infraestructura
- **Single Responsibility**: Cada componente tiene una única razón de cambiar
- **Open/Closed**: Abierto para extensión, cerrado para modificación
- **Dependency Injection**: Las dependencias se inyectan, no se crean internamente

## 🤝 Contribuir

1. Fork del proyecto
2. Crear una feature branch (`git checkout -b feature/amazing-feature`)
3. Commit de los cambios (`git commit -m 'Add amazing feature'`)
4. Push a la branch (`git push origin feature/amazing-feature`)
5. Abrir un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.


**Nota**: Este proyecto es un ejemplo educativo de arquitectura hexagonal en Go. Puede ser utilizado como base para proyectos más complejos.
