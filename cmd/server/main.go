package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xsynch/learn-pub-sub-starter/internal/gamelogic"
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
	gamelogic.PrintServerHelp()
	
	rabbitmqChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating rabbit mq channel: %s",err)
	}
	defer rabbitmqChan.Close()
	
	_, _,err = pubsub.DeclareAndBind(conn,routing.ExchangePerilDirect,"game_logs",fmt.Sprintf("%s.*",routing.GameLogSlug),1)
	if err != nil {
		log.Fatalf("Error creating queue: %s",err)
	}


	if err != nil {
		log.Fatalf("Error with Declare and Bind: %s",err)
	}



	for {
		msg := gamelogic.GetInput()
		switch {
		case msg[0] == "pause":
			log.Printf("Sending a pause message\n")
			pstate := routing.PlayingState{IsPaused: true}
			err = pubsub.PublishJSON(rabbitmqChan, string(routing.ExchangePerilDirect),string(routing.PauseKey),pstate)
			if err != nil {
				log.Fatalf("Error with PublishJSON: %s",err)
			}
		case msg[0] == "resume": 
			log.Printf("Sending a resume message\n")
			pstate := routing.PlayingState{IsPaused: false}
			err = pubsub.PublishJSON(rabbitmqChan, string(routing.ExchangePerilDirect),string(routing.PauseKey),pstate)
			if err != nil {
				log.Fatalf("Error with PublishJSON: %s",err)
			}
		case msg[0] == "quit":
			return 
		default:
			log.Printf("Please enter a valid command\n")		
			
		}
	}
	// pstate := routing.PlayingState{IsPaused: true}

	// err = pubsub.PublishJSON(rabbitmqChan, string(routing.ExchangePerilDirect),string(routing.PauseKey),pstate)
	// if err != nil {
	// 	log.Fatalf("Error with PublishJSON: %s",err)
	// }
	

	// wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
}
