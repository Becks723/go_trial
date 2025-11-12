package repo

import (
	"math"
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

func isPageParamsValid(total int64, limit, page int) bool {
	if limit <= 0 {
		return false
	}
	pc := int(math.Ceil(float64(total) / float64(limit))) // page count
	return 0 <= page && page < pc
}
