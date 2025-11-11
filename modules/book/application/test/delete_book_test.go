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

func TestBookService_DeleteBook_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	// preparar el contexto y datos de prueba
	ctx := context.Background()
	id := uint(1)

	// mock para verificar existencia
	mockRepo.EXPECT().GetByID(ctx, id).Return(&domain.Book{}, nil)

	// mock para eliminar
	mockRepo.EXPECT().Delete(ctx, id).Return(nil)

	// Act
	err := service.DeleteBook(ctx, id)

	// Assert
	assert.NoError(t, err)
}

func TestBookService_DeleteBook_BookNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	// preparar el contexto
	ctx := context.Background()
	id := uint(1)

	// mock para verificar existencia
	mockRepo.EXPECT().GetByID(ctx, id).Return(nil, fmt.Errorf("not found"))

	// Act
	err := service.DeleteBook(ctx, id)

	// Assert
	assert.Error(t, err)                         // error esperado
	assert.Contains(t, err.Error(), "not found") // mensaje de error esperado

}

func TestBookService_DeleteBook_ZeroID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	service := application.NewBookService(mockRepo)

	// preparar el contexto
	ctx := context.Background()
	id := uint(0)

	// Act
	err := service.DeleteBook(ctx, id)

	// Assert
	assert.Error(t, err)                              // error esperado
	assert.Contains(t, err.Error(), "id is required") // mensaje de error esperado

}
