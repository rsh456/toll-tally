package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rsh456/toll-tally/types"
	"github.com/sirupsen/logrus"
)

type kafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService Calculator
}

func NewKafkaConsumer(topic string, svc Calculator) (*kafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &kafkaConsumer{
		consumer:    c,
		calcService: svc,
	}, nil
}

func (c *kafkaConsumer) Start() {
	logrus.Info("Starting Kafka Consumer")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *kafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("kafka consumer: unable to unmarshal message: %v", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("kafka consumer: unable to calculate distance: %v", err)
			continue
		}
		fmt.Printf("distance %.2f\n", distance)
	}
}
