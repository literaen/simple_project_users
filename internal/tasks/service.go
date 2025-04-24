package tasks

import (
	"context"
	"fmt"

	grpcclient "github.com/literaen/simple_project/users/internal/grpc/client"

	taskpb "github.com/literaen/simple_project/proto/gen"
)

type TaskService struct {
	grpc *grpcclient.TaskGRPCClient
}

func NewTaskService(grpc *grpcclient.TaskGRPCClient) *TaskService {
	svc := &TaskService{
		grpc: grpc,
	}

	return svc
}

func (s *TaskService) GetUserAllTasks(id uint64) (*taskpb.GetAllTasksResponse, error) {
	tasksResp, err := s.grpc.GetUserAllTasks(context.TODO(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users tasks: %v", err)
	}

	return tasksResp, nil
}
