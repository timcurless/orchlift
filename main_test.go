package main

import (
  "net/http"
  "testing"

  "github.com/stretchr/testify/suite"
)

// Setup

type MainTestSuite struct {
  suite.Suite
}

func (s *MainTestSuite) SetupTest() {

}

// Run Server

func (s *MainTestSuite) Test_RunServer_InvokesListenAndServe() {
  actual := ""
  httpListenAndServe := func(addr string, handler http.Handler) error {
    actual = addr
    return nil
  }

  RunServer()

  s.Equal(":8080", actual)
}

// Suite

func TestMainSuite(t *testing.T) {
  logFatalOrig := logFatal
  defer func() { logFatal = logFatalOrig }()
  logFatal = func(v ...interface{}) {}
  logPrintfOrig := logPrintf
  defer func() { logPrintf = logPrintfOrig }()
  logPrintf = func(format string, v ...interface{}) {}
  httpListenAndServeOrig := httpListenAndServe
  defer func() { httpListenAndServe = httpListenAndServeOrig }()
  httpListenAndServe = func(addr string, handler http.handler) error { return nil }
  suite.Run(t, new(MainTestSuite))
}
