package main

import (
	"api-go-gestion-libros-hexagonal/modules/book/application"
	"api-go-gestion-libros-hexagonal/modules/book/infrastructure"
	"api-go-gestion-libros-hexagonal/modules/book/presentation"
	"api-go-gestion-libros-hexagonal/shared/config"
	"api-go-gestion-libros-hexagonal/shared/database"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Conectar a base de datos
	db, err := database.OpenTurso(cfg.TursoURL, cfg.TursoToken)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Inicializar schema
	if err := database.InitSchema(db); err != nil {
		log.Fatal("Error initializing schema:", err)
	}

	// Crear instancias de la arquitectura hexagonal
	// Infrastructure -> Application -> Presentation
	bookRepo := infrastructure.NewSqlBookRepository(db)
	bookService := application.NewBookService(bookRepo)
	bookHandler := presentation.NewBookHandler(bookService)

	// Configurar Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(presentation.ErrorResponse{
				Success: false,
				Errors:  []string{err.Error()},
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(presentation.Response{
			Success: true,
			Message: "API is running",
		})
	})

	// Setup routes
	presentation.SetupBookRoutes(app, bookHandler)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + strconv.Itoa(cfg.Port)))
}
