package constants

import "time"

const (
	TOTPSecretExpiry = 10 * time.Minute
	TOTPInterval     = 30 // second
	TOTPFailureLimit = 10
	TOTPFailureReset = 5 * time.Minute
)
