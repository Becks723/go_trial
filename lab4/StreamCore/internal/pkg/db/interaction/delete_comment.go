package interaction

import "StreamCore/biz/repo/model"

func (repo *iactiondb) DeleteCommentById(cid, authorId uint) (err error) {
	err = repo.db.
		Where("id = ? AND author_id = ?", cid, authorId).
		Delete(&model.CommentModel{}).
		Error
	if err != nil {
		return
	}
	// delete all subs
	err = repo.db.
		Where("parent_id = ?", cid).
		Delete(&model.CommentModel{}).
		Error
	return
}
