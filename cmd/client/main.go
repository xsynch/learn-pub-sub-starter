package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xsynch/learn-pub-sub-starter/internal/pubsub"
	"github.com/xsynch/learn-pub-sub-starter/internal/routing"
	"github.com/xsynch/learn-pub-sub-starter/internal/gamelogic"
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
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error with the client welcome: %s",err)		
	}

	_, _,err = pubsub.DeclareAndBind(conn,routing.ExchangePerilDirect,fmt.Sprintf("%s.%s",routing.PauseKey,username),routing.PauseKey,0)
	if err != nil {
		log.Fatalf("Error with Declare and Bind: %s",err)
	}



	
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Starting Peril client...")
}
