package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"log"
	"os"
)

type Message struct {
	CharacterID int64 `json:"characterID"`
	AccessToken string `json:"accessToken"`
}

func PublishMessage(dao *db.DB, projectID string, topicEnvKey string, characterID int64, accessToken string) error {
	topicID := os.Getenv(topicEnvKey)

	if topicID == "" {
		return fmt.Errorf("%s must be set", topicEnvKey)
	}

	encryptedToken, err := dao.Encrypt(accessToken)

	if err != nil {
		return fmt.Errorf("dao.Encrypt: %v", err)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)

	message := Message{
		CharacterID: characterID,
		AccessToken: encryptedToken,
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