package push_email

type config struct {
	Email       string `yaml:"email"`
	AppPassword string `yaml:"app_password" mapstructure:"app_password"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
}

type PlaceHoldersEmailSignupConfirmation struct {
	ConfirmationCode string `json:"confirmation_code"`
}
