package repository

import "memo/repository/model"

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
