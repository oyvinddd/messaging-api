package app

const (
	// send a push message
	sendPushRoute		= "POST /api/v1/internal/push/send"
	// register device token for an account
	registerTokenRoute 	= "POST /api/v1/push/register"
	// delete all device tokens for an account
	deleteTokensRoute 	= "DELETE /api/v1/push/delete"
	// delete one device token for an account
	deleteTokenRoute	= "DELETE /api/v1/push/delete/{id}"
	// send an email
	sendEmailRoute 		= "POST /api/v1/internal/email/send"
)

