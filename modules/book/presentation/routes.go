package presentation

import "github.com/gofiber/fiber/v2"

func SetupBookRoutes(app *fiber.App, handler *BookHandler) {
	api := app.Group("/api/v1/books")

	// CRUD endpoints
	api.Post("/", handler.CreateBook)             // POST /api/v1/books
	api.Get("/", handler.GetAllBooks)             // GET /api/v1/books
	api.Get("/search", handler.SearchBooks)       // GET /api/v1/books/search?title=...&author=...
	api.Get("/:id", handler.GetBookByID)          // GET /api/v1/books/123
	api.Get("/isbn/:isbn", handler.GetBookByISBN) // GET /api/v1/books/isbn/978-3-16-148410-0
	api.Put("/:id", handler.UpdateBook)           // PUT /api/v1/books/123
	api.Delete("/:id", handler.DeleteBook)        // DELETE /api/v1/books/123
}
