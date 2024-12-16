package main

import (
	
	
	"log"
	"os"
	"os/signal"
	
	

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xsynch/learn-pub-sub-starter/internal/pubsub"
	"github.com/xsynch/learn-pub-sub-starter/internal/routing"
)

func main() {
	connString := "amqp://guest:guest@localhost:5672/"
	
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", connString, err)
	}
	defer conn.Close()
	log.Printf("Successfully connected to Rabbit MQ Server\n")
	rabbitmqChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating rabbit mq channel: %s",err)
	}
	defer rabbitmqChan.Close()
	pstate := routing.PlayingState{IsPaused: true}

	err = pubsub.PublishJSON(rabbitmqChan, string(routing.ExchangePerilDirect),string(routing.PauseKey),pstate)
	if err != nil {
		log.Fatalf("Error with PublishJSON: %s",err)
	}
	

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
