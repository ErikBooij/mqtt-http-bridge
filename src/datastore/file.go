package datastore

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

// Ensure fileStore implements the Store interface.
var _ Store = &fileStore{}

type fileStore struct {
	storage *storage
}

func File(filename string, reloadInterval time.Duration) (Store, error) {
	storage := &storage{
		GlobalParameters:      make(map[string]any),
		Subscriptions:         make(map[string]SubscriptionRecord),
		SubscriptionTemplates: make(map[string]SubscriptionTemplateRecord),

		filename: filename,
	}

	if err := storage.load(); err != nil {
		return nil, err
	}

	if err := storage.flush(); err != nil {
		return nil, err
	}

	go func() {
		for range time.Tick(reloadInterval) {
			if err := storage.load(); err != nil {
				log.Printf("Failed to reload file store: %v\n", err)
			}
		}
	}()

	return &fileStore{
		storage: storage,
	}, nil
}

func (s *fileStore) AddSubscription(sub SubscriptionRecord) (SubscriptionRecord, error) {
	defer s.storage.flush()

	s.storage.subscriptionsMu.Lock()
	defer s.storage.subscriptionsMu.Unlock()

	s.storage.Subscriptions[sub.ID] = sub

	return sub, nil
}

func (s *fileStore) GetSubscription(id string) (SubscriptionRecord, error) {
	s.storage.subscriptionsMu.RLock()
	defer s.storage.subscriptionsMu.RUnlock()

	sub, ok := s.storage.Subscriptions[id]

	if !ok {
		return SubscriptionRecord{}, ErrSubscriptionNotFound
	}

	return sub, nil
}

func (s *fileStore) GetSubscriptions() ([]SubscriptionRecord, error) {
	s.storage.subscriptionsMu.RLock()
	defer s.storage.subscriptionsMu.RUnlock()

	subscriptions := make([]SubscriptionRecord, 0, len(s.storage.Subscriptions))

	for _, sub := range s.storage.Subscriptions {
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (s *fileStore) UpdateSubscription(sub SubscriptionRecord) (SubscriptionRecord, error) {
	defer s.storage.flush()

	s.storage.subscriptionsMu.Lock()
	defer s.storage.subscriptionsMu.Unlock()

	if _, ok := s.storage.Subscriptions[sub.ID]; !ok {
		return SubscriptionRecord{}, ErrSubscriptionNotFound
	}

	s.storage.Subscriptions[sub.ID] = sub
	return sub, nil
}

func (s *fileStore) DeleteSubscription(id string) error {
	defer s.storage.flush()

	s.storage.subscriptionsMu.Lock()
	defer s.storage.subscriptionsMu.Unlock()

	if _, ok := s.storage.Subscriptions[id]; !ok {
		return ErrSubscriptionNotFound
	}

	delete(s.storage.Subscriptions, id)
	return nil
}

func (s *fileStore) AddSubscriptionTemplate(subTemp SubscriptionTemplateRecord) (SubscriptionTemplateRecord, error) {
	defer s.storage.flush()

	s.storage.subscriptionTemplatesMu.Lock()
	defer s.storage.subscriptionTemplatesMu.Unlock()

	s.storage.SubscriptionTemplates[subTemp.ID] = subTemp

	return subTemp, nil
}

func (s *fileStore) GetSubscriptionTemplate(id string) (SubscriptionTemplateRecord, error) {
	s.storage.subscriptionTemplatesMu.RLock()
	defer s.storage.subscriptionTemplatesMu.RUnlock()

	subTemp, ok := s.storage.SubscriptionTemplates[id]

	if !ok {
		return SubscriptionTemplateRecord{}, ErrSubscriptionTemplateNotFound
	}

	return subTemp, nil
}

func (s *fileStore) GetSubscriptionTemplates() ([]SubscriptionTemplateRecord, error) {
	s.storage.subscriptionTemplatesMu.RLock()
	defer s.storage.subscriptionTemplatesMu.RUnlock()

	subscriptionTemplates := make([]SubscriptionTemplateRecord, 0, len(s.storage.SubscriptionTemplates))

	for _, subTemp := range s.storage.SubscriptionTemplates {
		subscriptionTemplates = append(subscriptionTemplates, subTemp)
	}

	return subscriptionTemplates, nil
}

func (s *fileStore) UpdateSubscriptionTemplate(subTemp SubscriptionTemplateRecord) (SubscriptionTemplateRecord, error) {
	defer s.storage.flush()

	s.storage.subscriptionTemplatesMu.Lock()
	defer s.storage.subscriptionTemplatesMu.Unlock()

	if _, ok := s.storage.SubscriptionTemplates[subTemp.ID]; !ok {
		return SubscriptionTemplateRecord{}, ErrSubscriptionNotFound
	}

	s.storage.SubscriptionTemplates[subTemp.ID] = subTemp
	return subTemp, nil
}

func (s *fileStore) DeleteSubscriptionTemplate(id string) error {
	defer s.storage.flush()

	s.storage.subscriptionTemplatesMu.Lock()
	defer s.storage.subscriptionTemplatesMu.Unlock()

	if _, ok := s.storage.SubscriptionTemplates[id]; !ok {
		return ErrSubscriptionTemplateNotFound
	}

	delete(s.storage.SubscriptionTemplates, id)
	return nil
}

func (s *fileStore) SetGlobalParameter(key string, value any) error {
	defer s.storage.flush()

	s.storage.globalParametersMu.Lock()
	defer s.storage.globalParametersMu.Unlock()

	if value == "" {
		delete(s.storage.GlobalParameters, key)
	} else {
		s.storage.GlobalParameters[key] = value
	}

	return nil
}

func (s *fileStore) GetGlobalParameters() (map[string]any, error) {
	s.storage.globalParametersMu.RLock()
	defer s.storage.globalParametersMu.RUnlock()

	variables := make(map[string]any)

	for key, value := range s.storage.GlobalParameters {
		variables[key] = value
	}

	return variables, nil
}

func (s *fileStore) DeleteGlobalParameter(key string) error {
	defer s.storage.flush()
	s.storage.globalParametersMu.Lock()
	defer s.storage.globalParametersMu.Unlock()

	delete(s.storage.GlobalParameters, key)
	return nil
}

type storage struct {
	GlobalParameters      map[string]any                        `json:"globalParameters"`
	Subscriptions         map[string]SubscriptionRecord         `json:"subscriptions"`
	SubscriptionTemplates map[string]SubscriptionTemplateRecord `json:"subscriptionTemplates"`

	globalParametersMu      sync.RWMutex
	subscriptionsMu         sync.RWMutex
	subscriptionTemplatesMu sync.RWMutex

	filename string
	fsMu     sync.RWMutex
}

func (s *storage) flush() error {
	data, err := json.Marshal(s)

	if err != nil {
		return err
	}

	s.fsMu.Lock()
	defer s.fsMu.Unlock()

	return os.WriteFile(s.filename, data, 0644)
}

func (s *storage) load() error {
	s.fsMu.RLock()
	defer s.fsMu.RUnlock()

	data, err := os.ReadFile(s.filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return err
	}

	if s.GlobalParameters == nil {
		s.GlobalParameters = make(map[string]any)
	}

	if s.Subscriptions == nil {
		s.Subscriptions = make(map[string]SubscriptionRecord)
	}

	if s.SubscriptionTemplates == nil {
		s.SubscriptionTemplates = make(map[string]SubscriptionTemplateRecord)
	}

	return nil
}
