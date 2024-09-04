package datastore

import (
	"slices"
	"strings"
	"sync"
	"zigbee-coordinator/src/subscription"
	"zigbee-coordinator/src/utilities"
)

// Ensure staticStore implements the Store interface.
var _ Store = &staticStore{}

type staticStore struct {
	mu            sync.RWMutex
	subscriptions map[string]subscription.Subscription
}

func Memory() Store {
	return &staticStore{
		subscriptions: make(map[string]subscription.Subscription, 0),
	}
}

func (s *staticStore) AddSubscription(sub subscription.Subscription) (subscription.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub.ID = utilities.GenerateRandomID()

	s.subscriptions[sub.ID] = sub

	return sub, nil
}

func (s *staticStore) DeleteSubscription(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.subscriptions[id]; !ok {
		return ErrorSubscriptionNotFound
	}

	delete(s.subscriptions, id)
	return nil
}

func (s *staticStore) GetSubscriptionForTopic(topic string) (subscription.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.subscriptions {
		if sub.Topic == topic {
			return sub, nil
		}
	}

	return subscription.Subscription{}, ErrorSubscriptionNotFound
}

func (s *staticStore) GetSubscriptions() ([]subscription.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subscriptions := make([]subscription.Subscription, 0, len(s.subscriptions))

	for _, sub := range s.subscriptions {
		subscriptions = append(subscriptions, sub)
	}

	slices.SortStableFunc(subscriptions, func(a, b subscription.Subscription) int {
		return strings.Compare(a.ID, b.ID)
	})

	return subscriptions, nil
}

func (s *staticStore) UpdateSubscription(sub subscription.Subscription) (subscription.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.subscriptions[sub.ID]; !ok {
		return subscription.Subscription{}, ErrorSubscriptionNotFound
	}

	s.subscriptions[sub.ID] = sub
	return sub, nil
}
