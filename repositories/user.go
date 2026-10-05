package repositories

import "naughtfound.github.io/go2phd/models"

type UserRepository struct{}

func NewUserRepo() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByID(id int) (models.User, error) {
	return models.User{ID: id, Name: "Alice"}, nil // Mock DB
}
