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
		GlobalParameters: make(map[string]any),
		Subscriptions:    make(map[string]SubscriptionRecord),

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
	GlobalParameters map[string]any                `json:"globalParameters"`
	Subscriptions    map[string]SubscriptionRecord `json:"subscriptions"`

	globalParametersMu sync.RWMutex
	subscriptionsMu    sync.RWMutex

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

	return nil
}
