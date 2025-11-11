package application

import (
	"api-go-gestion-libros-hexagonal/modules/book/domain"
	"context"
	"fmt"
)

type BookService struct {
	bookRepo domain.BookRepository
}

func NewBookService(bookRepo domain.BookRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
	}
}

func (s *BookService) CreateBook(ctx context.Context, title, author string, year uint, genre, isbn string) (*domain.Book, error) {
	// Validaciones basicas por caso de uso (evita llamadas innecesarias al repositorio)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if author == "" {
		return nil, fmt.Errorf("author is required")
	}

	// Construir entidad y aplicar reglas de dominio
	book := domain.NewBook(title, author, year, genre, isbn)
	if err := book.ValidateBasic(); err != nil {
		return nil, err
	}

	// Unicidad por ISBN
	existing, err := s.bookRepo.GetByISBN(ctx, book.ISBN)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("isbn %s already exists", book.ISBN)
	}

	// Persistir en el repositorio
	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}
	return book, nil

}

func (s *BookService) UpdateBook(ctx context.Context, id uint, input domain.UpdateBookInput) (*domain.Book, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	// Traer actual
	current, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("book %d not found: %w", id, err)
	}

	// Validar input opcional
	if err := domain.ValidateUpdateInput(input); err != nil {
		return nil, err
	}
	// Si cambia ISBN, verificar unicidad
	if input.ISBN != nil {
		if other, err := s.bookRepo.GetByISBN(ctx, *input.ISBN); err == nil && other != nil && other.ID != current.ID {
			return nil, fmt.Errorf("isbn already registered by another book")
		}
	}
	// Aplicar cambios y revalidar
	current.Update(input)
	if err := current.ValidateBasic(); err != nil {
		return nil, err
	}
	// Persistir
	updated, err := s.bookRepo.Update(ctx, current)
	if err != nil {
		return nil, err
	}
	return updated, nil

}

// DeleteBook elimina un libro por ID
func (s *BookService) DeleteBook(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}
	// Verificar existencia (opcional, util para 404)
	if _, err := s.bookRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("book %d not found: %w", id, err)
	}
	return s.bookRepo.Delete(ctx, id)
}

func (s *BookService) GetBookByID(ctx context.Context, id uint) (*domain.Book, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookService) GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	if isbn == "" {
		return nil, fmt.Errorf("isbn is required")
	}
	// Normalizar ISBN antes de consultar
	isbn = domain.NormalizeISBN(isbn)
	return s.bookRepo.GetByISBN(ctx, isbn)
}

// SearchBooks delega la busqueda al repositorio con filtros
func (s *BookService) SearchBooks(ctx context.Context, filter domain.BookFilter) ([]*domain.Book, error) {
	if !filter.HasAny() {
		return s.bookRepo.GetAll(ctx)
	}
	return s.bookRepo.FindByFilter(ctx, filter)
}
