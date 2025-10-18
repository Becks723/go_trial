package service

import (
	"memo/dto"
	"memo/repository"
)

type MemoService struct {
	repo *repository.MemoRepository
}

func NewMemoService(repo *repository.MemoRepository) *MemoService {
	serv := MemoService{
		repo: repo,
	}
	return &serv
}

func (serv *MemoService) Add(req *dto.AddMemoReq) (resp *dto.Response, err error) {
	return nil, nil
}
