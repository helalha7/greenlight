package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/helalha7/greenlight.git/internal/validator"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`
	Runtime   int       `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int       `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= time.Now().Year(), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES (?, ?, ?, ?)
	`

	genres, err := json.Marshal(movie.Genres)
	if err != nil {
		return err
	}

	args := []any{movie.Title, movie.Year, movie.Runtime, genres}

	res, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	movie.ID = int(id)
	return nil
}

func (m MovieModel) Get(id int) (*Movie, error) {
	query := `
		SELECT id, create_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = ?
	`
	var genres []byte
	movie := &Movie{}

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&genres,
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if err := json.Unmarshal(genres, &movie.Genres); err != nil {
		return nil, err
	}

	return movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = ?, year = ?, runtime = ?, genres = ?, version = version + 1
		WHERE id = ? AND version = ?
	`
	genres, err := json.Marshal(movie.Genres)
	if err != nil {
		return err
	}

	args := []any{movie.Title, movie.Year, movie.Runtime, genres, movie.ID, movie.Version}

	res, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	movie.Version++
	return nil
}

func (m MovieModel) Delete(id int) error {
	query := `
		DELETE FROM movies
		WHERE id = ?
	`

	res, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) GetAll() ([]*Movie, error) {
	query := `
		SELECT id, create_at, title, year, runtime, genres, version
		FROM movies
		ORDER BY id
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	movies := []*Movie{}

	for rows.Next() {
		var movie Movie
		var genres []byte

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			&genres,
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(genres, &movie.Genres); err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
