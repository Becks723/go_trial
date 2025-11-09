package service

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/model/user"
	"StreamCore/biz/service"
	"StreamCore/test/service/mock"
	"context"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	t.Run("username_exists", testRegister_UsernameExists)
}

func testRegister_UsernameExists(t *testing.T) {
	var (
		usernamePh = "username_placeholder"
		passwordPh = "password_placeholder"
		userPh     = &domain.User{
			Username: usernamePh,
			Password: passwordPh,
		}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock repo
	repo := mock.NewMockUserRepo(ctrl)
	repo.EXPECT().
		GetByUsername(usernamePh).
		Return(userPh, nil)

	// action
	serv := service.NewUserService(repo)
	err := serv.Register(context.Background(), &user.RegisterReq{
		Username: usernamePh,
		Password: "123456",
	})
	if err == nil {
		t.Errorf("Register should fail if same username already exists.")
	}
}
