package store

import (
	"alc/model/auth"
	"time"
)

type Image struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
}

type Type string

const (
	MangaType Type = "MANGA"
)

type Category struct {
	Id          int    `json:"id"`
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Img         Image  `json:"img"`
	Slug        string `json:"slug"`
}

type Item struct {
	Id              int       `json:"id"`
	Category        Category  `json:"category"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	LongDescription string    `json:"longDescription"`
	Img             Image     `json:"img"`
	LargeImg        Image     `json:"largeImg"`
	Slug            string    `json:"slug"`
	CreatedBy       auth.User `json:"-"`
}

type Product struct {
	Id      int               `json:"id"`
	Item    Item              `json:"item"`
	Name    string            `json:"name"`
	Price   int               `json:"price"` // Stored in USD cents
	Stock   *int              `json:"stock"`
	Details map[string]string `json:"details"`
	Slug    string            `json:"slug"`
}

type ProductDiscount struct {
	Id            int
	Product       Product
	DiscountValue int
	ValidFrom     time.Time
	ValidUntil    time.Time
	CouponCode    *string
	MinimumAmount *int
	MaximumAmount *int
}

// Comment management
type ItemComment struct {
	Id          int
	Item        Item
	CommentedBy auth.User
	Title       string
	Message     string
	Rating      int
	UpVotes     int
	DownVotes   int
	IsEdited    bool
	CreatedAt   time.Time
	EditedAt    time.Time
}

func (t Type) ToSlug() string {
	return "manga"
}
