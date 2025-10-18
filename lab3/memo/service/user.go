package service

import (
	"errors"
	"memo/dto"
	"memo/pkg/ctl"
	"memo/pkg/util"
	"memo/repository"
	"memo/repository/model"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	serv := UserService{
		repo: repo,
	}
	return &serv
}

func (serv *UserService) Signup(req *dto.SignupReq) (resp *dto.Response, err error) {
	// 用户名已存在
	u, err := serv.repo.FindUserByName(req.Username)
	if u != nil {
		err = errors.New("用户已存在") // TODO: i18n
		return
	}

	user := model.UserModel{
		Username: req.Username,
		Password: req.Password,
	}
	// 加密密码失败
	if user, err = serv.EncryptPassword(user); err != nil {
		return
	}

	// 写入数据库失败
	if err = serv.repo.InsertUser(&user); err != nil {
		return
	}

	return ctl.ResponseSuccess(), nil
}

func (serv *UserService) Login(req *dto.LoginReq) (resp *dto.Response, err error) {
	record, err := serv.repo.FindUserByName(req.Username)
	if record == nil {
		err = errors.New("用户不存在") // TODO: i18n
		return
	}

	if !util.CheckPassword(req.Password, record.Password) {
		err = errors.New("密码错误") // TODO: i18n
		return
	}

	// 生成token错误
	token, err := util.GenerateToken(record.ID, record.Username)
	if err != nil {
		return
	}

	td := dto.TokenData{Token: token}
	return ctl.ResponseSuccessWithData(td), nil
}

func (serv *UserService) EncryptPassword(in model.UserModel) (out model.UserModel, err error) {
	encrypted, err := util.EncryptPassword(in.Password)
	if err == nil {
		out = model.UserModel{
			Model:    in.Model,
			Username: in.Username,
			Password: encrypted,
		}
	}
	return
}
