package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type Book struct {
	ID        uint      `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Year      uint      `json:"year" db:"year"`
	Genre     string    `json:"genre" db:"genre"`
	ISBN      string    `json:"isbn" db:"isbn"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BookFilter struct {
	Title  *string `json:"title"`
	Author *string `json:"author"`
	Year   *uint   `json:"year"`
	Genre  *string `json:"genre"`
}

// UpdateBookInput especifica campos opcionales para actualización.
type UpdateBookInput struct {
	Title  *string
	Author *string
	Year   *uint
	Genre  *string
	ISBN   *string
}

func NewBook(title, author string, year uint, genre, isbn string) *Book {
	// agregar el tiempo en UTC
	now := time.Now().UTC()
	return &Book{
		Title:     strings.TrimSpace(title),
		Author:    strings.TrimSpace(author),
		Year:      year,
		Genre:     strings.TrimSpace(genre),
		ISBN:      NormalizeISBN(isbn),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *Book) Update(input UpdateBookInput) {
	if input.Title != nil {
		b.Title = strings.TrimSpace(*input.Title)
	}
	if input.Author != nil {
		b.Author = strings.TrimSpace(*input.Author)
	}
	if input.Year != nil {
		b.Year = *input.Year
	}
	if input.Genre != nil {
		b.Genre = strings.TrimSpace(*input.Genre)
	}
	if input.ISBN != nil {
		b.ISBN = NormalizeISBN(*input.ISBN)
	}

	// actualizar el tiempo en UTC
	b.UpdatedAt = time.Now().UTC()
}

// HasAny indica si el filtro tiene algún criterio.
func (f BookFilter) HasAny() bool {
	return f.Title != nil || f.Author != nil || f.Year != nil || f.Genre != nil
}

// ValidateBasic valida reglas de negocio esenciales del libro
func (b *Book) ValidateBasic() error {
	if strings.TrimSpace(b.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(b.Author) == "" {
		return errors.New("author is required")
	}
	if err := ValidateYear(b.Year); err != nil {
		return err
	}

	if err := ValidateISBN(b.ISBN); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateInput(in UpdateBookInput) error {
	if in.Title != nil && strings.TrimSpace(*in.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if in.Author != nil && strings.TrimSpace(*in.Author) == "" {
		return errors.New("author cannot be empty")
	}
	if in.Year != nil {
		if err := ValidateYear(*in.Year); err != nil {
			return err
		}
	}
	if in.ISBN != nil {
		if err := ValidateISBN(*in.ISBN); err != nil {
			return err
		}
	}
	return nil
}

// NormalizeISBN quita guiones/espacios y pone en mayúsculas.
func NormalizeISBN(isbn string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(isbn), "-", ""), " ", ""))
}

// ValidateISBN soporta ISBN-10 e ISBN-13 con verificación de dígito.
func ValidateISBN(isbn string) error {
	n := NormalizeISBN(isbn)

	// ISBN-10: 9 dígitos + (dígito o X)
	re10 := regexp.MustCompile(`^\d{9}(\d|X)$`)
	// ISBN-13: 13 dígitos
	re13 := regexp.MustCompile(`^\d{13}$`)

	switch {
	case re10.MatchString(n):
		sum := 0
		for i := 0; i < 9; i++ {
			sum += (10 - i) * int(n[i]-'0')
		}
		check := 0
		if n[9] == 'X' {
			check = 10
		} else {
			check = int(n[9] - '0')
		}
		sum += check
		if sum%11 != 0 {
			return errors.New("invalid ISBN-10 checksum")
		}
		return nil
	case re13.MatchString(n):
		sum := 0
		for i := 0; i < 12; i++ {
			d := int(n[i] - '0')
			if i%2 == 0 {
				sum += d
			} else {
				sum += 3 * d
			}
		}
		expected := (10 - (sum % 10)) % 10
		actual := int(n[12] - '0')
		if actual != expected {
			return errors.New("invalid ISBN-13 checksum")
		}
		return nil
	default:
		return errors.New("isbn must be ISBN-10 or ISBN-13 format")
	}
}

// ValidateYear asegura un año dentro de un rango razonable.
func ValidateYear(year uint) error {
	const minYear = 1450
	current := time.Now().UTC().Year()
	if int(year) < minYear || int(year) > current {
		return errors.New("year must be between 1450 and current year")
	}
	return nil
}
