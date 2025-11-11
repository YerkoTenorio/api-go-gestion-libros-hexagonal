package domain

import "context"

//go:generate mockgen -source=repository.go -destination=../application/mocks/mock_book_repository.go -package=mocks

// Crearemos interfaces contratos que describan las operaciones del repositorio (obtener todos, buscar por filtros, crear, actualizar, eliminar).
type BookRepository interface {
	// Create crea un nuevo libro en el repositorio
	Create(ctx context.Context, book *Book) error
	// Update actualiza un libro existente en el repositorio
	Update(ctx context.Context, book *Book) (*Book, error)
	// Delete elimina un libro existente en el repositorio
	Delete(ctx context.Context, id uint) error
	// GetAll obtiene todos los libros del repositorio
	GetAll(ctx context.Context) ([]*Book, error)
	// GetByID obtiene un libro por su ID del repositorio
	GetByID(ctx context.Context, id uint) (*Book, error)
	// FindByFilter obtiene libros por filtros del repositorio
	FindByFilter(ctx context.Context, filter BookFilter) ([]*Book, error)
	// GetByISBN obtiene libros por ISBN del repositorio
	GetByISBN(ctx context.Context, isbn string) (*Book, error)
}
