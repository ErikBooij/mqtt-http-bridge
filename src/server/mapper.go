package server

import (
	"mqtt-http-bridge/src/subscription"
)

type subscriptionResponse struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Topic   string            `json:"topic"`
	Extract map[string]string `json:"extract,omitempty"`
	Filter  string            `json:"filter,omitempty"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

func subscriptionToResponse(sub subscription.Subscription) any {
	return subscriptionResponse{
		ID:      sub.ID,
		Name:    sub.Name,
		Topic:   sub.Topic,
		Extract: sub.Extract,
		Filter:  sub.Filter,
		Method:  sub.Method,
		URL:     sub.URL,
		Headers: sub.Headers,
		Body:    sub.Body,
	}
}
