package main

import (
	
	
	"log"
	"os"
	"os/signal"
	
	

	amqp "github.com/rabbitmq/amqp091-go"
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

	

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
