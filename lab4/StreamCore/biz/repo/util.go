package repo

import (
	"time"

	"gorm.io/gorm"
)

func deletedAtToPtr(t gorm.DeletedAt) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func ptrToDeletedAt(t *time.Time) gorm.DeletedAt {
	var da gorm.DeletedAt
	if t != nil {
		da.Time = *t
		da.Valid = true
	} else {
		da.Valid = false
	}
	return da
}
