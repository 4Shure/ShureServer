package repository

import (
	"4shure/cmd/internal/domain/entity"
	"errors"
	"gorm.io/gorm"
)

type DefaultUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{db: db}
}

func (u *DefaultUserRepository) FindByID(id int) (*entity.User, error) {
	var user entity.User
	err := u.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
