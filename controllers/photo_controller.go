package controllers

import (
	"final-project-rakamin/app"
	"final-project-rakamin/database"
	"final-project-rakamin/helpers"
	"final-project-rakamin/models"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)


func GetAllPhotos(context *gin.Context) {
	var photos []app.PhotoResult
	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// munculkan foto berdasarkan user_id dan diurutkan berdasarkan created_at  db.Preload("UserResult").
	if err := database.Instance.Table("photos").Select("photos.id, photos.title, photos.caption, photos.photo_url, photos.created_at, photos.updated_at, users.email").Joins("JOIN users ON users.id = photos.user_id").Where("photos.user_id = ?", claims.ID).Order("photos.created_at desc").Scan(&photos).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": photos})
}

func GetPhotoByID(context *gin.Context) {
	var photo app.PhotoResult
	id := context.Param("id")
	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := database.Instance.Table("photos").Select("photos.id, photos.title, photos.caption, photos.photo_url, photos.created_at, photos.updated_at, users.email").Joins("JOIN users ON users.id = photos.user_id").Where("photos.id = ? AND photos.user_id = ?", id, claims.ID).Scan(&photo).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	
	if photo.ID == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "Foto tidak ditemukan"})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": photo})

}

func CreatePhoto(context *gin.Context) {
	var photoFormCreate app.PhotoFormCreate
	if err := context.ShouldBindJSON(&photoFormCreate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi data pengguna menggunakan govalidator
	if _, err := govalidator.ValidateStruct(photoFormCreate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	desiredExtensions := []string{"jpg", "jpeg", "png", "gif"}

	// Periksa apakah URL foto valid dan berakhir dengan salah satu ekstensi yang diinginkan
	if !helpers.IsValidURLWithDesiredExtension(photoFormCreate.PhotoUrl, desiredExtensions) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "URL foto tidak valid atau tidak berakhir dengan ekstensi yang diinginkan"})
		return
	}



	// Ambil user_id dari token
	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Buat objek photo baru
	photo := models.Photo{
		Title:    photoFormCreate.Title,
		Caption:  photoFormCreate.Caption,
		PhotoUrl: photoFormCreate.PhotoUrl,
		UserID:   claims.ID,
	}

	// Simpan objek photo ke database
	if err := database.Instance.Create(&photo).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Foto berhasil ditambahkan"})
}


func UpdatePhoto(context *gin.Context) {
	photoID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID foto tidak valid"})
		return
	}
	var photoFormUpdate app.PhotoFormUpdate
	if err := context.ShouldBindJSON(&photoFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi data pengguna menggunakan govalidator
	if _, err := govalidator.ValidateStruct(photoFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	desiredExtensions := []string{"jpg", "jpeg", "png", "gif"}

	// Periksa apakah URL foto valid dan berakhir dengan salah satu ekstensi yang diinginkan
	if !helpers.IsValidURLWithDesiredExtension(photoFormUpdate.PhotoUrl, desiredExtensions) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "URL foto tidak valid atau tidak berakhir dengan ekstensi yang diinginkan"})
		return
	}

	// Ambil user_id dari token
	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Ambil data photo dari database
	var photo models.Photo

	if err := database.Instance.Where("id = ? AND user_id = ?", photoID, claims.ID).First(&photo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Foto tidak ditemukan"})
		return
	}



	photo.Title = photoFormUpdate.Title
	photo.Caption = photoFormUpdate.Caption
	photo.PhotoUrl = photoFormUpdate.PhotoUrl
	if err := database.Instance.Save(&photo).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Foto berhasil diupdate"})

}

func DeletePhoto(context *gin.Context) {
	var photo models.Photo
	photoID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID foto tidak valid"})
		return
	}

	// Ambil user_id dari token
	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := database.Instance.First(&photo, photoID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Foto tidak ditemukan"})
		return
	}

	if photo.UserID != claims.ID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Anda tidak memiliki akses untuk menghapus foto ini"})
		return
	}

	if err := database.Instance.Delete(&photo).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Foto berhasil dihapus"})
}

