package orchlift

import (
  "context"
  "time"

  "github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
  AddTaskEndpoint   endpoint.Endpoint
  ListTasksEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
  return Endpoints{
    AddTaskEndpoint:   makeAddTaskEndpoint(s),
    ListTasksEndpoint: makeListTasksEndpoint(s),
  }
}

// Add a New Task Endpoint
type addTaskRequest struct {
  Name            string
  TaskType        string
  StartTime       time.Time
  EndTime         time.Time
  ExecutionStatus string
}

type addTaskResponse struct {
  ID  string    `json:"id,omitempty"`
  Err error     `json:"error,omitempty"`
}

func (r addTaskResponse) error() error { return r.Err }

func makeAddTaskEndpoint(s Service) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(addTaskRequest)
    id, err := s.CreateNewTask(req.Name, req.TaskType, req.StartTime, req.EndTime, req.ExecutionStatus)
    return addTaskResponse{ID: id, Err: err}, nil
  }
}
// ****************************************************

// List all Tasks Endpoint
type listTasksRequest struct {}

type listTasksResponse struct {
  Tasks []Task `json:"tasks,omitempty"`
  Err   error       `json:"error,omitempty"`
}

func (r listTasksResponse) error() error { return r.Err }

func makeListTasksEndpoint(s Service) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    _ = request.(listTasksRequest)
    return listTasksResponse{Tasks: s.Tasks(), Err: nil}, nil
  }
}
// ****************************************************
