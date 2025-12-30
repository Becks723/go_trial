package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/common"
	"StreamCore/biz/model/user"
	"StreamCore/biz/repo"
	cache "StreamCore/biz/repo/cache/user"
	"StreamCore/pkg/constants"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
	"StreamCore/pkg/util/jwt"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type UserService struct {
	repo  repo.UserRepo
	cache *cache.UserCache
}

func NewUserService(repo repo.UserRepo, cache *cache.UserCache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) Register(ctx context.Context, req *user.RegisterReq) (err error) {
	// check potential duplicated username
	if _, err = s.repo.GetByUsername(req.Username); err == nil {
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
	if err = s.repo.Create(&do); err != nil {
		return
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, req *user.LoginReq) (data *common.UserInfo, auth *common.AuthenticationInfo, err error) {
	// find user in db
	u, err := s.repo.GetByUsername(req.Username)
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
	atoken, err := jwt.GenerateAccessToken(u, ev.AccessToken_Secret, jwt.HoursOf(ev.AccessToken_ExpiryHours))
	if err != nil {
		return
	}
	rtoken, err := jwt.GenerateRefreshToken(u, ev.RefreshToken_Secret, jwt.HoursOf(ev.RefreshToken_ExpiryHours))
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

func (s *UserService) GetInfo(ctx context.Context, query *user.InfoQuery) (data *common.UserInfo, err error) {
	// convert string id to uint
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		err = errors.New("Bad uid format.")
		return
	}

	// find user in db
	u, err := s.repo.GetById(uint(uid))
	if err != nil {
		return
	}

	data = domain2Dto(u)
	err = nil
	return
}

func (s *UserService) UploadAvatar(ctx context.Context, fileHeader *multipart.FileHeader) (data *common.UserInfo, err error) {
	var (
		localPrefix  = "./uploads"
		accessPrefix = "/static"
	)
	curUid := retrieveUid(ctx)

	if !util.IsValidImage(fileHeader) {
		err = errors.New("Bad image format.")
		return
	}

	// exceeds image limit
	limit := env.Instance().IO_ImageSizeLimit
	if fileHeader.Size > util.ToByte(limit) {
		err = fmt.Errorf("Exceeds image size limit (current %dmb but limits %dmb)", limit, util.ToMb(fileHeader.Size))
		return
	}

	// save image locally
	dst := fmt.Sprintf(localPrefix+accessPrefix+"/avatars/%d_%d.png", // TODO: match extensions
		curUid, time.Now().Unix())
	err = util.SaveFile(fileHeader, dst)
	if err != nil {
		return
	}

	// update db
	newUrl, _ := strings.CutPrefix(dst, localPrefix)
	u, err := s.repo.UpdateAvatar(curUid, newUrl)
	if err != nil {
		return
	}

	// make resp
	data = domain2Dto(u)
	return
}

func (s *UserService) MFAQrcode(ctx context.Context, req *user.MFAQrcodeReq) (*user.MFAQrcodeResp_Data, error) {
	var err error

	uid, _ := util.RetrieveUserId(ctx)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "StreamCore",
		AccountName: strconv.FormatUint(uint64(uid), 10),
	})
	if err != nil {
		return nil, err
	}
	secret := key.Secret()

	w := env.Instance().MFA_QrcodeWidth
	h := env.Instance().MFA_QrcodeHeight
	img, err := key.Image(w, h)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)                                    // encode img to png (binary)
	qrcode := base64.StdEncoding.EncodeToString(buf.Bytes()) // base64

	// cache secret
	err = s.cache.SetTOTPPending(ctx, uid, secret)
	if err != nil {
		return nil, err
	}

	data := new(user.MFAQrcodeResp_Data)
	data.Secret = secret
	data.Qrcode = qrcode
	return data, nil
}

func (s *UserService) MFABind(ctx context.Context, req *user.MFABindReq) error {
	var err error

	// 1. check if current user has bound
	// 2. take pending secret (if not found, ask for another qrcode req)
	// 3. 检验 code（totp.ValidateCustom）
	// 3.1 防重放。一个code有效期为30s，30s内不允许重复提交
	// 3.2 防爆破。每个用户设置失败次数限制（如10次/5分钟）
	// 4. 正式绑定 secret （绑定状态和secret存db）
	uid, _ := util.RetrieveUserId(ctx)

	u, _ := s.repo.GetById(uid)
	if u.TOTPBound {
		return nil
	}

	pending, err := s.cache.GetTOTPPending(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get pending secret: %w", err)
	}

	failCount, err := s.cache.TOTPFailureCount(ctx, uid)
	if err != nil {
		return err
	}
	if failCount > constants.TOTPFailureLimit { // 防爆破
		return errors.New("failure exceeds limit, please try again later")
	}
	success, err := totp.ValidateCustom(req.Code, pending, time.Now(), totp.ValidateOpts{
		Period:    constants.TOTPInterval,
		Algorithm: otp.AlgorithmSHA1,
		Digits:    6,
		Skew:      1,
	})
	if success { // validation ok
		// 防重放机制 - 判断
		if s.cache.IsCurrentTOTPPeriodLocked(ctx, uid) {
			return errors.New("replay detected")
		} else {
			// 防重放机制 - 记录
			err = s.cache.LockCurrentTOTPPeriod(ctx, uid)
			if err != nil {
				return err
			}

			// save secret to db
			err = s.repo.UpdateTOTPSecret(uid, pending)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// validation fails
	// increase failure
	err = s.cache.IncreaseTOTPFailure(ctx, uid)
	if err != nil {
		return err
	}

	return errors.New("failed to bind mfa, request for a new qrcode and try again")
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
