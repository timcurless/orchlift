package orchlift

import (
  "bytes"
  "context"
  "encoding/json"
  "errors"
  "io/ioutil"
  "net/http"
  "net/url"
  //"strings"
  "time"

  //"golang.org/x/time/rate"
  "github.com/gorilla/mux"
  //"github.com/sony/gobreaker"
  //"github.com/go-kit/kit/endpoint"
  //"github.com/go-kit/kit/circuitbreaker"
  //"github.com/go-kit/kit/ratelimit"
  kitlog "github.com/go-kit/kit/log"
  kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
  opts := []kithttp.ServerOption{
          kithttp.ServerErrorLogger(logger),
          kithttp.ServerErrorEncoder(encodeError),
  }

  addTaskHandler := kithttp.NewServer(
    makeAddTaskEndpoint(s),
    decodeAddTaskRequest,
    encodeResponse,
    opts...,
  )
  listTasksHandler := kithttp.NewServer(
    makeListTasksEndpoint(s),
    decodeListTasksRequest,
    encodeResponse,
    opts...,
  )

  r := mux.NewRouter()

  r.Handle("/v1/tasks", addTaskHandler).Methods("POST")
  r.Handle("/v1/tasks", listTasksHandler).Methods("GET")

  return r
}

/*
func MakeHTTPHandler(s Service, logger kitlog.Logger) http.Handler{
  r := mux.NewRouter()
  e := MakeServerEndpoints(s)
  opts := []kithttp.ServerOption{
          kithttp.ServerErrorLogger(logger),
          kithttp.ServerErrorEncoder(encodeError),
  }

  // POST: /tasks/    Add a new task
  // GET:  /tasks/    Get all tasks

  r.Methods("POST").Path("/tasks/").Handler(kithttp.NewServer(
    e.AddTaskEndpoint,
    decodeAddTaskRequest,
    encodeResponse,
    opts...,
  ))
  r.Methods("GET").Path("/tasks/").Handler(kithttp.NewServer(
    e.ListTasksEndpoint,
    decodeListTasksRequest,
    encodeResponse,
    opts...,
  ))

  return r
}*/
/*
func NewHttpClient(instance string, logger kitlog.Logger) (Service, error) {
  if !strings.HasPrefix(instance, "http") {
    instance = "http://" + instance
  }
  u, err := url.Parse(instance)
  if err != nil {
    return nil, err
  }

  limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

  var addTaskEndpoint endpoint.Endpoint
  {
    addTaskEndpoint = kithttp.NewClient(
      "POST",
      copyURL(u, "/tasks"),
      encodeHTTPGenericRequest,
      decodeAddTaskResponse,
    ).Endpoint()
    addTaskEndpoint = limiter(addTaskEndpoint)
    addTaskEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
      Name: "tasks",
      Timeout: 30 * time.Second,
    }))(addTaskEndpoint)
  }

  var listTasksEndpoint endpoint.Endpoint
  {
    listTasksEndpoint = kithttp.NewClient(
      "GET",
      copyURL(u, "/tasks"),
      encodeHTTPGenericRequest,
      decodeListTasksResponse,
    ).Endpoint()
    listTasksEndpoint = limiter(listTasksEndpoint)
    listTasksEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
      Name: "tasks",
      Timeout: 30 * time.Second,
    }))(listTasksEndpoint)
  }

  return Endpoints{
    AddTaskEndpoint: addTaskEndpoint,
    ListTasksEndpoint: listTasksEndpoint,
  }, nil
}*/

func copyURL(base *url.URL, path string) *url.URL {
  next := *base
  next.Path = path
  return &next
}

var errBadRoute = errors.New("bad route")

func decodeAddTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
  var body struct {
    Name            string    `json:"name"`
    TaskType        string    `json:"taskType"`
    StartTime       time.Time `json:"startTime"`
    EndTime         time.Time `json:"endTime"`
    ExecutionStatus string    `json:"executionStatus"`
  }

  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    return nil, err
  }

  return addTaskRequest{
    Name: body.Name,
    TaskType: body.TaskType,
    StartTime: body.StartTime,
    EndTime: body.EndTime,
    ExecutionStatus: body.ExecutionStatus,
  }, nil
}

func decodeAddTaskResponse(_ context.Context, r *http.Response) (interface{}, error) {
  if r.StatusCode != http.StatusOK {
    return nil, errors.New(r.Status)
  }
  var resp addTaskResponse
  err := json.NewDecoder(r.Body).Decode(&resp)
  return resp, err
}

func decodeListTasksRequest(_ context.Context, r *http.Request) (interface{}, error) {
  return listTasksRequest{}, nil
}

func decodeListTasksResponse(_ context.Context, r *http.Response) (interface{}, error) {
  if r.StatusCode != http.StatusOK {
    return nil, errors.New(r.Status)
  }
  var resp listTasksResponse
  err := json.NewDecoder(r.Body).Decode(&resp)
  return resp, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
  if e, ok := response.(errorer); ok && e.error() != nil {
    encodeError(ctx, e.error(), w)
    return nil
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  return json.NewEncoder(w).Encode(response)
}

type errorer interface {
  error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  switch err {
  case ErrUnknown:
    w.WriteHeader(http.StatusNotFound)
  case ErrInvalidArgument:
    w.WriteHeader(http.StatusBadRequest)
  default:
    w.WriteHeader(http.StatusInternalServerError)
  }
  json.NewEncoder(w).Encode(map[string]interface{}{
    "error": err.Error(),
  })
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(request); err != nil {
    return err
  }
  r.Body = ioutil.NopCloser(&buf)
  return nil
}
