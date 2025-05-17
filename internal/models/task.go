package models

import "time"

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	Summary   string `gorm:"type:varchar(2500)"`
	Date      time.Time
	UserID    uint
	CreatedAt time.Time
}
