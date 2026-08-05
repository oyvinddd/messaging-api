package app

const (
	// register device token for an account
	registerTokenRoute 	= "POST /api/v1/push/tokens"
	// delete all device tokens for an account
	deleteTokensRoute 	= "DELETE /api/v1/push/tokens"
	// delete one device token for an account
	deleteTokenRoute	= "DELETE /api/v1/push/tokens/{id}"
	// send an email
	sendEmailRoute 		= "POST /api/v1/internal/email"
	// send a push message
	sendPushRoute		= "POST /api/v1/internal/push"

)

