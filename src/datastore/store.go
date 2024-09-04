package datastore

import (
	"errors"
	"zigbee-coordinator/src/subscription"
)

var (
	ErrorSubscriptionNotFound = errors.New("subscription not found")
)

type Store interface {
	AddSubscription(subscription subscription.Subscription) (subscription.Subscription, error)
	DeleteSubscription(id string) error
	GetSubscriptionForTopic(topic string) (subscription.Subscription, error)
	GetSubscriptions() ([]subscription.Subscription, error)
	UpdateSubscription(subscription subscription.Subscription) (subscription.Subscription, error)
}
