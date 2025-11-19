package redis

import "fmt"

const (
	VideoRankKey = "zVideoRank"
)

func VideoVisitCountKey(vid uint) string {
	return fmt.Sprintf("video_visit_count:%d", vid)
}

func VideoLikeKey(vid uint) string {
	return fmt.Sprintf("video_like:%d", vid)
}

func CommentLikeKey(cid uint) string {
	return fmt.Sprintf("comment_like:%d", cid)
}

func UserLikeVidKey(uid uint) string {
	return fmt.Sprintf("user_like_vid:%d", uid)
}

func FollowsCountKey(uid uint) string {
	return fmt.Sprintf("follows_count:%d", uid)
}

func FollowersCountKey(uid uint) string {
	return fmt.Sprintf("followers_count:%d", uid)
}
