namespace go api

include "user.thrift"
include "video.thrift"
include "interaction.thrift"
include "social.thrift"

service UserApi {
    user.RegisterResp  Register(1: required user.RegisterReq req) (api.post="/user/register")
    user.LoginResp     Login(1: required user.LoginReq req) (api.post="/user/login")
    user.InfoResp      GetInfo(1: required user.InfoQuery req) (api.get="/user/info")
    user.AvatarResp    UploadAvatar(1: required user.AvatarReq req) (api.put="/user/avatar/upload")
    user.MFAQrcodeResp MFAQrcode(1: required user.MFAQrcodeReq req) (api.get="/auth/mfa/qrcode")
    user.MFABindResp   MFABind(1: required user.MFABindReq req) (api.post="/auth/mfa/bind")
}

service VideoApi {
    video.FeedResp    Feed(1: required video.FeedQuery req) (api.get="/video/feed")
    video.PublishResp Publish(1: required video.PublishReq req) (api.post="/video/publish")
    video.ListResp    List(1: required video.ListQuery req) (api.get="/video/list")
    video.PopularResp Popular(1: required video.PopularQuery req) (api.get="/video/popular")
    video.SearchResp  Search(1: required video.SearchReq req) (api.post="/video/search")
    video.VisitResp   Visit(1: required video.VisitQuery req) (api.get="/video/:vid")
}

service InteractionApi {
    interaction.PublishLikeResp PublishLike(1: required interaction.PublishLikeReq req) (api.post="/like/action")
    interaction.ListLikeResp   ListLike(1: required interaction.ListLikeQuery req) (api.get="/like/list")
    interaction.PublishCommentResp PublishComment(1: required interaction.PublishCommentReq req) (api.post="/comment/publish")
    interaction.ListCommentResp ListComment(1: required interaction.ListCommentQuery query) (api.get="/comment/list")
    interaction.DeleteCommentResp DeleteComment(1: required interaction.DeleteCommentReq req) (api.delete="/comment/delete")
}

service SocialApi {
    social.FollowResp        Follow(1: required social.FollowReq req) (api.post="/relation/action")
    social.ListFollowsResp   ListFollows(1: required social.ListFollowsQuery req) (api.get="/following/list")
    social.ListFollowersResp ListFollowers(1: required social.ListFollowersQuery req) (api.get="/follower/list")
    social.ListFriendsResp   ListFriends(1: required social.ListFriendsQuery req) (api.get="/friends/list")
}