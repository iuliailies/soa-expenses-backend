package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Notification struct {
	UserID          int    `json:"user_id"`          
	Message         string `json:"message"`          
	CurrentExpenses int    `json:"current_expenses"` 
	Limit           int    `json:"limit"`            
}

type NotificationPublisher interface {
	Publish(notification Notification) error 
}

// RabbitMQPublisher is an implementation of NotificationPublisher using RabbitMQ
type RabbitMQPublisher struct {
	conn    *amqp.Connection // Connection to RabbitMQ
	channel *amqp.Channel    // Channel to communicate with RabbitMQ
	queue   amqp.Queue       // Queue to which notifications will be published
}

func NewRabbitMQPublisher(rabbitMQURL string) (*RabbitMQPublisher, error) {
	// Establish connection to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err 
	}

	// Open a channel for communication
	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // Clean up connection on error
		return nil, err
	}

	// Declare the queue where messages will be sent
	queue, err := ch.QueueDeclare(
		"notifications_queue", // Queue name
		true,                  // Durable (survives RabbitMQ restarts)
		false,                 // Auto-delete when unused
		false,                 // Not exclusive to a single connection
		false,                 // No-wait for confirmation
		nil,                   // Additional queue arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Return the initialized publisher
	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

// Publish sends a notification to the RabbitMQ queue
func (p *RabbitMQPublisher) Publish(notification Notification) error {
	// Serialize the notification into JSON format
	body, err := json.Marshal(notification)
	if err != nil {
		return err 
	}

	// Publish the message to the RabbitMQ queue
	err = p.channel.Publish(
		"",           // Default exchange (direct routing to a queue)
		p.queue.Name, // Queue name as the routing key
		false,        // Mandatory flag (not used here)
		false,        // Immediate flag (not used here)
		amqp.Publishing{
			ContentType: "application/json", // Specify content type
			Body:        body,              // Message body
		},
	)
	if err != nil {
		return err // Return error if publishing fails
	}

	// Log the published notification for debugging
	log.Printf("Notification published: %+v", notification)
	return nil
}

// Close releases RabbitMQ resources
func (p *RabbitMQPublisher) Close() {
	p.channel.Close() // Close the channel
	p.conn.Close()    // Close the connection
}
