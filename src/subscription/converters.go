package subscription

import "mqtt-http-bridge/src/datastore"

func subscriptionToStore(sub Subscription) datastore.SubscriptionRecord {
	return datastore.SubscriptionRecord{
		ID:      sub.ID,
		Name:    sub.Name,
		Topic:   sub.Topic,
		Extract: sub.Extract,
		Filter:  sub.Filter,
		URL:     sub.URL,
		Method:  sub.Method,
		Headers: sub.Headers,
		Body:    sub.Body,
	}
}

func subscriptionFromStore(sub datastore.SubscriptionRecord) Subscription {
	return Subscription{
		ID:      sub.ID,
		Name:    sub.Name,
		Topic:   sub.Topic,
		Extract: sub.Extract,
		Filter:  sub.Filter,
		URL:     sub.URL,
		Method:  sub.Method,
		Headers: sub.Headers,
		Body:    sub.Body,
	}
}
