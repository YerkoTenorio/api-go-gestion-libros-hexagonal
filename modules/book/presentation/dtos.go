package presentation

import "time"

// CreateBookRequest define la estructura para crear un libro via API

type CreateBookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Year   int    `json:"year" validate:"required,min=1450"`
	Genre  string `json:"genre"`
	ISBN   string `json:"isbn" validate:"required"`
}

//UpdateBookRequest defune la estructua para actualizar un libro via API
// usan punteros para que los campos sean opcionales
type UpdateBookRequest struct {
	Title  *string `json:"title,omitempty" validate:"omitempty,min=1"`
	Author *string `json:"author,omitempty" validate:"omitempty,min=1"`
	Year   *uint   `json:"year,omitempty" validate:"omitempty,min=1450"`
	Genre  *string `json:"genre,omitempty" validate:"omitempty,min=1"`
	ISBN   *string `json:"isbn,omitempty" validate:"omitempty,min=1"`
}

//BookResponse define la estructura para devolver un libro via API
type BookResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      uint      `json:"year"`
	Genre     string    `json:"genre"`
	ISBN      string    `json:"isbn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//BookFilterRequest define la estructura para filtrar libros via API
type BookFilterRequest struct {
	Title  *string `json:"title,omitempty"`
	Author *string `json:"author,omitempty"`
	Year   *uint   `json:"year,omitempty"`
	Genre  *string `json:"genre,omitempty"`
}

// Response estándar para todas las APIs
//Data es una interface{} para que pueda devolver cualquier tipo de dato
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse para errores detallados
type ErrorResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}
