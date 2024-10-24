package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

type UserModel struct {
	DB *sql.DB
}

type Models struct {
	User UserModel
}

func New(db *sql.DB) Models {
	return Models{
		User: UserModel{DB: db},
	}
}

type User struct {
	ID           int64     `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordhash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"udpated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *UserModel) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, passwordhash,
			  created_at, updated_at FROM users ORDER BY last_name`

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.UserName,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *UserModel) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, username, passwordhash, created_at, updated_at FROM users WHERE email = $1`

	var user User
	row := u.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) GetUserByID(id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, username, passwordhash, created_at, updated_at 
	          FROM users 
	          WHERE id = $1`

	var user User
	row := u.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) InsertUser(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	query := `INSERT INTO users (email, username, passwordhash, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = u.DB.QueryRowContext(ctx, query,
		user.Email,
		user.UserName,
		hashedPassword,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (u *UserModel) DeleteUserByID(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserModel) UpdateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE users SET
		email = $1,
		username = $2,
				updated_at = $4
		WHERE id = $5
	`

	_, err := u.DB.ExecContext(ctx, query,
		user.Email,
		user.UserName,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserModel) ResetUserPassword(userID int64, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `UPDATE users SET passwordhash = $1, updated_at = $2 WHERE id = $3`
	_, err = u.DB.ExecContext(ctx, stmt, hashedPassword, time.Now(), userID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
func (u *UserModel) PasswordMatches(userID int64, plainText string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var hashedPassword string
	query := `SELECT passwordhash FROM users WHERE id = $1`
	err := u.DB.QueryRowContext(ctx, query, userID).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	log.Println("MODELS: Stored hashed password:", hashedPassword) // Debug Log
	log.Println("MODELS: Request password:", plainText)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Println("MODELS: Error Mismatched password:", err)
			// niepoprawne hasło
			return false, nil
		}
		log.Println("MODELS: Other error in CompareHashAndPassword:", err)
		return false, err

	}

	return true, nil
}
