package push

import (
	"log"
	"context"
	"github.com/google/uuid"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"firebase.google.com/go/v4/messaging"
)

const (
	// FCM the token is an Android token
	FCM Platform = "fcm"
	// APNS the token is an iOS push token
	APNS Platform = "apns"
)

type (
	Platform string

	DeviceToken struct {
		// ID the token id
		ID uuid.UUID `json:"-"`
		// Value is the push token
		Value string `json:"device_token"`
		// Platform the token belongs to
		Platform Platform `json:"platform"`
	}

	Message struct {
		// RecipientID the account that will get the push
		RecipientID uuid.UUID
		// Title the title of the push message
		Title string `json:"title"`
		// Body the content of the push message
		Body string `json:"body"`
	}

	Provider interface {
		Send(ctx context.Context, tokens []DeviceToken, message Message) error
	}

	firebaseProvider struct {
		client *messaging.Client
	}
)

func NewDeviceToken(id uuid.UUID, token string, platform Platform) *DeviceToken {
	return &DeviceToken{ID: id, Value: token, Platform: platform}
}

func NewFirebaseProvider(ctx context.Context, credentialsFilePath string) Provider {
	opt := option.WithCredentialsFile(credentialsFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase provider: %v\n", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing Firebase provider: %v\n", err)
	}
	return &firebaseProvider{client: client}
}

func (p *firebaseProvider) Send(ctx context.Context, tokens []DeviceToken, message Message) error {
	fcmMessages := fcmMessagesForTokens(tokens, message)
	br, err := p.client.SendEach(ctx, fcmMessages)
	if err != nil {
		return err
	}
	for i, r := range br.Responses {
		if !r.Success {
			// TODO: delete token on our side
			log.Printf("failed to send to %s: %v", fcmMessages[i].Token, r.Error)
		}
	}
	return nil
}

func fcmMessagesForTokens(deviceTokens []DeviceToken, message Message) []*messaging.Message {
	messages := make([]*messaging.Message, 0)
	for _, token := range deviceTokens {
		fcmMessage := &messaging.Message{
			Token: token.Value,
			Notification: &messaging.Notification{
				Title: message.Title,
				Body: message.Body,
			},
		}
		messages = append(messages, fcmMessage)
	}
	return messages
}

