package application_test

import (
	"api-go-gestion-libros-hexagonal/modules/book/application"
	"api-go-gestion-libros-hexagonal/modules/book/application/mocks"
	"api-go-gestion-libros-hexagonal/modules/book/domain"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBookService_UpdateBook_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()
	id := uint(1)
	newTitle := "Título Actualizado"
	newAuthor := "Autor Actualizado"
	newYear := uint(2023)

	// Libro actual que existe en BD
	currentBook := &domain.Book{
		ID:        id,
		Title:     "Título Original",
		Author:    "Autor Original",
		Year:      2020,
		Genre:     "Ficción",
		ISBN:      "9788418037016",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Libro actualizado que retornará el repositorio
	updatedBook := &domain.Book{
		ID:        id,
		Title:     newTitle,
		Author:    newAuthor,
		Year:      newYear,
		Genre:     "Ficción",
		ISBN:      "9788418037016",
		CreatedAt: currentBook.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}

	input := domain.UpdateBookInput{
		Title:  &newTitle,
		Author: &newAuthor,
		Year:   &newYear,
	}

	// Mock: obtener libro actual
	mockRepo.EXPECT().GetByID(ctx, id).Return(currentBook, nil)

	// Mock: actualizar libro
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(updatedBook, nil)

	// Act
	result, err := service.UpdateBook(ctx, id, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newTitle, result.Title)
	assert.Equal(t, newAuthor, result.Author)
	assert.Equal(t, newYear, result.Year)
}

func TestBookService_UpdateBook_ZeroID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()
	input := domain.UpdateBookInput{
		Title: stringPtr("Nuevo Título"),
	}

	// Act
	result, err := service.UpdateBook(ctx, 0, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "id is required")
}

func TestBookService_UpdateBook_BookNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()
	id := uint(999)
	input := domain.UpdateBookInput{
		Title: stringPtr("Nuevo Título"),
	}

	// Mock: libro no encontrado
	mockRepo.EXPECT().GetByID(ctx, id).Return(nil, fmt.Errorf("not found"))

	// Act
	result, err := service.UpdateBook(ctx, id, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestBookService_UpdateBook_ISBNAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()
	id := uint(1)
	existingISBN := "9788418037017"

	currentBook := &domain.Book{
		ID:   id,
		ISBN: "9788418037016",
	}

	otherBook := &domain.Book{
		ID:   2,
		ISBN: existingISBN,
	}

	input := domain.UpdateBookInput{
		ISBN: &existingISBN,
	}

	// Mock: obtener libro actual
	mockRepo.EXPECT().GetByID(ctx, id).Return(currentBook, nil)

	// Mock: ISBN ya existe en otro libro
	mockRepo.EXPECT().GetByISBN(ctx, existingISBN).Return(otherBook, nil)

	// Act
	result, err := service.UpdateBook(ctx, id, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "isbn already registered by another book")
}

func TestBookService_UpdateBook_SameISBN(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	ctx := context.Background()
	id := uint(1)
	sameISBN := "9788418037016"
	title := "Libro de prueba"
	author := "Autor de prueba"
	year := uint(2022)
	genre := "Ficción"

	currentBook := &domain.Book{
		ID:   id,
		ISBN: sameISBN,
	}

	updatedBook := &domain.Book{
		ID:     id,
		ISBN:   sameISBN,
		Title:  title,
		Author: author,
		Year:   year,
		Genre:  genre,
	}

	input := domain.UpdateBookInput{
		Title:  &title,
		Author: &author,
		Year:   &year,
		Genre:  &genre,
		ISBN:   &sameISBN,
	}

	// Mock: obtener libro actual
	mockRepo.EXPECT().GetByID(ctx, id).Return(currentBook, nil)

	// Mock: verificar ISBN (retorna el mismo libro, indicando que es el mismo ISBN)
	mockRepo.EXPECT().GetByISBN(ctx, sameISBN).Return(currentBook, nil)

	// Mock: actualizar libro
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(updatedBook, nil)

	// Act
	result, err := service.UpdateBook(ctx, id, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sameISBN, result.ISBN)
}

// Helper function para crear punteros a string
func stringPtr(s string) *string {
	return &s
}
