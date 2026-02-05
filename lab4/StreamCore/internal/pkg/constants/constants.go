package constants

import "time"

const (
	ApiServiceName         = "api"
	UserServiceName        = "user"
	VideoServiceName       = "video"
	InteractionServiceName = "interaction"
	SocialServiceName      = "social"

	MaxConnections = 1000
	MaxQPS         = 100

	JWT_AccessSecret           = "access_token_secret"
	JWT_RefreshSecret          = "refresh_token_secret"
	JWT_AccessTokenExpiration  = 12 * time.Hour
	JWT_RefreshTokenExpiration = 7 * 24 * time.Hour

	MFA_QrcodeWidth  = 256
	MFA_QrcodeHeight = 256

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

	FollowAction_Follow   = 0
	FollowAction_Unfollow = 1
	SocialCacheExpiration = 30 * time.Minute
)
