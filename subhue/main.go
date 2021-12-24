package main

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/pubsub"
)

var (
	projectID      = flag.String("project_id", "", "project ID")
	subscriptionID = flag.String("subscription_id", "", "subscription ID")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	log.Printf("Creating client")
	client, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating subscription")
	sub := client.Subscription(*subscriptionID)
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Received message: %q", string(m.Data))
		m.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
}
