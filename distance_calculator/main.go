package main

import (
	"log"
)

const kafkaTopic = "obudata"

// Attach transport to the business logic(HTTP, GRPC, Kafka)
func main() {
	var (
		err error
		svc Calculator
	)

	svc = NewCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
