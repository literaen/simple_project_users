package users

import (
	"errors"
	"fmt"

	"github.com/literaen/simple_project/pkg/postgres"
	"github.com/literaen/simple_project/pkg/redis"

	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(fn func(tx *gorm.DB) error) error

	// Получить всех пользователей
	GetAllUsers() ([]User, error)

	// Получить пользователя по ID
	GetUserByID(id uint64) (*User, error)

	// Создать нового пользователя
	PostUser(user *User) error

	// Изменить пользоваля по ID
	PatchUserByID(id uint64, user *User) (*User, error)

	// Удалить пользоваля по ID
	DeleteUserByID(tx *gorm.DB, id uint64) error
}

type userRepository struct {
	gdb   *postgres.GDB
	redis *redis.RDB
}

func NewUserRepository(gdb *postgres.GDB, redis *redis.RDB) UserRepository {
	return &userRepository{gdb: gdb, redis: redis}
}

func (r *userRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.gdb.DB.Transaction(fn)
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.gdb.DB.Find(&users).Error
	return users, err
}

func (r *userRepository) GetUserByID(id uint64) (*User, error) {
	var user User
	if err := r.gdb.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) PostUser(user *User) error {
	if err := r.gdb.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) PatchUserByID(id uint64, user *User) (*User, error) {
	var resp *User
	res := r.gdb.DB.
		Model(&User{}).
		Where("id = ?", id).
		Updates(&user).
		Scan(&resp)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}

	return resp, nil
}

func (r *userRepository) DeleteUserByID(tx *gorm.DB, id uint64) error {
	var user User
	res := tx.First(&user, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	if err := tx.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
