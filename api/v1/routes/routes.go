package routes

import (
	"personal-api/api/v1/handler"
	"personal-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	personalInfoHandler *handler.PersonalInfoHandler,
	contactInfoHandler *handler.ContactInfoHandler,
	socialLinkHandler *handler.SocialLinkHandler,
	postImageHandler *handler.PostImageHandler,
) {
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public info routes
		v1.GET("/personal", personalInfoHandler.GetPersonalInfo)
		v1.GET("/contact", contactInfoHandler.GetContactInfo)
		v1.GET("/social", socialLinkHandler.GetAllSocialLinks)
		v1.GET("/social/:id", socialLinkHandler.GetSocialLink)

		// Post routes - some public, some protected
		posts := v1.Group("/posts")
		{
			// Public post routes
			posts.GET("", postHandler.GetAllPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.GET("/:id/images", postImageHandler.GetPostImagesByPostID)

			// Protected post routes
			postsProtected := posts.Group("")
			postsProtected.Use(middleware.AuthMiddleware())
			{
				postsProtected.POST("", postHandler.CreatePost)
				postsProtected.GET("/user", postHandler.GetUserPosts)
				postsProtected.PATCH("/:id", postHandler.UpdatePost)
				postsProtected.DELETE("/:id", postHandler.DeletePost)
			}
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Admin only routes
			admin := protected.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				// Personal info management
				admin.PATCH("/personal", personalInfoHandler.UpsertPersonalInfo)
				
				// Contact info management
				admin.PATCH("/contact", contactInfoHandler.UpsertContactInfo)
				
				// Social links management
				admin.POST("/social", socialLinkHandler.CreateSocialLink)
				admin.PATCH("/social/:id", socialLinkHandler.UpdateSocialLink)
				admin.DELETE("/social/:id", socialLinkHandler.DeleteSocialLink)
			}

			// User routes
			users := protected.Group("/users")
			{
				users.GET("/", nil)    // TODO: Implement user listing
				users.GET("/:id", nil) // TODO: Implement user details
				users.PATCH("/:id", nil) // TODO: Implement user update
				users.DELETE("/:id", nil) // TODO: Implement user deletion
			}

			// Role routes
			roles := protected.Group("/roles")
			roles.Use(middleware.RoleMiddleware("admin"))
			{
				roles.POST("/", nil)      // TODO: Implement role creation
				roles.GET("/", nil)       // TODO: Implement role listing
				roles.GET("/:id", nil)    // TODO: Implement role details
				roles.PATCH("/:id", nil)    // TODO: Implement role update
				roles.DELETE("/:id", nil) // TODO: Implement role deletion
			}

			// Post image routes
			postImages := protected.Group("/post-images")
			{
				postImages.POST("", postImageHandler.CreatePostImage)
				postImages.PUT("/:id", postImageHandler.UpdatePostImage)
				postImages.DELETE("/:id", postImageHandler.DeletePostImage)
			}
		}
	}
}
