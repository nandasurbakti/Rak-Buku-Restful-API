package controllers

import (
	"database/sql"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/database/connection"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/models"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/utils"
)

// @Summary User login
// @Description Login dengan username dan password untuk mendapatkan JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users/login [post]
func Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "format salah", "details": err.Error(),})
        return
    }
    
    var user models.User
    query := "SELECT id, username, password, created_at, created_by, modified_at, modified_by FROM users WHERE username = $1"
    err := connection.DB.Get(&user, query, req.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "usernama atau password tidak valid",})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error",})
        return
    }
    
    if !utils.CheckPasswordHash(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "usernama atau password tidak valid",})
        return
    }
    
    token, err := utils.GenerateJWT(user.ID, user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghasilkan token",})
        return
    }
    
    response := models.LoginResponse{
        Token: token,
        User:  user,
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "login sukses","data": response,})
}

// @Summary User registration
// @Description daftar akun baru
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /api/users/register [post]
func Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "format tidak valid", "details": err.Error(),})
        return
    }
    
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
    err := connection.DB.Get(&exists, checkQuery, req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error",})
        return
    }
    if exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "username atau email sudah ada",})
        return
    }
    
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat hash kata sandi",})
        return
    }
    
    var user models.User
	now := time.Now()
	createdBy := req.Username
	modifiedBy := createdBy

    query := `INSERT INTO users (username, password, created_at, created_by, modified_at, modified_by) 
          VALUES ($1, $2, $3, $4, $5, $6) 
          RETURNING id, username, created_at, created_by, modified_at, modified_by`
    err = connection.DB.Get(&user, query, req.Username, hashedPassword, now, createdBy, now, modifiedBy)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat user",})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"message": "berhasil mendaftar",
									 "data": user,
    })
}