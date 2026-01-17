package constants

import "time"

const (
	TOTPSecretExpiry = 10 * time.Minute
	TOTPInterval     = 30 // second
	TOTPFailureLimit = 10
	TOTPFailureReset = 5 * time.Minute

	MuxConnection = 1

	LikeTarType_Video   = 1
	LikeTarType_Comment = 2
	LikeAction_Like     = 1
	LikeAction_Unlike   = 2

	VideoVisitQueueSize                   = 5000
	VideoVisitBatchSize                   = 1000
	VideoVisitFlushInterval time.Duration = 10 * time.Second
)
