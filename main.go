package main

import (
	"mirauserlab/controllers"
	"mirauserlab/middlewares"
	"mirauserlab/models"

	"mirauserlab/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	models.ConnectDataBase()

	r := gin.Default()

	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	// @Security BearerAuth
	public.GET("/userget", controllers.UserGet)
	protected := r.Group("/api/admin")
	docs.SwaggerInfo.BasePath = "/"
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.PATCH("/createpassword/:id", controllers.CreatePassword)
	protected.GET("/getallusers", controllers.GetAllUsers)
	protected.GET("/currentuser", controllers.CurrentUser)
	protected.POST("/createemployee", controllers.CreateEmployee)
	protected.POST("/createassets", controllers.CreateAssets)
	protected.GET("/getasset/:id", controllers.GetAssetById)
	protected.GET("/getallassets", controllers.GetAllAssets)
	protected.DELETE("/deleteasset/:id", controllers.DeleteAssetById)
	protected.PATCH("/update/asset/:id/:assetId", controllers.UpdateEmployeeForAssetId)
	protected.PATCH("/deregister-asset/:id", controllers.DeRegisterByAssetId)
	protected.POST("/logout", controllers.Logout)

	//	protected.PUT("/updates/:eid/:assetId", controllers.UpdateEmployeeForAssetId)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	r.POST("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	r.PUT("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	r.PATCH("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	r.DELETE("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	r.Run(":8080")
}
