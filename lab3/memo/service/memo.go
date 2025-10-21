package service

import (
	"errors"
	"memo/config"
	"memo/dto"
	"memo/pkg/ctl"
	"memo/repository"
	"memo/repository/model"
	"strings"
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

func (serv *MemoService) Update(uid uint, req *dto.UpdateMemoReq) (resp *dto.Response, err error) {
	// 找到这条备忘录的创建者，并与传入的uid进行比对
	whose, err := serv.repo.FindCreatorId(req.Id)
	if err != nil {
		return
	}

	// 每个用户只能修改自己创建的备忘录
	if whose != uid {
		err = errors.New("Unknown memo id.") // TODO: i18n
		return
	}

	// 更新数据库
	err = serv.repo.UpdateMemo(&model.MemoModel{
		Id:       req.Id,
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status,
		StartsAt: req.StartsAt,
		EndsAt:   req.EndsAt,
	})
	if err != nil {
		return
	}

	return ctl.ResponseSuccess(), nil
}

func (serv *MemoService) List(uid uint, params *dto.ListMemoParams) (resp *dto.Response, err error) {
	if params.Limit <= 0 {
		params.Limit = config.DefaultLimit
	}

	// 按过滤条件分页查询
	var records []*model.MemoModel
	switch params.Filter {
	case dto.ListFilterNone:
		records, err = serv.repo.FindAllMemos(uid, params.Limit, params.PageStart, params.PageEnd)
	case dto.ListFilterPending:
		records, err = serv.repo.FindPendingMemos(uid, params.Limit, params.PageStart, params.PageEnd)
	case dto.ListFilterCompleted:
		records, err = serv.repo.FindCompletedMemos(uid, params.Limit, params.PageStart, params.PageEnd)
	}
	if err != nil {
		return
	}

	// 作为resp中的data返回
	var memos []dto.MemoData
	for _, record := range records {
		memos = append(memos, dto.MemoData{
			Id:       record.Id,
			Title:    record.Title,
			Content:  record.Content,
			Status:   record.Status,
			StartsAt: record.StartsAt,
			EndsAt:   record.EndsAt,
		})
	}
	return ctl.ResponseSuccessWithData(memos), nil
}

func (serv *MemoService) Search(uid uint, params *dto.SearchMemoParams) (resp *dto.Response, err error) {
	// 处理关键词
	keywords := serv.normalizeKeywords(params.Keywords)

	// 数据库查询
	records, err := serv.repo.SearchMemos(uid, keywords, params.Limit, params.PageStart, params.PageEnd)
	if err != nil {
		return
	}

	// 作为resp中的data返回
	var memos []dto.MemoData
	for _, record := range records {
		memos = append(memos, dto.MemoData{
			Id:       record.Id,
			Title:    record.Title,
			Content:  record.Content,
			Status:   record.Status,
			StartsAt: record.StartsAt,
			EndsAt:   record.EndsAt,
		})
	}
	return ctl.ResponseSuccessWithData(memos), nil
}

func (serv *MemoService) DeleteById(uid uint, req *dto.DeleteMemoByIdReq) (resp *dto.Response, err error) {
	// 找到这条备忘录的创建者，并与传入的uid进行比对
	whose, err := serv.repo.FindCreatorId(req.Id)
	if err != nil {
		return
	}

	// 每个用户只能修改自己创建的备忘录
	if whose != uid {
		err = errors.New("Unknown memo id.") // TODO: i18n
		return
	}

	// 从数据库删除
	if err = serv.repo.DeleteMemoById(req.Id); err != nil {
		return
	}

	return ctl.ResponseSuccess(), nil
}

func (serv *MemoService) DeleteByFilter(uid uint, req *dto.DeleteMemoByFilterReq) (resp *dto.Response, err error) {
	// 从数据库删除
	switch req.Filter {
	case dto.DeleteFilterNone:
		err = serv.repo.DeleteAllMemos(uid)
	case dto.DeleteFilterPending:
		err = serv.repo.DeletePendingMemos(uid)
	case dto.DeleteFilterCompleted:
		err = serv.repo.DeleteCompletedMemos(uid)
	}
	if err != nil {
		return
	}

	return ctl.ResponseSuccess(), nil
}

func (serv *MemoService) normalizeKeywords(value string) string {
	result := strings.TrimSpace(value)
	return result
}
