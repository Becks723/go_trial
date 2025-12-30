package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (s *UserService) UploadAvatar(uid uint, data []byte) (*common.UserInfo, error) {
	var err error
	var (
		localPrefix  = "./uploads"
		accessPrefix = "/static"
	)
	curUid := uid

	if !util.IsValidImage(data) {
		return nil, errors.New("bad image format")
	}

	// exceeds image limit
	limit := env.Instance().IO_ImageSizeLimit
	size := len(data)
	if size > util.ToByte(limit) {
		return nil, fmt.Errorf("exceeds image size limit (current %.2fmb but limits %.2fmb)", util.ToMb(size), limit)
	}

	// save image locally
	dst := fmt.Sprintf(localPrefix+accessPrefix+"/avatars/%d_%d.png", // TODO: match extensions
		curUid, time.Now().Unix())
	err = util.SaveFile(data, dst)
	if err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// update db
	newUrl, _ := strings.CutPrefix(dst, localPrefix)
	u, err := s.db.UpdateAvatar(curUid, newUrl)
	if err != nil {
		return nil, fmt.Errorf("error updating database avatar: %w", err)
	}

	// make resp
	user := &common.UserInfo{
		Id:        util.Uint2String(u.Id),
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		DeletedAt: util.TimePtr2String(u.DeletedAt),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
	return user, nil
}
