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
