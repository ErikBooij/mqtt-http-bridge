package dev

import (
	"log"
	"mqtt-http-bridge/src/subscription"
)

func PopulateDataStore(service subscription.Service, logger *log.Logger) error {
	if err := service.Reset(); err != nil {
		return err
	}

	sub1, err := addSubscription(service)

	if err != nil {
		return err
	}

	logger.Printf("Added subscription with ID: %s", sub1.ID)

	if err := addGlobalVariable(service, "authToken", "abcdef"); err != nil {
		return err
	}

	return nil
}

func addGlobalVariable(service subscription.Service, name, value string) error {
	return service.SetGlobalParameter(name, value)
}

func addSubscription(service subscription.Service) (subscription.Subscription, error) {
	sub, err := service.AddSubscription(subscription.Subscription{
		Name:  "Shortcut Button 001",
		Topic: "zigbee2mqtt/shortcut-button-001",
		Extract: map[string]string{
			"action":  "action",
			"battery": "battery",
		},
		Filter: "custom.action='1_short_release'",
		Body:   `{"action":"{{.custom.action}}","battery":"{{.custom.battery}}"}`,
		URL:    "https://straight-application-12.webhook.cool",
		Method: "PATCH",
		Headers: map[string]string{
			"Authorization": "Bearer 123",
			"Content-Type":  "application/json",
		},
	})

	if err != nil {
		return subscription.Subscription{}, err
	}

	return sub, err
}
