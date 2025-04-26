package grpchandler

import (
	"context"

	"github.com/literaen/simple_project/users/internal/users"

	taskpb "github.com/literaen/simple_project/api"
	userpb "github.com/literaen/simple_project/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	service *users.UserService
}

func NewUserHandler(service *users.UserService) *UserHandler {
	handler := &UserHandler{service: service}

	return handler
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*userpb.User, len(users))
	for i := 0; i < len(users); i++ {
		resp[i] = &userpb.User{
			Id:    users[i].ID,
			Name:  users[i].Name,
			Email: users[i].Email,
		}
	}

	return &userpb.GetAllUsersResponse{
		Users: resp,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	id := req.GetId()
	user, err := h.service.GetUserByID(id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoTasks []*taskpb.Task
	for _, t := range user.Tasks {
		protoTasks = append(protoTasks, &taskpb.Task{
			Id:          t.ID,
			UserId:      t.UserID,
			Description: t.Description,
		})
	}

	// Формируем gRPC-ответ
	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:    id,
			Name:  user.User.Name,
			Email: user.User.Email,
		},
		Tasks: protoTasks,
	}, nil
}

func (h *UserHandler) AddUser(ctx context.Context, req *userpb.AddUserRequest) (*userpb.AddUserResponse, error) {
	data := req.GetUser()

	user := &users.User{
		Name:  data.GetName(),
		Email: data.GetEmail(),
	}

	if err := h.service.PostUser(user); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.AddUserResponse{
		Id: user.ID,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	id := req.GetId()
	err := h.service.DeleteUserByID(id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.DeleteUserResponse{Success: true}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	id := req.GetId()
	data := req.GetUser()

	user, err := h.service.PatchUserByID(id, &users.User{
		Name:  data.GetName(),
		Email: data.GetEmail(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.UpdateUserResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
