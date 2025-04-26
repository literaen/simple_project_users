package users

import (
	"context"
	"fmt"

	"github.com/literaen/simple_project/users/internal/outbox"

	"github.com/literaen/simple_project/dto"
	grpcclient "github.com/literaen/simple_project/users/internal/grpc/client"
	"gorm.io/gorm"
)

type UserService struct {
	repo           UserRepository
	outboxService  *outbox.OutBoxService
	taskGRPCClient *grpcclient.TaskGRPCClient
}

func NewUserService(repo UserRepository, outboxService *outbox.OutBoxService, taskGRPCClient *grpcclient.TaskGRPCClient) *UserService {
	return &UserService{repo: repo, outboxService: outboxService, taskGRPCClient: taskGRPCClient}
}

func (s *UserService) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id uint64) (*dto.UserWithTasks, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	tasksResp, err := s.taskGRPCClient.GetUserAllTasks(context.TODO(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users tasks: %v", err)
	}

	var dtoTasks []dto.Task
	for _, t := range tasksResp.Tasks {
		dtoTasks = append(dtoTasks, dto.Task{
			ID:          t.Id,
			UserID:      t.UserId,
			Description: t.Description,
		})
	}

	dtoUser := &dto.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	dto := &dto.UserWithTasks{
		User:  dtoUser,
		Tasks: dtoTasks,
	}

	return dto, nil
}

func (s *UserService) PostUser(user *User) error {
	return s.repo.PostUser(user)
}

func (s *UserService) PatchUserByID(id uint64, user *User) (*User, error) {
	return s.repo.PatchUserByID(id, user)
}

// func (s *UserService) DeleteUserByID(id uint64) error {
// 	return s.repo.DeleteUserByID(id)
// }

func (s *UserService) DeleteUserByID(id uint64) error {
	return s.repo.WithTx(func(tx *gorm.DB) error {
		if err := s.repo.DeleteUserByID(tx, id); err != nil {
			return err
		}

		return s.outboxService.AddEvent(tx, "user.deleted", map[string]interface{}{
			"user_id": id,
		})
	})
}
