package model

import "time"

type SkillGeneratedArticle struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TaskID          string    `gorm:"size:64;uniqueIndex" json:"taskId"`
	Platform        string    `gorm:"size:32;index" json:"platform"`
	Keyword         string    `gorm:"size:128;index" json:"keyword"`
	Category        string    `gorm:"size:64;index" json:"category"`
	Title           string    `gorm:"size:512" json:"title"`
	TitleOptions    string    `gorm:"type:text" json:"titleOptions"`
	Summary         string    `gorm:"type:text" json:"summary"`
	MarkdownContent string    `gorm:"type:longtext" json:"markdownContent"`
	CoverImageURL   string    `gorm:"size:1024" json:"coverImageUrl"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (SkillGeneratedArticle) TableName() string {
	return "skill_generated_articles"
}
