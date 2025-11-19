package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/user"
	"StreamCore/biz/repo"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

type UserService struct {
	repo repo.UserRepo
}

func NewUserService(repo repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (serv *UserService) Register(ctx context.Context, req *user.RegisterReq) (err error) {
	// check potential duplicated username
	if _, err = serv.repo.GetByUsername(req.Username); err == nil {
		err = errors.New("用户名已存在") // TODO: i18n
		return
	}

	// encrypt password
	var encrypted string
	if encrypted, err = util.EncryptPassword(req.Password); err != nil {
		return
	}

	// repo stuff
	do := domain.User{
		Username:  req.Username,
		Password:  encrypted,
		AvatarUrl: "", // TODO: offer default avatar
	}
	if err = serv.repo.Create(&do); err != nil {
		return
	}

	return nil
}

func (serv *UserService) Login(ctx context.Context, req *user.LoginReq) (data *common.UserInfo, auth *common.AuthenticationInfo, err error) {
	// find user in db
	u, err := serv.repo.GetByUsername(req.Username)
	if err != nil {
		return
	}

	// password correct?
	if !util.CheckPassword(req.Password, u.Password) {
		err = errors.New("密码错误") // TODO: i18n
		return
	}

	// generate access, refresh tokens
	ev := env.Instance()
	atoken, err := util.GenerateAccessToken(u, ev.AccessToken_Secret, util.HoursOf(ev.AccessToken_ExpiryHours))
	if err != nil {
		return
	}
	rtoken, err := util.GenerateRefreshToken(u, ev.RefreshToken_Secret, util.HoursOf(ev.RefreshToken_ExpiryHours))
	if err != nil {
		return
	}

	data = domain2Dto(u)
	auth = &common.AuthenticationInfo{
		AccessToken:  atoken,
		RefreshToken: rtoken,
	}
	err = nil
	return
}

func (serv *UserService) GetInfo(ctx context.Context, query *user.InfoQuery) (data *common.UserInfo, err error) {
	// convert string id to uint
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		err = errors.New("Bad uid format.")
		return
	}

	// find user in db
	u, err := serv.repo.GetById(uint(uid))
	if err != nil {
		return
	}

	data = domain2Dto(u)
	err = nil
	return
}

func (serv *UserService) UploadAvatar(ctx context.Context, fileHeader *multipart.FileHeader) (data *common.UserInfo, err error) {
	var (
		localPrefix  = "./uploads"
		accessPrefix = "/static"
	)
	curUid := retrieveUid(ctx)

	if !isValidImage(fileHeader) {
		err = errors.New("Bad image format.")
		return
	}

	// save image locally
	dst := fmt.Sprintf(localPrefix+accessPrefix+"/avatars/%d_%d.png", // TODO: match extensions
		curUid, time.Now().Unix())
	err = saveFile(fileHeader, dst)
	if err != nil {
		return
	}

	// update db
	newUrl, _ := strings.CutPrefix(dst, localPrefix)
	u, err := serv.repo.UpdateAvatar(curUid, newUrl)
	if err != nil {
		return
	}

	// make resp
	data = domain2Dto(u)
	return
}

func domain2Dto(u *domain.User) *common.UserInfo {
	return &common.UserInfo{
		Id:        util.Uint2String(u.Id),
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		DeletedAt: util.TimePtr2String(u.DeletedAt),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
}

func retrieveUid(ctx context.Context) uint {
	obj := ctx.Value("uid")
	uid, _ := obj.(uint)
	return uid
}
