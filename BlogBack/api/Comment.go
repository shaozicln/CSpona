package model

import "time"

type Comment struct {
	Id        uint      `gorm:"primary key; autoIncrement" column:"id"`
	ArticleId int       `gorm:"column:article_id"`
	UserId    int       `gorm:"column:user_id"`
	Idea      string    `gorm:"column:idea"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	User      User
}
