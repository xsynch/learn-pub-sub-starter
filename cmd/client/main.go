package main

import (
	"fmt"
	"log"


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

	gs := gamelogic.NewGameState(username)
	for {
		commands := gamelogic.GetInput()
		cmd := commands[0]
		switch {
		case cmd == "spawn":
			log.Printf("Spawning a new military with: %s\n",commands[1:])
			err = gs.CommandSpawn(commands)
			
			if err != nil {
				log.Printf("Error spawning army: %s\n",err)
			}
		case cmd == "move":
			log.Printf("Moving the army\n")
			_,err = gs.CommandMove(commands)
			if err != nil {
				log.Printf("Error moving the army: %s",err)
			}
		case cmd == "help":
			gamelogic.PrintClientHelp()
		case cmd == "status":
			gs.CommandStatus()
		case cmd == "spam":
			log.Printf("Spamming not allowed yet\n")
		case cmd == "quit":
			gamelogic.PrintQuit()
			return 
		default:
			log.Printf("Command not found\n")

		}
	}



	
	// wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// fmt.Println("Starting Peril client...")
}
