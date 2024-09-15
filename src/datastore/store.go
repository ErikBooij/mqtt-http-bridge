package datastore

import (
	"errors"
)

var (
	ErrSubscriptionNotFound    = errors.New("subscription not found")
	ErrSubscriptionIDConflicts = errors.New("subscription ID conflicts with existing subscription")

	ErrSubscriptionTemplateNotFound    = errors.New("subscription template not found")
	ErrSubscriptionTemplateIDConflicts = errors.New("subscription template ID conflicts with existing subscription")
)

type Store interface {
	// Subscriptions

	AddSubscription(subscription SubscriptionRecord) (SubscriptionRecord, error)
	GetSubscription(id string) (SubscriptionRecord, error)
	GetSubscriptions() ([]SubscriptionRecord, error)
	UpdateSubscription(subscription SubscriptionRecord) (SubscriptionRecord, error)
	DeleteSubscription(id string) error

	AddSubscriptionTemplate(subscriptionTemplate SubscriptionTemplateRecord) (SubscriptionTemplateRecord, error)
	GetSubscriptionTemplate(id string) (SubscriptionTemplateRecord, error)
	GetSubscriptionTemplates() ([]SubscriptionTemplateRecord, error)
	UpdateSubscriptionTemplate(subscriptionTemplate SubscriptionTemplateRecord) (SubscriptionTemplateRecord, error)
	DeleteSubscriptionTemplate(id string) error

	// Global Variables

	SetGlobalParameter(key string, value any) error
	GetGlobalParameters() (map[string]any, error)
	DeleteGlobalParameter(key string) error
}

type SubscriptionRecord struct {
	// Name is the name of the subscription
	Name string `json:"name"`

	// SubscriptionTemplateID is the optional ID of the template this subscription was derived from
	SubscriptionTemplateID *string `json:"subscriptionTemplateId"`
	// SubscriptionTemplateParameters is a map of parameters to use when deriving a subscription from a template
	SubscriptionTemplateParameters map[string]any `json:"subscriptionTemplateParameters"`

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
	// BodyTemplate is the template to use for rendering the HTTP response body
	BodyTemplate string `json:"template"`
}

type SubscriptionTemplateRecord struct {
	SubscriptionRecord

	// RequiredParameters is a list of parameters that must be provided when deriving a subscription from a template.
	RequiredParameters []string `json:"requiredParameters"`
}
