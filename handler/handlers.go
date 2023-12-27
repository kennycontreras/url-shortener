package handler

import (
	"eenrique/url-shortener/shortener"
	"eenrique/url-shortener/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UrlCreationRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	shortenedUrl := shortener.GenerateShortLink(creationRequest.OriginalUrl, creationRequest.UserId)
	store.SaveUrlMapping(shortenedUrl, creationRequest.OriginalUrl, creationRequest.UserId)

	host := "http://localhost:9808/"
	c.JSON(200, gin.H{
		"message":   "Short url created successfully",
		"short_url": host + shortenedUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortenedUrl := c.Param("shortUrl")
	initialUrl := store.RetrieveInitialUrl(shortenedUrl)
	c.Redirect(302, initialUrl)
}
