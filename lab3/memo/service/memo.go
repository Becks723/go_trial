package service

import (
	"memo/dto"
	"memo/pkg/ctl"
	"memo/repository"
	"memo/repository/model"
)

type MemoService struct {
	repo     *repository.MemoRepository
	userRepo *repository.UserRepository
}

func NewMemoService(repo *repository.MemoRepository, userRepo *repository.UserRepository) *MemoService {
	serv := MemoService{
		repo:     repo,
		userRepo: userRepo,
	}
	return &serv
}

func (serv *MemoService) Add(uid uint, req *dto.AddMemoReq) (resp *dto.Response, err error) {
	// 通过uid找到user
	user, err := serv.userRepo.FindUserById(uid)
	if err != nil {
		return
	}

	// 添加一条代办事项。要求用户输入标题、内容、开始时间、截至时间
	memo := model.MemoModel{
		Title:    req.Title,
		Content:  req.Content,
		Status:   model.MemoStatusPending,
		StartsAt: req.StartsAt,
		EndsAt:   req.EndsAt,
		Uid:      uid,
		User:     *user,
	}
	// 插入数据库
	if err = serv.repo.InsertMemo(&memo); err != nil {
		return
	}

	return ctl.ResponseSuccess(), nil
}
