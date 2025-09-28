package models

import "time"

type Book struct {
	ID          int       `json:"id" db:"id" example:"1"`
	Title       string    `json:"title" db:"title" example:"Judul buku"`
	Description string    `json:"description" db:"description" example:"Deskripsi buku"`
	ImageUrl    string    `json:"image_url" db:"image_url"`
	ReleaseYear int       `json:"release_year" db:"release_year" binding:"required,min=1980,max=2024" example:"2025"`
	Price       int       `json:"price" db:"price" example:"100000"`
	TotalPage   int       `json:"total_page" db:"total_page" example:"100"`
	Thickness   string    `json:"thickness" db:"thickness"`
	CategoryID  int       `json:"category_id" db:"category_id" example:"1"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy 	string 	`json:"created_by" db:"created_by"`
	ModifiedAt 	time.Time `json:"modified_at" db:"modified_at"`
	ModifiedBy 	string `json:"modified_by" db:"modified_by"`
}

type BookRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
    ImageURL    string `json:"image_url"`
    ReleaseYear int    `json:"release_year" binding:"required"`
    Price       int `json:"price" binding:"required"`
    TotalPage   int    `json:"total_page" binding:"required,min=1"`
    CategoryID  uint   `json:"category_id" binding:"required"`
}

type BookWithCategory struct {
	Book
	CategoryName string `json:"category_name" db:"category_name"`
}