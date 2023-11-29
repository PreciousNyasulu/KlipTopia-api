package rabbitmqprocesses

import (
	conf "kliptopia-api/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var config = conf.LoadConfig()
var	logger = *log.New()

func CreateSessionQueue(queueName string) (*amqp.Channel , error){
	conn,err := ConnectChannel(config.RabbitMQ.Url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	channel,err := conn.Channel()	
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

func ConnectChannel(connectionString string) (*amqp.Connection, error){
	// Connect to RabbitMQ
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		logger.Error("Failed to connect to RabbitMQ: ", err)
		return nil, err
	}
	return conn, nil
}