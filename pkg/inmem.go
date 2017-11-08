package orchlift

import (
  "sync"

)

type taskRepository struct {
  mtx sync.RWMutex
  tasks map[string]*Task
}

func (r *taskRepository) Store(t *Task) error {
  r.mtx.Lock()
  defer r.mtx.Unlock()
  r.tasks[t.ID] = t
  return nil
}

func (r *taskRepository) Find(id string) (*Task, error) {
  r.mtx.Lock()
  defer r.mtx.Unlock()
  if val, ok := r.tasks[id]; ok {
    return val, nil
  }
  return nil, ErrUnknown
}

func (r *taskRepository) FindAll() ([]*Task) {
  r.mtx.Lock()
  defer r.mtx.Unlock()
  t := make([]*Task, 0, len(r.tasks))
  for _, val := range r.tasks {
    t = append(t, val)
  }
  return t
}

func NewTaskRepository() Repository {
  return &taskRepository{
    tasks: make(map[string]*Task),
  }
}
