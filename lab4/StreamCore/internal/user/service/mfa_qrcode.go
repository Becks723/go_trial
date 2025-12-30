package service

import (
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/env"
	"bytes"
	"encoding/base64"
	"image/png"
	"strconv"

	"github.com/pquerna/otp/totp"
)

func (s *UserService) MFAQrcode(uid uint) (*common.MFAInfo, error) {
	var err error

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

	data := new(common.MFAInfo)
	data.Secret = secret
	data.Qrcode = qrcode
	return data, nil
}
