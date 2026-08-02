package push

const (
	// Android the token is an Android token
	Android Platform = "android"
	// IOS the token is an iOS push token
	IOS Platform = "ios"
)

type (
	Platform string

	Token struct {
		// Value is the push token
		Value string `json:"push_token"`
		// Platform the token belongs to
		Platform Platform `json:"platform"`
	}
)

