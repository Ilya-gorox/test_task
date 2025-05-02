package repository

import (
	"database/sql"
	"fmt"
	"os"
	"test_service/internal/model"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewPostgresDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *model.User) error {
	if exists, err := r.EmailExists(user.Email); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	query := `INSERT INTO users (id, firstname, lastname, email, age, created) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, user.ID, user.Firstname, user.Lastname, user.Email, user.Age, user.Created)
	return err
}

func (r *Repository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *Repository) GetUser(id string) (*model.User, error) {
	query := `SELECT id, firstname, lastname, email, age, created FROM users WHERE id = $1`
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Age, &user.Created)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(id string, input *model.UpdateUserInput) error {
	if _, err := r.GetUser(id); err != nil {
		return err
	}

	if input.Email != nil {
		exists, err := r.EmailExists(*input.Email)
		if err != nil {
			return err
		}
		if exists {
			var currentEmail string
			err := r.db.QueryRow("SELECT email FROM users WHERE id = $1", id).Scan(&currentEmail)
			if err != nil {
				return err
			}
			if currentEmail != *input.Email {
				return fmt.Errorf("email %s already in use", *input.Email)
			}
		}
	}

	query := `UPDATE users SET 
			  firstname = COALESCE($1, firstname),
			  lastname = COALESCE($2, lastname),
			  email = COALESCE($3, email),
			  age = COALESCE($4, age)
			  WHERE id = $5`

	_, err := r.db.Exec(query, input.Firstname, input.Lastname, input.Email, input.Age, id)
	return err
}
