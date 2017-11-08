package orchlift

import (
  "errors"
  "time"

)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
  CreateNewTask(name, taskType string, startTime, endTime time.Time, executionStatus string) (string, error)
  Tasks() []Task
}

type service struct {
  tasks Repository
}

func (s *service) CreateNewTask(name, taskType string, startTime, endTime time.Time, executionStatus string) (string, error) {
  if name == "" || taskType == "" {
    return "", ErrInvalidArgument
  }

  id := NextID()

  var status string
  if executionStatus == "" {
    status = "NOT_STARTED"
  } else {
    status = executionStatus
  }

  var thisStartTime time.Time
  if startTime.IsZero() {
    thisStartTime = time.Now()
  } else {
    thisStartTime = startTime
  }

  t := NewTask(id, name, taskType, thisStartTime, endTime, status)

  if err := s.tasks.Store(t); err != nil {
    return "", err
  }
  return t.ID, nil
}

func (s *service) Tasks() []Task {
  var result []Task
  for _, t := range s.tasks.FindAll() {
    result = append(result, assemble(t))
  }
  return result
}

func NewService(tasks Repository) Service {
  return &service{
    tasks: tasks,
  }
}

func assemble(t *Task) Task {
  return Task{
    ID:              t.ID,
    Name:            t.Name,
    TaskType:        t.TaskType,
    StartTime:       t.StartTime,
    EndTime:         t.EndTime,
    ExecutionStatus: t.ExecutionStatus,
  }
}
