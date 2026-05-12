package model

import "time"

type CollectedNews struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UniqueKey    string    `gorm:"size:64;uniqueIndex" json:"uniqueKey"`
	Keyword      string    `gorm:"size:128;index" json:"keyword"`
	SearchType   string    `gorm:"size:32;index" json:"searchType"`
	ChannelID    int       `gorm:"index" json:"channelId"`
	ChannelName  string    `gorm:"size:64" json:"channelName"`
	ThirdPartyID string    `gorm:"size:128;index" json:"thirdPartyId"`
	URL          string    `gorm:"size:1024" json:"url"`
	Ctime        string    `gorm:"size:64" json:"ctime"`
	Title        string    `gorm:"size:512" json:"title"`
	PicURL       string    `gorm:"size:1024" json:"picUrl"`
	Source       string    `gorm:"size:128" json:"source"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (CollectedNews) TableName() string {
	return "collected_news"
}

type GeneratedArticle struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Keyword         string    `gorm:"size:128;index" json:"keyword"`
	Title           string    `gorm:"size:512" json:"title"`
	Summary         string    `gorm:"type:text" json:"summary"`
	MarkdownContent string    `gorm:"type:longtext" json:"markdownContent"`
	HTMLPreview     string    `gorm:"type:longtext" json:"htmlPreview"`
	Tags            string    `gorm:"size:512" json:"tags"`
	CoverImageURL   string    `gorm:"size:1024" json:"coverImageUrl"`
	ImageStatus     string    `gorm:"size:32;default:missing" json:"imageStatus"`
	ModelName       string    `gorm:"size:64" json:"modelName"`
	SourceNewsIDs   string    `gorm:"type:text" json:"sourceNewsIds"`
	Status          string    `gorm:"size:32;index;default:draft" json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (GeneratedArticle) TableName() string {
	return "generated_articles"
}

type PublishPackage struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ArticleID       uint      `gorm:"index" json:"articleId"`
	Platform        string    `gorm:"size:32;index" json:"platform"`
	Title           string    `gorm:"size:512" json:"title"`
	MarkdownContent string    `gorm:"type:longtext" json:"markdownContent"`
	CoverImageURL   string    `gorm:"size:1024" json:"coverImageUrl"`
	SourceNewsIDs   string    `gorm:"type:text" json:"sourceNewsIds"`
	Status          string    `gorm:"size:32;index" json:"status"`
	Remark          string    `gorm:"type:text" json:"remark"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (PublishPackage) TableName() string {
	return "publish_packages"
}
