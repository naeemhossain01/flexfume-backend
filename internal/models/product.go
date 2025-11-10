package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductName string         `gorm:"not null" json:"productName"`
	ProductCode string         `gorm:"uniqueIndex" json:"productCode"`
	Description string         `gorm:"type:text" json:"description"`
	LongDescription string    `gorm:"type:text" json:"longDescription"`
	ImageURL        string         `json:"imageUrl"`
	YouTubeVideoUrl string         `gorm:"column:youtube_video_url" json:"youtubeVideoUrl"`
	Price           float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	CategoryID         string           `gorm:"type:uuid;not null" json:"categoryId"`
	KeyFeatures        pq.StringArray   `gorm:"type:text[]" json:"keyFeatures,omitempty"`
	WhyChooseBenefits  pq.StringArray   `gorm:"type:text[];column:why_choose_benefits" json:"whyChooseBenefits,omitempty"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy          string         `json:"-"`
	UpdatedBy          string         `json:"-"`
	
	// Relationships
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Discount *Discount `gorm:"foreignKey:ProductID" json:"discount,omitempty"`
}

// TableName specifies the table name for the Product model
func (Product) TableName() string {
	return "products"
}

// ProductResponse represents the product data returned in API responses
type ProductResponse struct {
	ID           string            `json:"id"`
	ProductName  string            `json:"productName"`
	ProductCode  string            `json:"productCode"`
	Description  string            `json:"description"`
	LongDescription string         `json:"longDescription"`
	ImageURL        string            `json:"imageUrl"`
	YouTubeVideoUrl string            `json:"youtubeVideoUrl"`
	Price           float64           `json:"price"`
	Stock        int               `json:"stock"`
	CategoryID         string            `json:"categoryId"`
	KeyFeatures        []string          `json:"keyFeatures,omitempty"`
	WhyChooseBenefits  []string          `json:"whyChooseBenefits,omitempty"`
	CategoryInfo       *CategoryResponse `json:"categoryInfo,omitempty"`
	DiscountPercentage *int              `json:"discountPercentage,omitempty"`
	DiscountPrice      *float64          `json:"discountPrice,omitempty"`
	CreatedAt          time.Time         `json:"createdAt"`
	UpdatedAt          time.Time         `json:"updatedAt"`
}

// ToResponse converts a Product model to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	response := ProductResponse{
		ID:                p.ID,
		ProductName:       p.ProductName,
		ProductCode:       p.ProductCode,
		Description:       p.Description,
		LongDescription:   p.LongDescription,
		ImageURL:          p.ImageURL,
		YouTubeVideoUrl:   p.YouTubeVideoUrl,
		Price:             p.Price,
		Stock:             p.Stock,
		CategoryID:        p.CategoryID,
		KeyFeatures:       p.KeyFeatures,
		WhyChooseBenefits: p.WhyChooseBenefits,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
	
	if p.Category != nil {
		categoryResp := p.Category.ToResponse()
		response.CategoryInfo = &categoryResp
	}
	
	if p.Discount != nil {
		response.DiscountPercentage = &p.Discount.Percentage
		response.DiscountPrice = &p.Discount.DiscountPrice
	}
	
	return response
}
