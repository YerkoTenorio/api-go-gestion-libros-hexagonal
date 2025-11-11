package presentation

import (
	"api-go-gestion-libros-hexagonal/modules/book/application"
	"api-go-gestion-libros-hexagonal/modules/book/domain"
	"context"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookService application.BookServiceInterface
	validator   *validator.Validate
}

func NewBookHandler(bookService application.BookServiceInterface) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		validator:   validator.New(),
	}
}

// Helper functions para convertir entre DTOs y Domain
func domainToResponse(book *domain.Book) *BookResponse {
	return &BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Year:      book.Year,
		Genre:     book.Genre,
		ISBN:      book.ISBN,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
}

func requestToDomain(req CreateBookRequest) domain.UpdateBookInput {
	year := uint(req.Year)
	return domain.UpdateBookInput{
		Title:  &req.Title,
		Author: &req.Author,
		Year:   &year,
		Genre:  &req.Genre,
		ISBN:   &req.ISBN,
	}
}

func filterRequestToDomain(req BookFilterRequest) domain.BookFilter {
	return domain.BookFilter{
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
		Genre:  req.Genre,
	}
}

// HTTP Handlers
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var req CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	book, err := h.bookService.CreateBook(context.Background(), req.Title, req.Author, uint(req.Year), req.Genre, req.ISBN)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    domainToResponse(book),
		Message: "Book created successfully",
	})

}

func (h *BookHandler) GetBookById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{"Invalid book ID"},
		})
	}

	book, err := h.bookService.GetBookByID(context.Background(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	return c.JSON(Response{
		Success: true,
		Data:    domainToResponse(book),
	})
}

func (h *BookHandler) GetBookByISBN(c *fiber.Ctx) error {
	// Aqui lo que hacemos es obtener el isbn que viene como string
	isbn := c.Params("isbn")
	// Aqui lo que hacemos es obtener el libro por isbn
	book, err := h.bookService.GetBookByISBN(context.Background(), isbn)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	return c.JSON(Response{
		Success: true,
		Data:    domainToResponse(book),
	})
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	// Aqui lo que hacemos es obtener el id que viene como string
	idParam := c.Params("id")

	// Aqui lo que hacemos es convertir el id que viene como string a uint
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	// Aqui lo que hacemos es obtener el body que viene como json
	var req UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	// Aqui lo que hacemos es convertir el body que viene como json a domain.UpdateBookInput
	input := domain.UpdateBookInput{
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
		Genre:  req.Genre,
		ISBN:   req.ISBN,
	}

	// Aqui lo que hacemos es actualizar el libro
	book, err := h.bookService.UpdateBook(context.Background(), uint(id), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	// Aqui lo que hacemos es devolver el libro actualizado
	return c.JSON(Response{
		Success: true,
		Data:    domainToResponse(book),
		Message: "book updated successfully",
	})

}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	if err := h.bookService.DeleteBook(context.Background(), uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	return c.JSON(Response{
		Success: true,
		Message: "book deleted successfully",
	})

}

func (h *BookHandler) SearchBooks(c *fiber.Ctx) error {
	// Aqui lo que hacemos es obtener el body que viene como json
	var req BookFilterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}
	// Aqui lo que hacemos es convertir el body que viene como json a domain.BookFilter
	filter := filterRequestToDomain(req)

	books, err := h.bookService.SearchBooks(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}
	// Convertir a response
	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *domainToResponse(book)
	}

	return c.JSON(Response{
		Success: true,
		Data:    responses,
	})
}

func (h *BookHandler) GetBookByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{"Invalid book ID"},
		})
	}

	book, err := h.bookService.GetBookByID(context.Background(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	return c.JSON(Response{
		Success: true,
		Data:    domainToResponse(book),
	})
}

func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	books, err := h.bookService.SearchBooks(context.Background(), domain.BookFilter{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Success: false,
			Errors:  []string{err.Error()},
		})
	}

	// Convertir a response
	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *domainToResponse(book)
	}

	return c.JSON(Response{
		Success: true,
		Data:    responses,
	})

}
