package model

import (
	"github.com/hhr0815hhr/gint/internal/database"
	"gorm.io/gorm"
)

type Test struct {
	Id   int    `gorm:"primarykey"`
	Name string `gorm:"size:50"`
}

type TestRepo struct {
	*database.BaseRepository[Test]
}

func NewTestRepo(db *gorm.DB) *TestRepo {
	return &TestRepo{
		BaseRepository: database.NewBaseRepository[Test](db),
	}
}
