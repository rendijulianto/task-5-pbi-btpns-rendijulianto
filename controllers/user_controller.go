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

// CreateUser membuat pengguna baru.
func CreateUser(context *gin.Context) {
	var userFormRegister app.UserFormRegister
	if err := context.ShouldBindJSON(&userFormRegister); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	  // Validasi data pengguna menggunakan govalidator
	  if _, err := govalidator.ValidateStruct(userFormRegister); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	    }
	// cek apakah email sudah terdaftar
	var user models.User

	if len(userFormRegister.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("email = ?", userFormRegister.Email).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		context.Abort()
		return
	}

	// cek apakah username sudah terdaftar	
	if err := database.Instance.Where("username = ?", userFormRegister.Username).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar"})
		context.Abort()
		return
	}


	user = models.User{
		Username: userFormRegister.Username,
		Email:    userFormRegister.Email,
		Password: userFormRegister.Password,
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message" : "Berhasil Membuat Akun"})
    }

    // Login mengautentikasi pengguna.
    func Login(context *gin.Context) {
	var userFormLogin app.UserFormLogin
	if err := context.ShouldBindJSON(&userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	  // Validasi data pengguna menggunakan govalidator
	  if _, err := govalidator.ValidateStruct(userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	  }


	var user models.User
	if err := database.Instance.Where("email = ?", userFormLogin.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := user.CheckPassword(userFormLogin.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email atau password salah"})
		return
	}



	token, err := helpers.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "token": token})
	}

    
    // GetUserByID mengambil pengguna berdasarkan ID.
    func GetUserByID(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
	    context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
	    return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diizinkan"})
		context.Abort()
		return
	}
	var user models.User
	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

    
	
	if err := database.Instance.First(&user, userID).Error; err != nil {
	    context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
	    return
	}

	var userResult app.UserResult
	userResult.ID = user.ID
	userResult.Username = user.Username
	userResult.Email = user.Email
	userResult.CreatedAt = user.CreatedAt.String()
	userResult.UpdatedAt = user.UpdatedAt.String()
    
	context.JSON(http.StatusOK, gin.H{"data": userResult})
    }
    
    // UpdateUser mengupdate pengguna berdasarkan ID.
    func UpdateUser(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
	    context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
	    return
	}
	var userFormUpdate app.UserFormUpdate
	if err := context.ShouldBindJSON(&userFormUpdate); err != nil {
	    context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	    return
	}

	  // Validasi data pengguna menggunakan govalidator
	  if _, err := govalidator.ValidateStruct(userFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	  }
	  var user models.User 
	//   password minimal 6 karakter
	if len(userFormUpdate.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
		context.Abort()
		return
	}

	//   validasi email dan username sudah terdaftar atau belum selain data user yang sedang login
	if err := database.Instance.Where("email = ? AND id != ?", userFormUpdate.Email, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("username = ? AND id != ?", userFormUpdate.Username, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar"})
		context.Abort()
		return
	}

		
	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diizinkan"})
		context.Abort()
		return
	}


	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}


	if err := database.Instance.First(&user, userID).Error; err != nil {
	    context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
	    return
	}
    
	user.Username = userFormUpdate.Username
	user.Email = userFormUpdate.Email
	if userFormUpdate.Password != "" {
		if err := user.HashPassword(userFormUpdate.Password); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
	}


    
	if err := database.Instance.Save(&user).Error; err != nil {
	    context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	    return
	}
    
	context.JSON(http.StatusOK, gin.H{"message": "Berhasil mengupdate pengguna"})
    }
    
    // DeleteUser menghapus pengguna berdasarkan ID.
    func DeleteUser(context *gin.Context) {
	var user models.User
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
	    context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
	    return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if (userID != int(claims.ID)) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diizinkan"})
		context.Abort()
		return
	}
	
	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}
	

	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := database.Instance.Delete(&user).Error; err != nil {
	 context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	 return
	}
	
	
	context.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus pengguna"})
    }