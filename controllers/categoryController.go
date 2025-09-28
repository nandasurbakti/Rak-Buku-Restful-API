package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/database/connection"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/models"
)

// @Summary Get all categories
// @Description Get list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 400
// @Router /api/categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category

	query := "SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories ORDER BY id"
	err := connection.DB.Select(&categories, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal ambil data",})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": categories,})
}

// @Summary Get category by ID
// @Description Get category details by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400
// @Router /api/categories/:id [get]
func GetCategoryByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid",})
        return
    }
    
    var category models.Category
    query := "SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories WHERE id = $1"
    err = connection.DB.Get(&category, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan",})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data kategori",})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"data": category,})
}

// @Summary Get books by category
// @Description Get all books in a specific category
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {array} models.BookWithCategory
// @Failure 400
// @Router /api/categories/:id/books [get]
func GetBooksByCategory(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid",})
        return
    }
    
    var books []models.BookWithCategory
    query := `SELECT b.id, b.title, b.description, b.image_url, b.release_year, b.price, 
              b.total_page, b.thickness, b.category_id, b.created_at, b.created_by, b.modified_at, b.modified_by, 
			  c.name as category_name
              FROM books b
              JOIN categories c ON b.category_id = c.id
              WHERE b.category_id = $1
              ORDER BY b.id`
    
    err = connection.DB.Select(&books, query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data buku",})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"data": books,})
}

// @Summary Create new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CategoryRequest true "Category data"
// @Success 201 {object} models.Category
// @Failure 400
// @Security BearerAuth
// @Router /api/categories [post]
func PostCategory(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	var category models.Category
	query := `INSERT INTO categories (name, created_at, created_by, modified_at, modified_by)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, name, created_at, created_by, modified_at, modified_by`

	now := time.Now()
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan",})
		return
	}
	createdBy := username.(string)
	modifiedBy := username.(string)
	err := connection.DB.Get(&category, query, req.Name, now, createdBy, now, modifiedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat kategori",})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "berhasil membuat kategori",
									 "data": category,})
}

// @Summary Update category
// @Description Update category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.CategoryRequest true "Category data"
// @Success 200 {object} models.Category
// @Failure 400
// @Security BearerAuth
// @Router /api/categories/:id [put]
func UpdateCategory(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid",})
        return
    }
    
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)"
    err = connection.DB.Get(&exists, checkQuery, id)
    if err != nil || !exists {c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan",})
        return
    }
    
    var req models.CategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var old models.Category
    err = connection.DB.Get(&old, "SELECT created_at, created_by FROM categories WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal ambil data lama"})
        return
    }

    now := time.Now()
    username, ok := c.Get("username")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan"})
        return
    }
    modifiedBy := username.(string)

    var category models.Category
    query := `UPDATE categories 
              SET name = $1, created_at = $2, created_by = $3, modified_at = $4, modified_by = $5
              WHERE id = $6 
              RETURNING id, name, created_at, created_by, modified_at, modified_by`

    err = connection.DB.Get(&category, query, req.Name, old.CreatedAt, old.CreatedBy, now, modifiedBy, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui kategori"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil diperbarui",
        						"data": category,})
}

// @Summary Delete category
// @Description Delete category by ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400
// @Security BearerAuth
// @Router /api/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid",})
        return
    }
    
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)"
    err = connection.DB.Get(&exists, checkQuery, id)
    if err != nil || !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan",})
        return
    }
    
    query := "DELETE FROM categories WHERE id = $1"
    _, err = connection.DB.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus kategori",})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "kategori berhasil dihapus",
    })
}