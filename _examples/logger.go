package main

import (
	"github.com/liutme/gorabbit"
	"log"
)

type CustomLogger struct {
	gorabbit.Logger
}

func (c CustomLogger) Info(msg string, args ...any) {
	log.Println(msg, args)
}

func (c CustomLogger) Warn(msg string, args ...any) {
	log.Println(msg, args)
}

func (c CustomLogger) Error(msg string, args ...any) {
	log.Println(msg, args)
}

func main() {
	rabbitClient := &gorabbit.Client{
		Config: gorabbit.ConnectionConfig{
			Host:     "127.0.0.1",
			Port:     "5672",
			UserName: "admin",
			Password: "admin",
			VHost:    "/",
		},
		Consumers: []gorabbit.IConsumer{
			&Consume2{},
		},
		Logger: CustomLogger{},
	}

	rabbitClient.Init()
	forever := make(chan bool)
	<-forever
}
