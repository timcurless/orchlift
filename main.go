package main

import (
  "flag"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  "github.com/go-kit/kit/log"
  "github.com/timcurless/orchlift/pkg"
)

const (
  defaultPort = "8080"
)

func main() {
  RunServer()
}

func RunServer() {
  var (
    addr = envString("PORT", defaultPort)
    httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
  )

  flag.Parse()

  var logger log.Logger
  logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
  logger = log.With(logger, "ts", log.DefaultTimestampUTC)

  var (
    tasks = orchlift.NewTaskRepository()
  )

  var s orchlift.Service
  s = orchlift.NewService(tasks)

  httpLogger := log.With(logger, "component", "http")

  mux := http.NewServeMux()

  mux.Handle("/v1/", orchlift.MakeHandler(s, httpLogger))

  http.Handle("/", mux)

  errs := make(chan error, 2)
  go func() {
    logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
    errs <- http.ListenAndServe(*httpAddr, nil)
  }()
  go func() {
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT)
    errs <- fmt.Errorf("%s", <-c)
  }()

  logger.Log("terminated", <-errs)
}

func envString(env, failback string) string {
  e := os.Getenv(env)
  if e == "" {
    return failback
  }
  return e
}
