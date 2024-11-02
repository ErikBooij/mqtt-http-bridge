package datastore

import (
	"errors"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Store interface {
	// Subscriptions

	AddSubscription(subscription SubscriptionRecord) (SubscriptionRecord, error)
	GetSubscription(id string) (SubscriptionRecord, error)
	GetSubscriptions() ([]SubscriptionRecord, error)
	UpdateSubscription(subscription SubscriptionRecord) (SubscriptionRecord, error)
	DeleteSubscription(id string) error

	// Global Variables

	SetGlobalParameter(key string, value any) error
	GetGlobalParameters() (map[string]any, error)
	DeleteGlobalParameter(key string) error
}

type SubscriptionRecord struct {
	// Name is the name of the subscription
	Name string `json:"name"`

	// ID is the unique identifier for the subscription
	ID string `json:"id"`
	// Topic is the MQTT topic the subscription is for
	Topic string `json:"topic"`
	// Extract is a map of variable names to JSONata expressions
	Extract map[string]string `json:"extract"`
	// Filter is a JSONata expression to filter messages, returning true if the message should be processed
	Filter string `json:"filter"`

	// Method is the HTTP method to use for the request
	Method string `json:"method"`
	// URL is the URL to send the HTTP request to
	URL string `json:"URL"`
	// Headers is a map of headers to include in the request
	Headers map[string]string `json:"headers"`
	// Body is the template to use for rendering the HTTP response body
	Body string `json:"template"`
}
