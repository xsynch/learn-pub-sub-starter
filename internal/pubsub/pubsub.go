package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	
)


func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{
	data,err := json.Marshal(val)
	if err != nil {
		return err 
	}
	newMsg := amqp.Publishing{
		ContentType: "application/json",
		Body: data,
	}
	err = ch.PublishWithContext(context.Background(),exchange,key,false,false,newMsg)
	if err != nil {
		return err 
	}
	return nil 
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	var durability bool 
	var ad bool
	var ex bool
	tChannel, err := conn.Channel()
	if err != nil {
		return nil,amqp.Queue{}, err 
	}
	if simpleQueueType == 1 {
		durability = true
		ad =  false
		ex = false 

	} else {
		durability = false 
		ad = true 
		ex = true 
	}
	tQueue, err := tChannel.QueueDeclare(queueName,durability,ad, ex,false,nil)
	if err != nil {
		return nil, amqp.Queue{},err
	}
	err = tChannel.QueueBind(queueName,key,exchange,false,nil)
	if err != nil {
		return nil,amqp.Queue{},nil 
	}
	return tChannel, tQueue,nil 
	
}