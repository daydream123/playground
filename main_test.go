package main

import (
	"errors"
	"gorm.io/gorm"
	"testing"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

func TestGORM(t *testing.T) {
	if err := DB.Callback().Create().Before("gorm:create").Register("user_create_callback", func(db *gorm.DB) {
		if db.Statement.Schema == nil {
			return
		}

		var language Language
		if err := db.
			Model(&Language{}).
			Where(&Language{Code: "1111"}).
			Take(&language).
			Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				// no such column: companies.code,
				// 即便我指定了查询 `Language`
				// 我不能指定 db.Session(&gorm.Session{NewDB: true})，因为sqlite不能同时存在2个connect(准确讲是2个writeable conn)
				panic(err)
			}
		}
	}); err != nil {
		panic(err)
	}

	if err := DB.Save(&Company{
		Name: "hello world",
	}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}
}
