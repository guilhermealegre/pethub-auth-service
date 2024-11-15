package producer

const (

	// Notification queue
	EmailTemplateSignupConfirmationKey = "signup_confirmation"
	SignupConfirmationCode             = "signup_confirmation_code"
	NotificationExchange               = "notification-exchange"
	NotificationQueueSendEmailBidding  = "notification.send.email"
)

type NotificationSendEmailRequest struct {
	To           []string       `json:"to"`
	Template     string         `json:"template"`
	PlaceHolders map[string]any `json:"place_holders"`
}
