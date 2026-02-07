package service

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strconv"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/common"
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

	w := constants.MFA_QrcodeWidth
	h := constants.MFA_QrcodeHeight
	img, err := key.Image(w, h)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil { // encode img to png (binary)
		return nil, err
	}
	qrcode := base64.StdEncoding.EncodeToString(buf.Bytes()) // base64

	// cache secret
	err = s.cache.SetTOTPPending(s.ctx, uid, secret)
	if err != nil {
		return nil, err
	}

	data := new(common.MFAInfo)
	data.Secret = secret
	data.Qrcode = qrcode
	return data, nil
}
