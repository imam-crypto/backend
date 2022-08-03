package main

import (
	"backend/auth"
	"backend/handler"
	"backend/helper"
	"backend/user"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func main() {
	dsn := "host=localhost user=postgres password=root dbname=db_backend port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//er := db.AutoMigrate(user.User{})
	//if er != nil {
	//	log.Fatal(er)
	//}
	//fmt.Println("Migrated")
	authService := auth.NewService()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)
	//user, err := userRepository.GetUsers()
	//token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IpXVCJ9.eyJ1c2VyX2lkIjozMn0.0fyseXfC6MpSkY8X8Rnr2L_pK9XqqYIBg4QtSy3BQi4")

	if err != nil {
		fmt.Println("error")
	}
	//fmt.Println("data di main", user)

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.DELETE("/delete/:id", authMiddleware(authService, userService), userHandler.DeleteUser)
	api.PUT("/updated/:id", authMiddleware(authService, userService), userHandler.Update)
	api.GET("/usersPdf", authMiddleware(authService, userService), userHandler.GetUsers)
	api.PUT("/change-password/", authMiddleware(authService, userService), userHandler.ChangesPassword)
	api.GET("/paginations", authMiddleware(authService, userService), userHandler.PaginationUser)
	//api.GET("/users", func(c *gin.Context) {
	//
	//	var users []user.User
	//	page, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))
	//	perPage := 5
	//	//Perpage, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	//	var total int64
	//
	//	sql := "SELECT * FROM users"
	//	db.Raw(sql).Count(&total)
	//	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perPage, (page-1)*perPage)
	//	db.Raw(sql).Scan(&users)
	//
	//	c.JSON(http.StatusOK, gin.H{
	//		"data":  users,
	//		"total": total,
	//		"page":  page,
	//	})
	//})
	router.StaticFS("/assets", gin.Dir("file", true))
	router.Run()

	fmt.Println("Connected To db_backend", db)
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse(http.StatusUnauthorized, false, "unauthorize", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse(http.StatusUnauthorized, false, "unauthorize", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse(http.StatusUnauthorized, false, "unauthorize", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(claim["user_id"].(float64))
		user, err := userService.FindById(userID)
		if err != nil {
			response := helper.APIResponse(http.StatusUnauthorized, false, "unauthorize", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("current_user", user)
	}

}
