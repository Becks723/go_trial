package repository

import (
	"math"
	"memo/repository/model"
)

type MemoRepository struct {
	baseRepository
}

func NewMemoRepo() *MemoRepository {
	return &MemoRepository{
		baseRepository{
			db: db,
		},
	}
}

func (repo *MemoRepository) InsertMemo(memo *model.MemoModel) error {
	return repo.db.
		Model(&model.MemoModel{}).
		Create(memo).
		Error
}

func (repo *MemoRepository) UpdateMemo(newMemo *model.MemoModel) error {
	// 找到数据库中的值
	var record model.MemoModel
	err := repo.db.
		Model(&model.MemoModel{}).
		Where("id = ?", newMemo.Id).
		First(&record).
		Error
	if err != nil {
		return err
	}

	// 更新字段
	if newMemo.Title != "" {
		record.Title = newMemo.Title
	}
	if newMemo.Content != "" {
		record.Content = newMemo.Content
	}
	if newMemo.Status != 0 {
		record.Status = newMemo.Status
	}
	if newMemo.StartsAt != nil {
		record.StartsAt = newMemo.StartsAt
	}
	if newMemo.EndsAt != nil {
		record.EndsAt = newMemo.EndsAt
	}

	// 保存
	return repo.db.
		Save(&record).
		Error
}

func (repo *MemoRepository) FindMemoById(id uint) (*model.MemoModel, error) {
	var match model.MemoModel
	err := repo.db.
		Model(&model.MemoModel{}).
		Where("id = ?", id).
		First(&match).
		Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

/* 根据备忘录的id查找它属于哪个用户*/
func (repo *MemoRepository) FindCreatorId(id uint) (uid uint, err error) {
	memo, err := repo.FindMemoById(id)
	if err != nil {
		return
	}
	uid = memo.Uid
	return
}

func (repo *MemoRepository) CountMemos(uid uint) (count int64, err error) {
	err = repo.db.
		Model(&model.MemoModel{}).
		Where("uid = ?", uid).
		Count(&count).
		Error
	return
}

func (repo *MemoRepository) FindAllMemos(uid uint, limit, ps, pe int) (records []*model.MemoModel, err error) {
	records, err = repo.findMemosCore(uid, limit, ps, pe, "")
	return
}

func (repo *MemoRepository) FindPendingMemos(uid uint, limit, ps, pe int) (records []*model.MemoModel, err error) {
	records, err = repo.findMemosCore(uid, limit, ps, pe, "status = ?", model.MemoStatusPending)
	return
}

func (repo *MemoRepository) FindCompletedMemos(uid uint, limit, ps, pe int) (records []*model.MemoModel, err error) {
	records, err = repo.findMemosCore(uid, limit, ps, pe, "status = ?", model.MemoStatusCompleted)
	return
}

func (repo *MemoRepository) findMemosCore(uid uint, limit, ps, pe int, query string, args ...any) (records []*model.MemoModel, err error) {
	tx := repo.db.
		Model(&model.MemoModel{}).
		Where("uid = ?", uid)

	// 额外的条件
	if query != "" {
		tx = tx.Where(query, args)
	}

	// 分页查询
	total, err := repo.CountMemos(uid)
	if err != nil {
		return
	}
	pages := int64(math.Ceil(float64(total) / float64(limit))) // 总页数
	if 1 <= ps && int64(ps) <= pages &&
		1 <= pe && int64(pe) <= pages &&
		ps <= pe &&
		limit > 0 {
		// 仅分页参数合法时分页，否则查询全部
		tx = tx.Limit((pe - ps + 1) * limit).
			Offset((ps - 1) * limit)
	}

	err = tx.Find(&records).Error
	return
}
