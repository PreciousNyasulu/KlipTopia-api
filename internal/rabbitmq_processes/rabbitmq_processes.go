package rabbitmqprocesses

import (
	conf "kliptopia-api/internal/config"
	"kliptopia-api/internal/models"

	"bytes"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var config = conf.LoadConfig()
var	logger = *log.New()
var channel *amqp.Channel

//rabbitmq_processes: establishes a queue with the user's username
func CreateSessionQueue(queueName string) (*amqp.Channel , error){
	conn,err := ConnectChannel()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	channel,err = conn.Channel()	
	if err != nil {
		logger.Error("Failed to initiate channel")
		return nil, err
	}
	defer channel.Close()

	//create queue
	_,err = channel.QueueDeclare(
		queueName,
		true, //durable
		false, // delete when unused
		false,	// exclusive
		false, // no-wait
		nil, //arguments
	)

	if err != nil {
		logger.Error("Failed to declare queue.")
		return nil, err
	}

	return channel, nil
}

func ConnectChannel() (*amqp.Connection, error){
	// Connect to RabbitMQ
	conn, err := amqp.Dial(config.RabbitMQ.Url)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ: ", err)
		return nil, err
	}
	return conn, nil
}

func PushMessageToQueue(queueName string, message *models.QueueMessage) error {
	// Encode message 
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(message); err != nil {
		logger.Error("Failed to encode message: ", err)
		return err
	}

	conn, _ := ConnectChannel()
	channel,err := conn.Channel()
	if err != nil {
		logger.Error("Failed to establish RabbitMQ channel.")
		return err
	}
	defer channel.Close()

	err = channel.Publish(
		"",         // exchange
		queueName,  // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        buf.Bytes(),
		},
	)

	if err != nil {
		logger.Error("Failed to publish message to the queue: ", err)
		return err
	}

	logger.Info("Message sent to the queue successfully.")
	return nil
}