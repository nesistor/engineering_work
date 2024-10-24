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

type AdminModel struct {
	DB *sql.DB
}

type Models struct {
	Admin AdminModel
}

func New(db *sql.DB) Models {
	return Models{
		Admin: AdminModel{DB: db},
	}
}

type Admin struct {
	ID           int64     `json:"id"`
	AdminName    string    `json:"admin_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordhash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetAdminByEmail retrieves an admin by their email address
func (a *AdminModel) GetAdminByEmail(email string) (*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, admin_name, passwordhash, created_at, updated_at FROM admins WHERE email = $1`

	var admin Admin
	row := a.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&admin.ID,
		&admin.Email,
		&admin.AdminName,
		&admin.PasswordHash,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (a *AdminModel) GetAdminByID(id int64) (*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, admin_name, passwordhash, created_at, updated_at 
	          FROM admins 
	          WHERE id = $1`

	var admin Admin
	row := a.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&admin.ID,
		&admin.Email,
		&admin.AdminName,
		&admin.PasswordHash,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (a *AdminModel) InsertAdmin(admin Admin) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	query := `INSERT INTO admins (email, admin_name, passwordhash, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = a.DB.QueryRowContext(ctx, query,
		admin.Email,
		admin.AdminName,
		hashedPassword,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (a *AdminModel) DeleteAdminByID(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM admins WHERE id = $1`

	_, err := a.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *AdminModel) UpdateAdmin(admin Admin) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE admins SET
		email = $1,
		admin_name = $2,
		updated_at = $3
		WHERE id = $4
	`

	_, err := a.DB.ExecContext(ctx, query,
		admin.Email,
		admin.AdminName,
		time.Now(),
		admin.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (a *AdminModel) ResetAdminPassword(adminID int64, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `UPDATE admins SET passwordhash = $1, updated_at = $2 WHERE id = $3`
	_, err = a.DB.ExecContext(ctx, stmt, hashedPassword, time.Now(), adminID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare an admin supplied password
func (a *AdminModel) PasswordMatches(adminID int64, plainText string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var hashedPassword string
	query := `SELECT passwordhash FROM admins WHERE id = $1`
	err := a.DB.QueryRowContext(ctx, query, adminID).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

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
