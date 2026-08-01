package mailer

type (
	Mailer interface {
		Send(recipient, subject, title string) error
	}

	postmarkMailer struct {
		apiKey string
	}
)

func NewPostmarkMailer(apiKey string) Mailer {
	return &postmarkMailer{apiKey: apiKey}
}

