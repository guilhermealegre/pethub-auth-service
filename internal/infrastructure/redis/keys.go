package redis

import (
	"fmt"
	"strings"
)

type KeyFormat string

const (
	SignupConfirmationCode KeyFormat = "env.%s.auth.signup_confirmation_code.%s"
)

func (kf KeyFormat) Format(args ...any) string {
	return fmt.Sprintf(string(kf), kf.prepareKeys(args...)...)
}

func (kf KeyFormat) prepareKeys(args ...any) []any {
	for i, arg := range args {
		if str, ok := arg.(string); ok {
			args[i] = strings.ReplaceAll(strings.TrimSpace(str), " ", "_")
		}
	}

	return args

}
