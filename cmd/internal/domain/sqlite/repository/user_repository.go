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

func (u *DefaultUserRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	err := u.db.Find(&users).Error
	return users, err
}

func (u *DefaultUserRepository) FindBySub(sub string) (*entity.User, error) {
	var user entity.User
	err := u.db.Where("sub_uuid = ?", sub).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (u *DefaultUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (u *DefaultUserRepository) ExistsByEmail(email string) (bool, error) {
	var exists int
	err := u.db.
		Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).
		Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (u *DefaultUserRepository) Save(user *entity.User) error {
	return u.db.Save(user).Error
}
