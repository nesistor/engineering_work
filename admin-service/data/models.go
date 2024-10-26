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

type NewAdminModel struct {
	DB *sql.DB
}

type Models struct {
	Admin    AdminModel
	NewAdmin NewAdminModel
	Token    TokenModel
}

func New(db *sql.DB) Models {
	return Models{
		Admin:    AdminModel{DB: db},
		NewAdmin: NewAdminModel{DB: db},
		Token:    TokenModel{DB: db},
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

// NewAdmin struct for storing email addresses of admins who can create accounts
type NewAdmin struct {
	Email string `json:"email"`
}

// GetAllAdmins returns a slice of all admins, sorted by their admin name
func (a *AdminModel) GetAllAdmins() ([]*Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, admin_name, passwordhash, created_at, updated_at 
			  FROM admins ORDER BY admin_name`

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*Admin

	for rows.Next() {
		var admin Admin
		err := rows.Scan(
			&admin.ID,
			&admin.Email,
			&admin.AdminName,
			&admin.PasswordHash,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		admins = append(admins, &admin)
	}

	return admins, nil
}

// InsertNewAdmin inserts a new admin email into the database
func (n *NewAdminModel) InsertNewAdmin(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO new_admins (email, created_at) VALUES ($1, $2)`

	_, err := n.DB.ExecContext(ctx, query, email, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetAllNewAdmins retrieves all new admin emails from the database
func (n *NewAdminModel) GetAllNewAdmins() ([]NewAdmin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT email FROM new_admins`
	rows, err := n.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newAdmins []NewAdmin
	for rows.Next() {
		var newAdmin NewAdmin
		if err := rows.Scan(&newAdmin.Email); err != nil {
			return nil, err
		}
		newAdmins = append(newAdmins, newAdmin)
	}

	return newAdmins, nil
}

// DeleteNewAdmin removes an admin email from the new admins table
func (n *NewAdminModel) DeleteNewAdmin(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM new_admins WHERE email = $1`

	_, err := n.DB.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
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

// UpdateAdminPassword updates the admin's password in the database.
func (a *AdminModel) UpdateAdminPassword(email string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `UPDATE admins SET passwordhash = $1, updated_at = $2 WHERE email = $3`
	_, err = a.DB.ExecContext(ctx, stmt, hashedPassword, time.Now(), email)
	if err != nil {
		return err
	}

	return nil
}
