package subscription

import (
	"errors"
	"fmt"
	"mqtt-http-bridge/src/datastore"
	"mqtt-http-bridge/src/utilities"
	"regexp"
	"slices"
	"strings"
)

var (
	ErrMissingRequiredParametersForTemplate         = errors.New("missing required parameters for template")
	ErrUnableToHydrateTemplatedSubscriptionProperty = errors.New("unable to hydrate templated subscription property")
	ErrInvalidGlobalParameterKey                    = errors.New("invalid key")
)

type Service interface {
	AddSubscription(subscription Subscription) (Subscription, error)
	GetSubscription(id string) (Subscription, error)
	GetSubscriptions() ([]Subscription, error)
	UpdateSubscription(subscription Subscription) (Subscription, error)
	DeleteSubscription(id string) error

	AddSubscriptionTemplate(subscriptionTemplate SubscriptionTemplate) (SubscriptionTemplate, error)
	GetSubscriptionTemplate(id string) (SubscriptionTemplate, error)
	GetSubscriptionTemplates() ([]SubscriptionTemplate, error)
	UpdateSubscriptionTemplate(subscriptionTemplate SubscriptionTemplate) (SubscriptionTemplate, error)
	DeleteSubscriptionTemplate(id string) error

	SetGlobalParameter(key string, value string) error
	DeleteGlobalParameter(key string) error
	GetGlobalParameters() (map[string]any, error)

	AddSubscriptionFromTemplate(subscriptionTemplateID string, parameters map[string]any) (Subscription, error)
	GetSubscriptionsForTopic(topic string) ([]Subscription, error)

	ApplyPlaceholdersOnSubscription(sub Subscription, params map[string]any) (Subscription, error)

	// Reset removes everything from the store, mostly only used for development purposes.
	Reset() error
}

func NewService(store datastore.Store) Service {
	return &service{
		store: store,

		topicMatcher: newTopicMatcher(),
	}
}

type service struct {
	store datastore.Store

	topicMatcher *topicMatcher
}

func (s *service) AddSubscription(subscription Subscription) (Subscription, error) {
	subscription.ID = utilities.GenerateRandomID()

	sub, err := s.store.AddSubscription(subscriptionToStore(subscription))

	if err != nil {
		return Subscription{}, err
	}

	return subscriptionFromStore(sub), nil
}

func (s *service) GetSubscription(id string) (Subscription, error) {
	sub, err := s.store.GetSubscription(id)

	if err != nil {
		return Subscription{}, err
	}

	return subscriptionFromStore(sub), nil
}

func (s *service) GetSubscriptions() ([]Subscription, error) {
	subscriptions := make([]Subscription, 0)

	subs, err := s.store.GetSubscriptions()

	if err != nil {
		return subscriptions, err
	}

	for _, sub := range subs {
		converted := subscriptionFromStore(sub)

		if hydrated, err := s.hydrateTemplatedSubscription(converted); err != nil {
			subscriptions = append(subscriptions, converted)
		} else {
			subscriptions = append(subscriptions, hydrated)
		}
	}

	slices.SortStableFunc(subscriptions, func(a, b Subscription) int {
		// Sort by name first, but if they're for whatever reason the same, sort by ID
		if name := strings.Compare(a.Name, b.Name); name != 0 {
			return name
		}

		return strings.Compare(a.ID, b.ID)
	})

	return subscriptions, nil
}

func (s *service) UpdateSubscription(subscription Subscription) (Subscription, error) {
	sub, err := s.store.UpdateSubscription(subscriptionToStore(subscription))

	if err != nil {
		return Subscription{}, err
	}

	return subscriptionFromStore(sub), nil
}

func (s *service) DeleteSubscription(id string) error {
	return s.store.DeleteSubscription(id)
}

func (s *service) AddSubscriptionTemplate(subscriptionTemplate SubscriptionTemplate) (SubscriptionTemplate, error) {
	subscriptionTemplate.ID = utilities.GenerateRandomID()

	sub, err := s.store.AddSubscriptionTemplate(subscriptionTemplateToStore(subscriptionTemplate))

	if err != nil {
		return SubscriptionTemplate{}, err
	}

	return subscriptionTemplateFromStore(sub), nil
}

func (s *service) GetSubscriptionTemplate(id string) (SubscriptionTemplate, error) {
	sub, err := s.store.GetSubscriptionTemplate(id)

	if err != nil {
		return SubscriptionTemplate{}, err
	}

	return subscriptionTemplateFromStore(sub), nil
}

func (s *service) GetSubscriptionTemplates() ([]SubscriptionTemplate, error) {
	subscriptionTemplates := make([]SubscriptionTemplate, 0)

	subTemps, err := s.store.GetSubscriptionTemplates()

	if err != nil {
		return subscriptionTemplates, err
	}

	for _, subTemp := range subTemps {
		converted := subscriptionTemplateFromStore(subTemp)

		subscriptionTemplates = append(subscriptionTemplates, converted)
	}

	slices.SortStableFunc(subscriptionTemplates, func(a, b SubscriptionTemplate) int {
		// Sort by name first, but if they're for whatever reason the same, sort by ID
		if name := strings.Compare(a.Name, b.Name); name != 0 {
			return name
		}

		return strings.Compare(a.ID, b.ID)
	})

	return subscriptionTemplates, nil
}

func (s *service) UpdateSubscriptionTemplate(subscriptionTemplate SubscriptionTemplate) (SubscriptionTemplate, error) {
	sub, err := s.store.UpdateSubscriptionTemplate(subscriptionTemplateToStore(subscriptionTemplate))

	if err != nil {
		return SubscriptionTemplate{}, err
	}

	subs, err := s.GetSubscriptions()

	if err != nil {
		return SubscriptionTemplate{}, err
	}

	templateSubs := utilities.FilterSlice(subs, func(sub Subscription) bool {
		return sub.SubscriptionTemplateID != nil && *sub.SubscriptionTemplateID == subscriptionTemplate.ID
	})

	for _, sub := range templateSubs {
		sub.Name = subscriptionTemplate.Name
		sub.Topic = subscriptionTemplate.Topic
		sub.Extract = subscriptionTemplate.Extract
		sub.Filter = subscriptionTemplate.Filter
		sub.Method = subscriptionTemplate.Method
		sub.URL = subscriptionTemplate.URL
		sub.Headers = subscriptionTemplate.Headers
		sub.BodyTemplate = subscriptionTemplate.BodyTemplate

		if _, err := s.UpdateSubscription(sub); err != nil {
			return SubscriptionTemplate{}, err
		}
	}

	return subscriptionTemplateFromStore(sub), nil
}

func (s *service) DeleteSubscriptionTemplate(id string) error {
	return s.store.DeleteSubscriptionTemplate(id)
}

var globalParameterKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)

func (s *service) SetGlobalParameter(key string, value string) error {
	if !globalParameterKeyRegex.MatchString(key) {
		return fmt.Errorf("%w: %s", ErrInvalidGlobalParameterKey, key)
	}

	return s.store.SetGlobalParameter(key, value)
}

func (s *service) GetGlobalParameters() (map[string]any, error) {
	return s.store.GetGlobalParameters()
}

func (s *service) DeleteGlobalParameter(key string) error {
	return s.store.DeleteGlobalParameter(key)
}

func (s *service) AddSubscriptionFromTemplate(subscriptionTemplateID string, parameters map[string]any) (Subscription, error) {
	template, err := s.store.GetSubscriptionTemplate(subscriptionTemplateID)

	if err != nil {
		return Subscription{}, err
	}

	subTemp := subscriptionTemplateFromStore(template)

	for _, value := range subTemp.RequiredParameters {
		missing := make([]string, 0, len(subTemp.RequiredParameters))

		if _, ok := parameters[value]; !ok {
			missing = append(missing, value)
		}

		if len(missing) > 0 {
			return Subscription{}, fmt.Errorf("%w: %s", ErrMissingRequiredParametersForTemplate, strings.Join(missing, ","))
		}
	}

	newSubscriptionInput := Subscription{
		Name: subTemp.Name,

		Topic:   subTemp.Topic,
		Extract: subTemp.Extract,
		Filter:  subTemp.Topic,

		Method:       subTemp.Method,
		URL:          subTemp.URL,
		Headers:      subTemp.Headers,
		BodyTemplate: subTemp.BodyTemplate,

		SubscriptionTemplateID:         &subscriptionTemplateID,
		SubscriptionTemplateParameters: parameters,
	}

	return s.AddSubscription(newSubscriptionInput)
}

func (s *service) GetSubscriptionsForTopic(topic string) ([]Subscription, error) {
	subscriptions := make([]Subscription, 0)

	subs, err := s.store.GetSubscriptions()

	if err != nil {
		return subscriptions, err
	}

	for _, sub := range subs {
		if s.topicMatcher.match(topic, sub.Topic) {
			converted := subscriptionFromStore(sub)

			if hydrated, err := s.hydrateTemplatedSubscription(converted); err != nil {
				subscriptions = append(subscriptions, converted)
			} else {
				subscriptions = append(subscriptions, hydrated)
			}
		}
	}

	return subscriptions, nil
}

func (s *service) Reset() error {
	// Delete all subscriptions
	subs, err := s.store.GetSubscriptions()

	if err != nil {
		return err
	}

	for _, sub := range subs {
		if err := s.store.DeleteSubscription(sub.ID); err != nil {
			return err
		}
	}

	// Delete all subscription templates
	subsTemp, err := s.store.GetSubscriptionTemplates()

	if err != nil {
		return err
	}

	for _, subTemp := range subsTemp {
		if err := s.store.DeleteSubscriptionTemplate(subTemp.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) ApplyPlaceholdersOnSubscription(sub Subscription, params map[string]any) (Subscription, error) {
	if params == nil {
		params = make(map[string]any)
	}

	params["tpl"] = sub.SubscriptionTemplateParameters

	var err error

	subClone, _ := utilities.DeepCopy(sub)

	if subClone.Name, err = utilities.RenderInlineTemplate(subClone.Name, params); err != nil {
		return Subscription{}, fmt.Errorf("%w name: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, err)
	}

	if subClone.Topic, err = utilities.RenderInlineTemplate(subClone.Topic, params); err != nil {
		return Subscription{}, fmt.Errorf("%w topic: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, err)
	}

	if subClone.Filter, err = utilities.RenderInlineTemplate(subClone.Filter, params); err != nil {
		return Subscription{}, fmt.Errorf("%w filter: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, err)
	}

	if subClone.BodyTemplate, err = utilities.RenderInlineTemplate(subClone.BodyTemplate, params); err != nil {
		return Subscription{}, fmt.Errorf("%w body template: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, err)
	}

	if subClone.Method, err = utilities.RenderInlineTemplate(subClone.Method, params); err != nil {
		return Subscription{}, fmt.Errorf("%w method: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, err)
	}

	if subClone.URL, err = utilities.RenderInlineTemplate(subClone.URL, params); err != nil {
		return Subscription{}, fmt.Errorf("%w url: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, err)
	}

	for key, value := range subClone.Headers {
		if subClone.Headers[key], err = utilities.RenderInlineTemplate(value, params); err != nil {
			return Subscription{}, fmt.Errorf("%w header %s: %w", ErrUnableToHydrateTemplatedSubscriptionProperty, key, err)
		}
	}

	return subClone, nil
}

func (s *service) hydrateTemplatedSubscription(sub Subscription) (Subscription, error) {
	return s.ApplyPlaceholdersOnSubscription(sub, nil)
}
