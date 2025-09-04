package repository

import (
	"context"
	"userService/config"
	"userService/internal/dto"
	"userService/internal/model"
)

type InterfaceUserRepository interface {
	Create(ctx context.Context, user dto.UserRegister) (int, error)
	Delete(ctx context.Context, id int) error
	Find(username string) (model.User, error)
	Exists(username string) (bool, error)
}

type UserRepository struct {
	db *config.Database
}

func NewUserRepository(db *config.Database) InterfaceUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, user dto.UserRegister) (int, error) {
	var id int
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	err = r.db.DB.QueryRow(`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		user.Username, user.Password).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *UserRepository) Find(username string) (model.User, error) {
	var user model.User
	err := r.db.DB.QueryRow(`SELECT id, username, password FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Exists(username string) (bool, error) {
	var exists bool
	err := r.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
