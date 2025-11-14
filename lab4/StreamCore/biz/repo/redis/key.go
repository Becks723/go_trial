package redis

import "fmt"

const (
	VideoRankKey = "zVideoRank"
)

func VideoLikeKey(vid uint) string {
	return fmt.Sprintf("video_like:%d", vid)
}

func CommentLikeKey(cid uint) string {
	return fmt.Sprintf("comment_like:%d", cid)
}
