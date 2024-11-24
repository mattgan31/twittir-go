package controllers

import (
	"context"
	"io"
	"net/http"
	"twittir-go/config"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	// Set HTTP 303 status and redirect to the generated URL
	c.Redirect(http.StatusSeeOther, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "randomstate" {
		c.String(http.StatusBadRequest, "States don't Match!!")
		return
	}

	code := c.Query("code")

	googlecon := config.GoogleConfig()

	// Exchange code for token
	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Code-Token Exchange Failed")
		return
	}

	// Fetch user data with the access token
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "User Data Fetch Failed")
		return
	}
	defer resp.Body.Close()

	// Read and parse the user data response
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "JSON Parsing Failed")
		return
	}

	// Return user data as a string
	c.String(http.StatusOK, string(userData))
}
