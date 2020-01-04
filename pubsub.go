package esi

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type EsiMessage struct {
	CharacterID int64 `json:"characterID"`
	AccessToken string `json:"accessToken"`
}

func PublishMessage(projectID string, topicID string, characterID int64, accessToken string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)

	message := EsiMessage{
		CharacterID: characterID,
		AccessToken: accessToken,
	}

	data, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}

	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("result.Get: %v", err)
	}
	log.Printf( "Published message; msg ID: %v\n", id)
	return nil
}