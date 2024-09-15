package subscription

import "mqtt-http-bridge/src/datastore"

func subscriptionToStore(sub Subscription) datastore.SubscriptionRecord {
	return datastore.SubscriptionRecord{
		ID:           sub.ID,
		Name:         sub.Name,
		Topic:        sub.Topic,
		Extract:      sub.Extract,
		Filter:       sub.Filter,
		URL:          sub.URL,
		Method:       sub.Method,
		Headers:      sub.Headers,
		BodyTemplate: sub.BodyTemplate,

		SubscriptionTemplateID:         sub.SubscriptionTemplateID,
		SubscriptionTemplateParameters: sub.SubscriptionTemplateParameters,
	}
}

func subscriptionFromStore(sub datastore.SubscriptionRecord) Subscription {
	return Subscription{
		ID:           sub.ID,
		Name:         sub.Name,
		Topic:        sub.Topic,
		Extract:      sub.Extract,
		Filter:       sub.Filter,
		URL:          sub.URL,
		Method:       sub.Method,
		Headers:      sub.Headers,
		BodyTemplate: sub.BodyTemplate,

		SubscriptionTemplateID:         sub.SubscriptionTemplateID,
		SubscriptionTemplateParameters: sub.SubscriptionTemplateParameters,
	}
}

func subscriptionTemplateToStore(sub SubscriptionTemplate) datastore.SubscriptionTemplateRecord {
	return datastore.SubscriptionTemplateRecord{
		SubscriptionRecord: subscriptionToStore(sub.Subscription),
		RequiredParameters: sub.RequiredParameters,
	}
}

func subscriptionTemplateFromStore(sub datastore.SubscriptionTemplateRecord) SubscriptionTemplate {
	return SubscriptionTemplate{
		Subscription:       subscriptionFromStore(sub.SubscriptionRecord),
		RequiredParameters: sub.RequiredParameters,
	}
}
