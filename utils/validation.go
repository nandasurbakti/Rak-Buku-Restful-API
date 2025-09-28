package utils

import (
	"errors"

	"github.com/nandasurbakti/Rak-Buku-Restful-API/models"
)

func ValidateBook(req models.BookRequest) error {
	if req.ReleaseYear < 1980 || req.ReleaseYear > 2024 {
		return errors.New("tahun rilis harus di antara 1980 dan 2024")
	}
	return nil
}

func CheckThickness(totalPage int) string {
	if totalPage > 100 {
		return "tebal"
	}
	return "tipis"
}