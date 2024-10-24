package publisher

import (
	"bytes"
	"context"
	"log"
	"mqtt-http-bridge/src/subscription"
	"net/http"
)

type Publisher interface {
	Publish(body []byte, subscription subscription.Subscription)
}

type publisher struct {
	jobs   chan publisherJob
	logger *log.Logger
}

type publisherJob struct {
	body         []byte
	subscription subscription.Subscription
}

func New(ctx context.Context, parallel int, clientFactory func() *http.Client, logger *log.Logger) *publisher {
	p := &publisher{
		jobs:   make(chan publisherJob, 100),
		logger: logger,
	}

	for range parallel {
		client := clientFactory()

		go p.start(ctx, client)
	}

	return p
}

func (p *publisher) Publish(body []byte, subscription subscription.Subscription) {
	p.jobs <- publisherJob{
		body:         body,
		subscription: subscription,
	}
}

func (p *publisher) start(ctx context.Context, client *http.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.jobs:
			p.doPublish(job, client)
		}
	}
}

func (p *publisher) doPublish(job publisherJob, client *http.Client) {
	p.logger.Printf("Publishing message to subscription %s (%s %s %s)\n", job.subscription.ID, job.subscription.Method, job.subscription.URL, job.body)

	req, err := http.NewRequest(job.subscription.Method, job.subscription.URL, bytes.NewReader(job.body))
	if err != nil {
		p.logger.Printf("Error creating request for subscription %s: %s\n", job.subscription.ID, err)
		return
	}

	for k, v := range job.subscription.Headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Subscription-ID", job.subscription.ID)
	req.Header.Add("Subscription-Name", job.subscription.Name)

	resp, err := client.Do(req)
	if err != nil {
		p.logger.Printf("Error publishing message to subscription %s: %s\n", job.subscription.ID, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		p.logger.Printf("Unexpected status code publishing message to subscription %s: %s\n", job.subscription.ID, resp.Status)
		return
	}
}
