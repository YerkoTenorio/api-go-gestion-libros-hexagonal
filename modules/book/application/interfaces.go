package application

import (
	"api-go-gestion-libros-hexagonal/modules/book/domain"
	"context"
)

// BookService define los casos de uso de libros expuestos a la capa de presentacion

type BookServiceInterface interface {
	CreateBook(ctx context.Context, title, author string, year uint, genre, isbn string) (*domain.Book, error) // Crea un nuevo libro
	UpdateBook(ctx context.Context, id uint, input domain.UpdateBookInput) (*domain.Book, error)               // Actualiza un libro existente
	DeleteBook(ctx context.Context, id uint) error                                                             // Elimina un libro por ID
	GetBookByID(ctx context.Context, id uint) (*domain.Book, error)                                            // Obtiene un libro por ID
	GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error)                                      // Obtiene un libro por ISBN
	SearchBooks(ctx context.Context, filter domain.BookFilter) ([]*domain.Book, error)                         // Busca libros por filtro
}
