package main

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer
var deliveryChan chan kafka.Event

// Connect to Kafka in a persistent way
func Connect(args map[string]interface{}) *kafka.Producer {
	if producer != nil {
		return producer
	}
	// extract broker list from map
	brokers := ""
	for i, s := range args["kafka_brokers_sasl"].([]interface{}) {
		if i == 0 {
			brokers = s.(string)
		} else {
			brokers += "," + s.(string)
		}
	}
	// generate configuration
	config := kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"security.protocol": "sasl_ssl",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     args["user"],
		"sasl.password":     args["password"],
	}
	// create a producer and return it
	p, err := kafka.NewProducer(&config)
	if err != nil {
		log.Println(err)
		return nil
	}
	producer = p
	deliveryChan = make(chan kafka.Event)
	return producer
}

// Send a message
func Send(p *kafka.Producer, topic string, msg []byte) error {
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, deliveryChan)
	e := <-deliveryChan
	m := e.(*kafka.Message)
	return m.TopicPartition.Error
}

// mkErr makes an error from different sources
func mkErr(message interface{}) map[string]interface{} {
	result := "unknown error"
	switch v := message.(type) {
	case string:
		result = v
	case error:
		result = v.Error()
	}
	return map[string]interface{}{
		"error": result,
	}
}

// Main is the entry point
func Main(args map[string]interface{}) map[string]interface{} {
	// retrieving the connection
	p := Connect(args)
	if p == nil {
		return mkErr("cannot connect")
	}
	// getting the topic
	t, ok := args["topic"].(string)
	if !ok {
		return mkErr("no topic")
	}
	m, ok := args["message"]
	if !ok {
		return mkErr("no message")
	}
	// getting the message
	msg, err := json.Marshal(m)
	if err != nil {
		return mkErr(err)
	}
	// sending the message
	err = Send(p, t, msg)
	if err != nil {
		return mkErr(err)
	}
	return map[string]interface{}{
		"ok": true,
	}
}
