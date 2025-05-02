package repository

import (
	"database/sql"
	"test_service/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestEmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").WithArgs(email).WillReturnRows(rows)

	exists, err := repo.EmailExists(email)
	assert.NoError(t, err)
	assert.True(t, exists)

	rows = sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").WithArgs(email).WillReturnRows(rows)

	exists, err = repo.EmailExists(email)
	assert.NoError(t, err)
	assert.False(t, exists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userId := uuid.New()
	now := time.Now()
	user := &model.User{
		ID:        userId,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		Created:   now,
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("john@example.com").WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(userId, "John", "Doe", "john@example.com", uint(30), now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userId := uuid.New()
	user := &model.User{
		ID:        userId,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@example.com",
		Age:       30,
		Created:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("john@example.com").WillReturnRows(rows)

	err = repo.CreateUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userId := uuid.New()
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "age", "created"}).
		AddRow(userId, "John", "Doe", "john@example.com", uint(30), now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE").
		WithArgs(userId.String()).
		WillReturnRows(rows)

	user, err := repo.GetUser(userId.String())
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userId, user.ID)
	assert.Equal(t, "John", user.Firstname)
	assert.Equal(t, "Doe", user.Lastname)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, uint(30), user.Age)
	assert.Equal(t, now, user.Created)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userId := uuid.New().String()
	firstname := "Jane"
	email := "jane@example.com"
	updateInput := &model.UpdateUserInput{
		Firstname: &firstname,
		Email:     &email,
	}

	userId_uuid, _ := uuid.Parse(userId)
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "age", "created"}).
		AddRow(userId_uuid, "John", "Doe", "john@example.com", uint(30), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").
		WithArgs(userId).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("jane@example.com").WillReturnRows(rows)

	mock.ExpectExec("UPDATE users SET").
		WithArgs(firstname, nil, email, nil, userId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUser(userId, updateInput)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser_EmailAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userId := uuid.New().String()
	email := "jane@example.com"
	updateInput := &model.UpdateUserInput{
		Email: &email,
	}

	userId_uuid, _ := uuid.Parse(userId)
	rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "age", "created"}).
		AddRow(userId_uuid, "John", "Doe", "john@example.com", uint(30), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").
		WithArgs(userId).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("jane@example.com").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"email"}).AddRow("john@example.com")
	mock.ExpectQuery("SELECT email FROM users WHERE").WithArgs(userId).WillReturnRows(rows)

	err = repo.UpdateUser(userId, updateInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already in use")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	userId := uuid.New().String()
	mock.ExpectQuery("SELECT (.+) FROM users WHERE").
		WithArgs(userId).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUser(userId)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
