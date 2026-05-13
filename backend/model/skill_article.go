package model

import "time"

type SkillGeneratedArticle struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TaskID          string    `gorm:"size:64;uniqueIndex" json:"taskId"`
	Platform        string    `gorm:"size:32;index" json:"platform"`
	Keyword         string    `gorm:"size:128;index" json:"keyword"`
	KeywordSegments string    `gorm:"type:text" json:"keywordSegments"`
	Category        string    `gorm:"size:64;index" json:"category"`
	Title           string    `gorm:"size:512" json:"title"`
	TitleOptions    string    `gorm:"type:text" json:"titleOptions"`
	Summary         string    `gorm:"type:text" json:"summary"`
	MarkdownContent string    `gorm:"type:longtext" json:"markdownContent"`
	HTMLPreview     string    `gorm:"type:longtext" json:"htmlPreview"`
	Tags            string    `gorm:"type:text" json:"tags"`
	CoverImageURL   string    `gorm:"size:1024" json:"coverImageUrl"`
	CoverImageType  string    `gorm:"size:32" json:"coverImageType"`
	Images          string    `gorm:"type:longtext" json:"images"`
	Sources         string    `gorm:"type:longtext" json:"sources"`
	HotTopics       string    `gorm:"type:longtext" json:"hotTopics"`
	StyleProfile    string    `gorm:"type:text" json:"styleProfile"`
	WordCount       int       `json:"wordCount"`
	ModelProvider   string    `gorm:"size:64" json:"modelProvider"`
	ModelName       string    `gorm:"size:128" json:"modelName"`
	PromptVersion   string    `gorm:"size:64" json:"promptVersion"`
	SkillVersion    string    `gorm:"size:64" json:"skillVersion"`
	HumanizeStatus  string    `gorm:"size:32;default:pending" json:"humanizeStatus"`
	Status          string    `gorm:"size:32;index;default:draft" json:"status"`
	PublishStatus   string    `gorm:"size:32;index;default:unpublished" json:"publishStatus"`
	PublishPayload  string    `gorm:"type:longtext" json:"publishPayload"`
	ErrorMessage    string    `gorm:"type:text" json:"errorMessage"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (SkillGeneratedArticle) TableName() string {
	return "skill_generated_articles"
}
