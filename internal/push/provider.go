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
		Send(ctx context.Context, deviceToken string, message Message) error
	}

	firebaseProvider struct {
		client *messaging.Client
	}
)

func NewFirebaseProvider(ctx context.Context, credentialsPath string) Provider {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("error initializing Firebase provider: %v\n", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing Firebase provider: %v\n", err)
	}

	return &firebaseProvider{client: client}
}

func (p *firebaseProvider) Send(ctx context.Context, deviceToken string, message Message) error {
	fcmMessage := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: message.Title,
			Body: message.Body,
		},
	}
	_, err := p.client.Send(ctx, fcmMessage)
	return err
}

