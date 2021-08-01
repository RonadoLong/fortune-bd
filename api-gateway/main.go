package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wq-fotune-backend/api-gateway/handler"
	"wq-fotune-backend/libs/logger"
)

func main() {
	eng := handler.InitEngine()
	port := "0.0.0.0:9530"
	s := &http.Server{
		Addr:           port,
		Handler:        eng,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	handleSigterm(func() {
		logger.Infof("%v", "do something when server exit")
	})

	logger.Infof("Listening server at ====== %s", port)
	if err := s.ListenAndServe(); err == nil {
		log.Println("Listening server err ", port)
	}
}

func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
