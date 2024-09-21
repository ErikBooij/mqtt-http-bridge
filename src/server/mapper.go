package server

import (
	"mqtt-http-bridge/src/subscription"
)

type subscriptionResponse struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Topic              string            `json:"topic"`
	Extract            map[string]string `json:"extract,omitempty"`
	Filter             string            `json:"filter,omitempty"`
	Method             string            `json:"method"`
	URL                string            `json:"url"`
	Headers            map[string]string `json:"headers,omitempty"`
	BodyTemplate       string            `json:"bodyTemplate,omitempty"`
	TemplateID         *string           `json:"templateId,omitempty"`
	TemplateParameters map[string]any    `json:"templateParameters,omitempty"`
}

func subscriptionToResponse(sub subscription.Subscription) any {
	return subscriptionResponse{
		ID:           sub.ID,
		Name:         sub.Name,
		Topic:        sub.Topic,
		Extract:      sub.Extract,
		Filter:       sub.Filter,
		Method:       sub.Method,
		URL:          sub.URL,
		Headers:      sub.Headers,
		BodyTemplate: sub.BodyTemplate,

		TemplateID:         sub.SubscriptionTemplateID,
		TemplateParameters: sub.SubscriptionTemplateParameters,
	}
}

type subscriptionTemplateResponse struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Topic              string            `json:"topic"`
	Extract            map[string]string `json:"extract,omitempty"`
	Filter             string            `json:"filter,omitempty"`
	Method             string            `json:"method"`
	URL                string            `json:"url"`
	Headers            map[string]string `json:"headers,omitempty"`
	BodyTemplate       string            `json:"bodyTemplate,omitempty"`
	RequiredParameters []string          `json:"requiredParameters,omitempty"`
}

func subscriptionTemplateToResponse(sub subscription.SubscriptionTemplate) any {
	return subscriptionTemplateResponse{
		ID:                 sub.ID,
		Name:               sub.Name,
		Topic:              sub.Topic,
		Extract:            sub.Extract,
		Filter:             sub.Filter,
		Method:             sub.Method,
		URL:                sub.URL,
		Headers:            sub.Headers,
		BodyTemplate:       sub.BodyTemplate,
		RequiredParameters: sub.RequiredParameters,
	}
}
