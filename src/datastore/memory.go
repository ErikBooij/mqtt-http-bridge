package datastore

import (
	"sync"
)

// Ensure fileStore implements the Store interface.
var _ Store = &memoryStore{}

type memoryStore struct {
	globalParameters   map[string]any
	globalParametersMu sync.RWMutex
	subscriptions      map[string]SubscriptionRecord
	subscriptionsMu    sync.RWMutex
}

func Memory() (Store, error) {
	return &memoryStore{
		globalParameters: make(map[string]any),
		subscriptions:    make(map[string]SubscriptionRecord),
	}, nil
}

func (s *memoryStore) AddSubscription(sub SubscriptionRecord) (SubscriptionRecord, error) {
	s.subscriptionsMu.Lock()
	defer s.subscriptionsMu.Unlock()

	s.subscriptions[sub.ID] = sub

	return sub, nil
}

func (s *memoryStore) GetSubscription(id string) (SubscriptionRecord, error) {
	s.subscriptionsMu.RLock()
	defer s.subscriptionsMu.RUnlock()

	sub, ok := s.subscriptions[id]

	if !ok {
		return SubscriptionRecord{}, ErrSubscriptionNotFound
	}

	return sub, nil
}

func (s *memoryStore) GetSubscriptions() ([]SubscriptionRecord, error) {
	s.subscriptionsMu.RLock()
	defer s.subscriptionsMu.RUnlock()

	subscriptions := make([]SubscriptionRecord, 0, len(s.subscriptions))

	for _, sub := range s.subscriptions {
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (s *memoryStore) UpdateSubscription(sub SubscriptionRecord) (SubscriptionRecord, error) {
	s.subscriptionsMu.Lock()
	defer s.subscriptionsMu.Unlock()

	if _, ok := s.subscriptions[sub.ID]; !ok {
		return SubscriptionRecord{}, ErrSubscriptionNotFound
	}

	s.subscriptions[sub.ID] = sub
	return sub, nil
}

func (s *memoryStore) DeleteSubscription(id string) error {
	s.subscriptionsMu.Lock()
	defer s.subscriptionsMu.Unlock()

	if _, ok := s.subscriptions[id]; !ok {
		return ErrSubscriptionNotFound
	}

	delete(s.subscriptions, id)
	return nil
}

func (s *memoryStore) SetGlobalParameter(key string, value any) error {
	s.globalParametersMu.Lock()
	defer s.globalParametersMu.Unlock()

	if value == "" {
		delete(s.globalParameters, key)
	} else {
		s.globalParameters[key] = value
	}

	return nil
}

func (s *memoryStore) GetGlobalParameters() (map[string]any, error) {
	s.globalParametersMu.RLock()
	defer s.globalParametersMu.RUnlock()

	globalVariables := make(map[string]any, len(s.globalParameters))

	for key, value := range s.globalParameters {
		globalVariables[key] = value
	}

	return globalVariables, nil
}

func (s *memoryStore) DeleteGlobalParameter(key string) error {
	s.globalParametersMu.Lock()
	defer s.globalParametersMu.Unlock()

	delete(s.globalParameters, key)
	return nil
}
