package router

import (
	"final-project-rakamin/controllers"
	"final-project-rakamin/middlewares"

	"github.com/gin-gonic/gin"
)
	    


func SetupRouter() *gin.Engine {
	r := gin.Default()
    
	// Public routes (without authentication middleware)
	public := r.Group("/api")
	{
		public.GET("/users/login", controllers.Login)
	    	public.POST("/users/register", controllers.CreateUser)
	}
    
	// Protected routes (with authentication middleware)
	protected := r.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
	    // Routes for the resource: user
	    protected.GET("/users/:id", controllers.GetUserByID)
	    protected.PUT("/users/:id", controllers.UpdateUser)
	    protected.DELETE("/users/:id", controllers.DeleteUser)
	    // Routes for the resource: photo
	    protected.GET("/photos", controllers.GetAllPhotos)
	    protected.GET("/photos/:id", controllers.GetPhotoByID)
	    protected.POST("/photos", controllers.CreatePhoto)
	    protected.PUT("/photos/:id", controllers.UpdatePhoto)
	    protected.DELETE("/photos/:id", controllers.DeletePhoto)
	}
    
	return r
    }