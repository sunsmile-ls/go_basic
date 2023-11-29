package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func func1(c *gin.Context) {
	fmt.Println("func1")
}

func func2(c *gin.Context) {
	fmt.Println("func2 before")
	c.Next()
	fmt.Println("func2 after")
}

func func3(c *gin.Context) {
	fmt.Println("func3")
}
func func4(c *gin.Context) {
	fmt.Println("func4")
	c.Set("name", "sunSmile")
}
func func5(c *gin.Context) {
	fmt.Println("func5")
	val, ok := c.Get("name")
	if ok {
		vStr := val.(string)
		fmt.Println(vStr)
	}
}
func main() {
	r := gin.Default()
	r.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "search",
		})
	})
	r.GET("/blog", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "blog",
		})
	})
	r.GET("/support", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "support",
		})
	})
	shopGroup := r.Group("/shop", func1, func2)
	shopGroup.Use(func3)
	{
		shopGroup.GET("/index", func4, func5)
		shopGroup.GET("/sunsmile", func(c *gin.Context) {
			fmt.Println("func sunsmile")
		})
	}
	r.Run(":8080")
}
