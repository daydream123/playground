package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

func TestGORM(t *testing.T) {
	var (
		err     error
		student User
	)

	if err = DB.Callback().Update().Register("user_update_callback", func(db *gorm.DB) {
		if db.Statement.Schema == nil {
			return
		}

		for _, field := range db.Statement.Schema.Fields {
			if field.Name != "Name" {
				continue
			}

			fmt.Printf("----------------- field[%s] changed? %v ------------------\n", field.Name, db.Statement.Changed(field.Name))
		}
	}); err != nil {
		panic(err)
	}

	if err = DB.Model(&User{}).First(&student).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		student.Name = fmt.Sprintf("name: %d", time.Now().Unix())
		if err = DB.Save(&student).Error; err != nil {
			panic(err)
		}
	}

	newName := fmt.Sprintf("name: %d", time.Now().Unix())
	if err = DB.Model(&User{
		Model: gorm.Model{
			ID: student.ID,
		}}).
		Updates(&User{Name: newName}).
		Error; err != nil {
		panic(err)
	}
}
