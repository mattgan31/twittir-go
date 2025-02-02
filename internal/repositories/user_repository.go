package repositories

import (
	"errors"

	"twittir-go/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	SaveNewUser(user *domain.User) error
	FindByID(id int) (*domain.User, error)
	SearchByUsername(username string) ([]domain.User, error)
	UpdateProfile(userID int, updateUser *domain.User) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Debug().Where("username = ?", username).Take(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SaveNewUser(user *domain.User) error {
	err := r.db.Debug().Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SearchByUsername(username string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Debug().Where("username LIKE ?", "%"+username+"%").Find(&users).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateProfile(userID int, updateUser *domain.User) (*domain.User, error) {

	user := domain.User{}

	// Pastikan user dengan ID ini ada
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := r.db.Debug().Model(&user).Updates(map[string]interface{}{"FullName": updateUser.FullName, "Username": updateUser.Username, "Bio": updateUser.Bio}).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
