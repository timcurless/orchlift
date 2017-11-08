package orchlift

import (
  "errors"
  "time"

  "github.com/pborman/uuid"
)

type Task struct {
  ID              string
  Name            string
  TaskType        string
  StartTime       time.Time
  EndTime         time.Time
  ExecutionStatus string
}

func NewTask(id, name, taskType string, startTime, endTime time.Time, executionStatus string) *Task {
  return &Task{
    ID:              id,
    Name:            name,
    TaskType:        taskType,
    StartTime:       startTime,
    EndTime:         endTime,
    ExecutionStatus: executionStatus,
  }
}

type Repository interface {
  Store(task *Task) error
  Find(id string) (*Task, error)
  FindAll() []*Task
}

func NextID() string {
  return uuid.New()
}

var ErrUnknown = errors.New("Unknown Task")
