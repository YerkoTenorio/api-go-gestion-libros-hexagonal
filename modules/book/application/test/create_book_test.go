package application_test

import (
	"api-go-gestion-libros-hexagonal/modules/book/application"
	"api-go-gestion-libros-hexagonal/modules/book/application/mocks"
	"api-go-gestion-libros-hexagonal/modules/book/domain"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// estructura dek book

//type Book struct {
//ID        uint      `json:"id" db:"id"`
//Title     string    `json:"title" db:"title"`
//Author    string    `json:"author" db:"author"`
//Year      uint      `json:"year" db:"year"`
//Genre     string    `json:"genre" db:"genre"`
//ISBN      string    `json:"isbn" db:"isbn"`
//CreatedAt time.Time `json:"created_at" db:"created_at"`
//UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
//}

func TestBookService_CreateBook_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Crear mock del repositorio
	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	// Prepara el contexto y datos de prueba
	ctx := context.Background()
	title := "Libro de prueba"
	author := "Autor de prueba"
	year := uint(2022)
	genre := "Ficción"
	isbn := "978-84-18037-01-6"
	normalizedISBN := domain.NormalizeISBN(isbn)

	// Mock para verificar que no existe el ISBN
	mockRepo.EXPECT().GetByISBN(ctx, normalizedISBN).Return(nil, fmt.Errorf("not found"))

	// Mock para la creación del libro
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	// Act
	result, err := service.CreateBook(ctx, title, author, year, genre, normalizedISBN)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, title, result.Title)
	assert.Equal(t, author, result.Author)
	assert.Equal(t, year, result.Year)
	assert.Equal(t, genre, result.Genre)
	assert.Equal(t, normalizedISBN, result.ISBN)
}

func TestBookService_CreateBook_TitleRequired(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()

	// Act
	result, err := service.CreateBook(ctx, "", "Autor", 2022, "Ficción", "978-84-18037-01-6")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "title is required")
}

func TestBookService_CreateBook_AuthorRequired(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()

	// Act
	result, err := service.CreateBook(ctx, "Título", "", 2022, "Ficción", "978-84-18037-01-6")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "author is required")
}

func TestBookService_CreateBook_ISBNAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()
	isbn := "978-84-18037-01-6"

	// Mock: ISBN ya existe
	existingBook := &domain.Book{ID: 1, Title: "Otro libro", ISBN: domain.NormalizeISBN(isbn)}
	mockRepo.EXPECT().GetByISBN(ctx, domain.NormalizeISBN(isbn)).Return(existingBook, nil)

	// Act
	result, err := service.CreateBook(ctx, "Nuevo Título", "Nuevo Autor", 2023, "Ficción", domain.NormalizeISBN(isbn))

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
}
