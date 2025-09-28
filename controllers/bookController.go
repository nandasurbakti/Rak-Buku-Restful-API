package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/database/connection"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/models"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/utils"
)

// @Summary Get all books
// @Description Get list of all books with category information
// @Tags books
// @Produce json
// @Success 200 {array} models.BookWithCategory
// @Failure 400
// @Router /api/books [get]
func GetBooks(c *gin.Context) {
	var books []models.BookWithCategory

	query := `SELECT b.id, b.title, b.description, b.image_url, b.release_year, b.price,
            b.total_page, b.thickness, b.category_id, b.created_at, b.created_by, b.modified_at, b.modified_by,
            c.name as category_name
            FROM books b
            JOIN categories c ON b.category_id = c.id
            ORDER BY b.id`

	err := connection.DB.Select(&books, query)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "gagal mengambil data buku",
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// @Summary Get book by ID
// @Description Get book details by ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.BookWithCategory
// @Failure 400
// @Router /api/books/:id [get]
func GetBookByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID",})
        return
    }
    
    var book models.BookWithCategory
    query := `SELECT b.id, b.title, b.description, b.image_url, b.release_year, b.price,
            b.total_page, b.thickness, b.category_id, b.created_at, b.created_by, b.modified_at, b.modified_by,
            c.name as category_name
              FROM books b
              JOIN categories c ON b.category_id = c.id
              WHERE b.id = $1`
    
    err = connection.DB.Get(&book, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "buku tidak ditemukan",})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data buku",})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"data": book,})
}

// @Summary Create new book
// @Description Create a new book with automatic thickness calculation
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.BookRequest true "Book data"
// @Success 201 {object} models.BookWithCategory
// @Failure 400
// @Router /api/books [post]
func PostBook(c *gin.Context) {
    var req models.BookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
        return
    }
    
    var categoryExists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)"
    err := connection.DB.Get(&categoryExists, checkQuery, req.CategoryID)
    if err != nil || !categoryExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "kategori tidak ditemukan",})
        return
    }

    if err := utils.ValidateBook(req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
    }
    thickness := utils.CheckThickness(req.TotalPage)
    now := time.Now()
    username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan",})
		return
	}

    var book models.BookWithCategory
    query := `INSERT INTO books (title, description, image_url, release_year, price, 
              total_page, thickness, category_id, created_at, created_by, modified_at, modified_by) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
              RETURNING id, title, description, image_url, release_year, price, total_page, thickness,
              category_id, created_at, created_by, modified_at, modified_by`

	createdBy := username.(string)
	modifiedBy := username.(string)
    err = connection.DB.Get(&book, query, req.Title, req.Description, req.ImageURL, 
        req.ReleaseYear, req.Price, req.TotalPage, thickness, req.CategoryID, 
        now, createdBy, now, modifiedBy)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat data buku",})
        return
    }
    
    categoryQuery := "SELECT name FROM categories WHERE id = $1"
    err = connection.DB.Get(&book.CategoryName, categoryQuery, req.CategoryID)
    if err != nil {
        book.CategoryName = ""
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "berhasil dibuat","data": book,})
}

// @Summary Update book
// @Description Update book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.BookRequest true "Book data"
// @Success 200 {object} models.BookWithCategory
// @Failure 400
// @Router /api/books/:id [put]
func UpdateBook(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID buku tidak valid",})
        return
    }
    
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)"
    err = connection.DB.Get(&exists, checkQuery, id)
    if err != nil || !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "buku tidak ditemukan",})
        return
    }
    
    var req models.BookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
        return
    }

    var categoryExists bool
    checkCategoryQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)"
    err = connection.DB.Get(&categoryExists, checkCategoryQuery, req.CategoryID)
    if err != nil || !categoryExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "kategori tidak ditemukan",})
        return
    }

    if err := utils.ValidateBook(req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
    }
    thickness := utils.CheckThickness(req.TotalPage)
    now := time.Now()
    username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan",})
		return
	}
	createdBy := username.(string)
	modifiedBy := username.(string)
    
    var book models.BookWithCategory
    query := `UPDATE books SET title = $1, description = $2, image_url = $3, release_year = $4, 
              price = $5, total_page = $6, thickness = $7, category_id = $8, created_at = $9, 
			  created_by = $10, modified_at = $11, modified_by = $12
              WHERE id = $13
              RETURNING id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by`
    
    err = connection.DB.Get(&book, query, req.Title, req.Description, req.ImageURL, 
		req.ReleaseYear, req.Price, req.TotalPage, thickness, req.CategoryID, now, createdBy, now, modifiedBy, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui buku",})
        return
    }
    
    categoryQuery := "SELECT name FROM categories WHERE id = $1"
    err = connection.DB.Get(&book.CategoryName, categoryQuery, req.CategoryID)
    if err != nil {
        book.CategoryName = ""
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "buku berhasil diperbarui",
								"data": book,})
}

// @Summary Delete book
// @Description Delete book by ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 400
// @Router /api/books/:id [delete]
func DeleteBook(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID buku tidak valid",})
        return
    }
    
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)"
    err = connection.DB.Get(&exists, checkQuery, id)
    if err != nil || !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "buku tidak ditemukan",})
        return
    }
    
    query := "DELETE FROM books WHERE id = $1"
    _, err = connection.DB.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus buku",})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "buku berhasil dihapus",}) 
}