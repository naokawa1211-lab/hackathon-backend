// internal/repository/user_repository.go
package repository

import (
	"database/sql"
	"hackathon-backend/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// UpsertUser はユーザー情報を保存または更新します
func (r *UserRepository) UpsertUser(user *model.User) error {
	query := `
		INSERT INTO users (id, username, avatar_url, space_credits, created_at)
		VALUES (?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE username = ?, avatar_url = ?`
	
	_, err := r.db.Exec(query, 
		user.ID, user.Username, user.AvatarURL, user.SpaceCredits, 
		user.Username, user.AvatarURL,
	)
	return err
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	query := `SELECT id, username, avatar_url, space_credits FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.AvatarURL, &user.SpaceCredits)
	if err == sql.ErrNoRows {
		return nil, nil // 見つからない場合はnil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}