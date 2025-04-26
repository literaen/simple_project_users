package grpcclients

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/literaen/simple_project/users/internal/config"

	grpcclient "github.com/literaen/simple_project/pkg/grpc/client"

	taskpb "github.com/literaen/simple_project/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type TaskClientConstructor struct{}

func (c *TaskClientConstructor) NewClient(conn *grpc.ClientConn) interface{} {
	return taskpb.NewTaskServiceClient(conn)
}

type TaskGRPCClient struct {
	client *grpcclient.Client
}

func NewTaskGRPCClient(cfg *config.Config) *TaskGRPCClient {
	client := &TaskGRPCClient{
		client: grpcclient.NewClient(5*time.Second, &TaskClientConstructor{}),
	}

	client.Start(cfg)

	return client
}

func (s *TaskGRPCClient) GetUserAllTasks(ctx context.Context, id uint64) (*taskpb.GetAllTasksResponse, error) {
	if !s.client.IsReady() {
		return nil, fmt.Errorf("user service unavailable")
	}

	tasksResp, err := s.GetUserClient().GetUserAllTasks(ctx, &taskpb.GetTaskRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, errors.New(st.Message())
		} else {
			return nil, fmt.Errorf("unknown error: %v", err)
		}
	}

	return tasksResp, nil
}

// func (s *TaskGRPCClient) GetUser(ctx context.Context, id uint64) error {
// 	if !s.client.IsReady() {
// 		return fmt.Errorf("user service unavailable")
// 	}

// 	_, err := s.GetUserClient().GetUser(ctx, &taskpb.GetUserRequest{Id: id})
// 	if err != nil {
// 		st, ok := status.FromError(err)
// 		if ok {
// 			return errors.New(st.Message())
// 		} else {
// 			return fmt.Errorf("unknown error: %v", err)
// 		}
// 	}

// 	return nil
// }

// GetUserClient возвращает типизированный клиент
func (s *TaskGRPCClient) GetUserClient() taskpb.TaskServiceClient {
	return s.client.GetClient().(taskpb.TaskServiceClient)
}

func (s *TaskGRPCClient) Start(cfg *config.Config) {
	target := fmt.Sprintf("%s:%s", cfg.TASK_SERVICE_HOST, cfg.TASK_SERVICE_PORT)
	s.client.AutoReconnect(context.TODO(), target)
}

func (s *TaskGRPCClient) Close() error {
	return s.client.Close()
}
