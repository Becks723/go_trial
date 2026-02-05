package service

import (
	"StreamCore/kitex_gen/user"
	"StreamCore/pkg/constants"
	"errors"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func (s *UserService) MFABind(uid uint, req *user.MFABindReq) error {
	var err error

	// 1. check if current user has bound
	// 2. take pending secret (if not found, ask for another qrcode req)
	// 3. 检验 code（totp.ValidateCustom）
	// 3.1 防重放。一个code有效期为30s，30s内不允许重复提交
	// 3.2 防爆破。每个用户设置失败次数限制（如10次/5分钟）
	// 4. 正式绑定 secret （绑定状态和secret存db）

	u, _ := s.db.GetById(uid)
	if u.TOTPBound {
		return nil
	}

	pending, err := s.cache.GetTOTPPending(s.ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get pending secret: %w", err)
	}

	failCount, err := s.cache.TOTPFailureCount(s.ctx, uid)
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
		if s.cache.IsCurrentTOTPPeriodLocked(s.ctx, uid) {
			return errors.New("replay detected")
		} else {
			// 防重放机制 - 记录
			err = s.cache.LockCurrentTOTPPeriod(s.ctx, uid)
			if err != nil {
				return err
			}

			// save secret to db
			err = s.db.UpdateTOTPSecret(uid, pending)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// validation fails
	// increase failure
	err = s.cache.IncreaseTOTPFailure(s.ctx, uid)
	if err != nil {
		return err
	}

	return errors.New("failed to bind mfa, try to refresh the qrcode")
}
