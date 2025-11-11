package infrastructure

import (
	"api-go-gestion-libros-hexagonal/modules/book/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type SqlBookRepository struct {
	db *sql.DB
}

func NewSqlBookRepository(db *sql.DB) *SqlBookRepository {
	return &SqlBookRepository{db: db}
}

// Create crea un nuevo libro en el repositorio
func (r *SqlBookRepository) Create(ctx context.Context, book *domain.Book) error {
	isbn := domain.NormalizeISBN(book.ISBN)
	q := `INSERT INTO books (title, author, year, genre, isbn, created_at, updated_at)
	      VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q,
		strings.TrimSpace(book.Title),
		strings.TrimSpace(book.Author),
		book.Year,
		strings.TrimSpace(book.Genre),
		isbn,
		book.CreatedAt,
		book.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("duplicate isbn: %s", isbn)
		}
		return err
	}

	// Recuperar ID asignado por ISBN
	row := r.db.QueryRowContext(ctx, "SELECT id FROM books WHERE isbn = ?", isbn)
	var id int64
	if err := row.Scan(&id); err != nil {
		return err
	}
	book.ID = uint(id)
	return nil
}

// Update actualiza un libro existente en el repositorio
func (r *SqlBookRepository) Update(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	isbn := domain.NormalizeISBN(book.ISBN)
	q := `UPDATE books
	      SET title = ?, author = ?, year = ?, genre = ?, isbn = ?, updated_at = ?
	      WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q,
		strings.TrimSpace(book.Title),
		strings.TrimSpace(book.Author),
		int(book.Year),
		strings.TrimSpace(book.Genre),
		isbn,
		time.Now().UTC(),
		int(book.ID),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("duplicate isbn: %s", isbn)
		}
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, fmt.Errorf("book %d not found", book.ID)
	}
	return r.GetByID(ctx, book.ID)
}

// Delete elimina un libro existente en el repositorio
func (r *SqlBookRepository) Delete(ctx context.Context, id uint) error {
	q := `DELETE FROM books WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, int(id))
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("book %d not found", id)
	}
	return nil
}

// GetAll obtiene todos los libros del repositorio
func (r *SqlBookRepository) GetAll(ctx context.Context) ([]*domain.Book, error) {
	q := `SELECT id, title, author, year, genre, isbn, created_at, updated_at FROM books ORDER BY id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBooks(rows)
}

// GetByISBN obtiene un libro por ISBN del repositorio
func (r *SqlBookRepository) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	n := domain.NormalizeISBN(isbn)
	q := `SELECT id, title, author, year, genre, isbn, created_at, updated_at FROM books WHERE isbn = ?`
	row := r.db.QueryRowContext(ctx, q, n)
	return scanBook(row)
}

// GetByID obtiene un libro por ID
func (r *SqlBookRepository) GetByID(ctx context.Context, id uint) (*domain.Book, error) {
	q := `SELECT id, title, author, year, genre, isbn, created_at, updated_at FROM books WHERE id = ?`
	row := r.db.QueryRowContext(ctx, q, int(id))
	return scanBook(row)
}

// FindByFilter obtiene libros por filtros del repositorio
func (r *SqlBookRepository) FindByFilter(ctx context.Context, filter domain.BookFilter) ([]*domain.Book, error) {
	clauses := []string{}
	args := []any{}

	addLike := func(col string, p *string) {
		if p != nil {
			needle := "%" + strings.ToLower(strings.TrimSpace(*p)) + "%"
			clauses = append(clauses, fmt.Sprintf("LOWER(%s) LIKE ?", col))
			args = append(args, needle)
		}
	}
	addEqUint := func(col string, p *uint) {
		if p != nil {
			clauses = append(clauses, fmt.Sprintf("%s = ?", col))
			args = append(args, int(*p))
		}
	}

	addLike("title", filter.Title)
	addLike("author", filter.Author)
	addEqUint("year", filter.Year)
	addLike("genre", filter.Genre)

	q := `SELECT id, title, author, year, genre, isbn, created_at, updated_at FROM books`
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}
	q += " ORDER BY id"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBooks(rows)
}

// Helpers

func scanBook(row *sql.Row) (*domain.Book, error) {
	var (
		id        int64
		title     string
		author    string
		year      int64
		genre     string
		isbn      string
		createdAt time.Time
		updatedAt time.Time
	)
	if err := row.Scan(&id, &title, &author, &year, &genre, &isbn, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return &domain.Book{
		ID:        uint(id),
		Title:     title,
		Author:    author,
		Year:      uint(year),
		Genre:     genre,
		ISBN:      isbn,
		CreatedAt: createdAt.UTC(),
		UpdatedAt: updatedAt.UTC(),
	}, nil
}

func scanBooks(rows *sql.Rows) ([]*domain.Book, error) {
	out := []*domain.Book{}
	for rows.Next() {
		var (
			id        int64
			title     string
			author    string
			year      int64
			genre     string
			isbn      string
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&id, &title, &author, &year, &genre, &isbn, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, &domain.Book{
			ID:        uint(id),
			Title:     title,
			Author:    author,
			Year:      uint(year),
			Genre:     genre,
			ISBN:      isbn,
			CreatedAt: createdAt.UTC(),
			UpdatedAt: updatedAt.UTC(),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func isUniqueViolation(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint")
}
