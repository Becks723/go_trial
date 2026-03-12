package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/util/jwt"
)

// RefreshToken 刷新令牌
func (s *UserService) RefreshToken(token string) (*common.TokenInfo, error) {
	// 解析令牌
	tk, expiresAt, err := jwt.ParseToken[*domain.RefreshToken](token, constants.JWT_RefreshSecret)
	// 解析失败
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	// 令牌过期
	if expiresAt.UnixMilli() < time.Now().UnixMilli() {
		return nil, fmt.Errorf("token expired at %v", expiresAt)
	}
	// 校验令牌id
	tkId, err := s.db.GetTokenId(s.ctx, tk.Uid)
	if err != nil {
		return nil, fmt.Errorf("error GetTokenId: %w", err)
	}
	if tk.Id != tkId {
		return nil, fmt.Errorf("invalid token id")
	}

	// 生成新令牌，刷新令牌也要生成新的，
	// 旧的刷新令牌可以丢弃，也可以继续保存在db里，用来检测攻击者
	access, refresh, err := s.generateTokens(tk.Uid)
	if err != nil {
		return nil, err
	}
	data := new(common.TokenInfo)
	data.AccessToken = access
	data.RefreshToken = refresh
	return data, nil
}

// generateTokens 生成访问和刷新令牌
func (s *UserService) generateTokens(uid uint) (string, string, error) {
	// 生成访问令牌
	access, err := jwt.GenerateAccessToken(uid, constants.JWT_AccessSecret, constants.JWT_AccessTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed gen accessToken: %w", err)
	}

	// 生成刷新令牌
	tokenId := s.generateTokenId() // 生成令牌id
	refresh, err := jwt.GenerateRefreshToken(&domain.RefreshToken{
		Uid: uid,
		Id:  tokenId,
	}, constants.JWT_RefreshSecret, constants.JWT_RefreshTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed gen refreshToken: %w", err)
	}
	// 令牌id存入db
	if err = s.db.UpdateTokenId(s.ctx, uid, tokenId); err != nil {
		return "", "", fmt.Errorf("update token id failed: %w", err)
	}

	return access, refresh, nil
}

func (s *UserService) generateTokenId() string {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	return base64.RawStdEncoding.EncodeToString(buf)
}
