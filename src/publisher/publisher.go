package publisher

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"zigbee-coordinator/src/subscription"
)

type Publisher interface {
	Publish(body []byte, subscription subscription.Subscription)
}

type publisher struct {
	jobs chan publisherJob
}

type publisherJob struct {
	body         []byte
	subscription subscription.Subscription
}

func New(ctx context.Context, parallel int, clientFactory func() *http.Client) *publisher {
	p := &publisher{
		jobs: make(chan publisherJob, 100),
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
	req, err := http.NewRequest(job.subscription.Method, job.subscription.URL, bytes.NewReader(job.body))
	if err != nil {
		log.Printf("Error creating request for subscription %s: %s\n", job.subscription.ID, err)
		return
	}

	for k, v := range job.subscription.Headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error publishing message to subscription %s: %s\n", job.subscription.ID, err)
		return
	}

	defer resp.Body.Close()
}
