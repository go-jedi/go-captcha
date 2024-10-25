package main

import (
	"log"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// Маршрут для генерации CAPTCHA
	router.GET("/captcha/:captchaId", func(c *gin.Context) {
		captchaId := c.Param("captchaId")
		c.Header("Content-Type", "image/png")
		if err := captcha.WriteImage(c.Writer, captchaId, 240, 80); err != nil {
			log.Println("error write image")
		}
	})

	// Маршрут для генерации нового CAPTCHA ID
	router.GET("/captcha/new", func(c *gin.Context) {
		captchaId := captcha.New()
		c.JSON(http.StatusOK, gin.H{"captcha_id": captchaId})
	})

	// Маршрут для проверки CAPTCHA
	router.POST("/captcha/verify", func(c *gin.Context) {
		var request struct {
			CaptchaId    string `json:"captcha_id"`
			CaptchaValue string `json:"captcha_value"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if captcha.VerifyString(request.CaptchaId, request.CaptchaValue) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failure"})
		}
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
