package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Post model defined
type Posts struct {
    ID        		int    		`gorm:"primaryKey;autoIncrement"`
    Title     		string 		`gorm:"size:200;not null"`
    Content   		string 		`gorm:"type:text;not null"`
    Category  		string 		`gorm:"size:100;not null"`
	Created_date 	time.Time 	`gorm:"autoCreateTime"`
	Updated_date 	time.Time 	`gorm:"autoCreateTime"`
	Status   		string 		`gorm:"size:100;not null"`
}

// Initialize the database
func initDB() {
    dsn := "root:password@tcp(127.0.0.1:3307)/vision_db?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Auto-migrate the schema
    if err := DB.AutoMigrate(&Posts{}); err != nil {
        panic("failed to migrate schema")
    }
}

func main() {
    // Initialize database
    initDB()

    // Create a new Gin router
    router := gin.Default()

    // Define endpoints
    router.POST("/article", createPost)
    router.GET("/article/:id", getPost)
    router.PUT("/article/:id", updatePost)
    router.DELETE("/article/:id", deletePost)
    router.GET("/article/list", listPosts)

    // Start the server on port 8080
    router.Run(":8080")
}

func createPost(c *gin.Context) {
    var post Posts
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := DB.Create(&post)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusCreated, post)
}

func getPost(c *gin.Context) {
    id := c.Param("id")
    var post Posts
    result := DB.First(&post, id)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    c.JSON(http.StatusOK, post)
}

func updatePost(c *gin.Context) {
    id := c.Param("id")
    var post Posts
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := DB.Model(&post).Where("id = ?", id).Updates(post)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    c.JSON(http.StatusOK, post)
}

func deletePost(c *gin.Context) {
    id := c.Param("id")
    result := DB.Delete(&Posts{}, id)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    c.JSON(http.StatusNoContent, nil)
}

func listPosts(c *gin.Context) {
	// Get limit and offset from the query parameters
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")) // Default limit is 20
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0")) // Default offset is 0
	offset = (offset - 1) * limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	var posts []Posts
	result := DB.Limit(limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}