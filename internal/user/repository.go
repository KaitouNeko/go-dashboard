package user

import (
	"ai-workshop/internal/models"
	"ai-workshop/internal/utils/errorutils"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user models.User) error {
	query := `INSERT INTO users (name, email, password, permission) VALUES (:name, :email, :password, :permission)`

	_, err := r.DB.NamedExec(query, user)

	if err != nil {
		fmt.Println("Error when creating user:", err)
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *UserRepository) UpdatePassword(params UserUpdatePasswordParams) error {
	query := `UPDATE users SET password = :password WHERE id = :id`

	result, err := r.DB.NamedExec(query, params)
	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no member found with id: %v", params.ID)
	}

	return nil
}

func (r *UserRepository) UpdateInfo(params UserUpdateInfoParams, userId uuid.UUID) error {
	query := `UPDATE users SET name = :name, permission = :permission WHERE id = :id`

	result, err := r.DB.NamedExec(query, params)

	fmt.Println("result", result)
	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %v", params.ID)
	}

	return nil
}

func (r *UserRepository) GetByIdWithPassword(id uuid.UUID) (*models.User, error) {
	query := `SELECT * FROM users WHERE users.id = $1`

	var user models.User

	err := r.DB.Get(&user, query, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUsers() (*[]models.User, error) {
	query := `SELECT id, email, name, updated_at, permission, status FROM users`

	var users []models.User

	err := r.DB.Select(&users, query)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *UserRepository) GetById(id uuid.UUID) (*models.User, error) {
	query := `SELECT * FROM users WHERE users.id = $1`

	var user models.User

	err := r.DB.Get(&user, query, id)

	if err != nil {
		return nil, err
	}

	// Remove password from the struct
	user.Password = ""

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE users.email = $1`

	err := r.DB.Get(&user, query, email)
	fmt.Println("Error:", err)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateDefaultUsers(users []CreateDefaultUser) error {
	query := `
	INSERT INTO users(id, email, name, password, status)
	VALUES(:id, :email, :name, :password, :status)
	ON CONFLICT (id) DO NOTHING
	`
	_, err := r.DB.NamedExec(query, users)

	fmt.Printf("DEBUG: Error when creating default user: %s\n", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
