package datastore

import (
	"sync"
)

// Ensure fileStore implements the Store interface.
var _ Store = &memoryStore{}

type memoryStore struct {
	globalParameters        map[string]any
	globalParametersMu      sync.RWMutex
	subscriptions           map[string]SubscriptionRecord
	subscriptionsMu         sync.RWMutex
	subscriptionTemplates   map[string]SubscriptionTemplateRecord
	subscriptionTemplatesMu sync.RWMutex
}

func Memory() (Store, error) {
	return &memoryStore{
		globalParameters:      make(map[string]any),
		subscriptions:         make(map[string]SubscriptionRecord),
		subscriptionTemplates: make(map[string]SubscriptionTemplateRecord),
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

func (s *memoryStore) AddSubscriptionTemplate(subTemp SubscriptionTemplateRecord) (SubscriptionTemplateRecord, error) {
	if err := s.ensureNoSubscriptionTemplateConflict(subTemp, true); err != nil {
		return SubscriptionTemplateRecord{}, err
	}

	s.subscriptionTemplatesMu.Lock()
	defer s.subscriptionTemplatesMu.Unlock()

	s.subscriptionTemplates[subTemp.ID] = subTemp

	return subTemp, nil
}

func (s *memoryStore) GetSubscriptionTemplate(id string) (SubscriptionTemplateRecord, error) {
	s.subscriptionTemplatesMu.RLock()
	defer s.subscriptionTemplatesMu.RUnlock()

	subTemp, ok := s.subscriptionTemplates[id]

	if !ok {
		return SubscriptionTemplateRecord{}, ErrSubscriptionTemplateNotFound
	}

	return subTemp, nil
}

func (s *memoryStore) GetSubscriptionTemplates() ([]SubscriptionTemplateRecord, error) {
	s.subscriptionTemplatesMu.RLock()
	defer s.subscriptionTemplatesMu.RUnlock()

	subscriptionTemplates := make([]SubscriptionTemplateRecord, 0, len(s.subscriptionTemplates))

	for _, subTemp := range s.subscriptionTemplates {
		subscriptionTemplates = append(subscriptionTemplates, subTemp)
	}

	return subscriptionTemplates, nil
}

func (s *memoryStore) UpdateSubscriptionTemplate(subTemp SubscriptionTemplateRecord) (SubscriptionTemplateRecord, error) {
	s.subscriptionTemplatesMu.Lock()
	defer s.subscriptionTemplatesMu.Unlock()

	if _, ok := s.subscriptionTemplates[subTemp.ID]; !ok {
		return SubscriptionTemplateRecord{}, ErrSubscriptionNotFound
	}

	s.subscriptionTemplates[subTemp.ID] = subTemp
	return subTemp, nil
}

func (s *memoryStore) DeleteSubscriptionTemplate(id string) error {
	s.subscriptionTemplatesMu.Lock()
	defer s.subscriptionTemplatesMu.Unlock()

	if _, ok := s.subscriptionTemplates[id]; !ok {
		return ErrSubscriptionTemplateNotFound
	}

	delete(s.subscriptionTemplates, id)
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

func (s *memoryStore) ensureNoSubscriptionConflict(sub SubscriptionRecord, isNew bool) error {
	subs, err := s.GetSubscriptions()

	if err != nil {
		return err
	}

	for _, existingSub := range subs {
		if !isNew && existingSub.ID == sub.ID {
			continue
		}

		if existingSub.ID == sub.ID {
			return ErrSubscriptionIDConflicts
		}
	}

	return nil
}

func (s *memoryStore) ensureNoSubscriptionTemplateConflict(sub SubscriptionTemplateRecord, isNew bool) error {
	subs, err := s.GetSubscriptionTemplates()

	if err != nil {
		return err
	}

	for _, existingSubTemp := range subs {
		if !isNew && existingSubTemp.ID == sub.ID {
			continue
		}

		if existingSubTemp.ID == sub.ID {
			return ErrSubscriptionTemplateIDConflicts
		}
	}

	return nil
}
