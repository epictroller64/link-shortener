package handlers

import (
	"fmt"
	"link-shortener-backend/src/repository"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Redirect is a function that redirects a user to a new URL based on the redirect rules
func Redirect(c *gin.Context) {
	headers := c.Request.Header
	cookies := c.Request.Cookies()
	shortId := c.Param("shortId")
	link, err := repository.GetLinkByShortId(shortId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}
	// Record a click to the database for statistics
	// Country will be added later using a 3rd party service or IP geolocation
	click := repository.Click{
		LinkID:    link.ID,
		UserAgent: headers.Get("User-Agent"),
		Referer:   headers.Get("Referer"),
		CreatedAt: time.Now(),
		IP:        c.ClientIP(),
	}
	repository.CreateClick(click)
	// Update the link's click count by 1
	err = repository.UpdateLinkClickCount(link.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update link click count"})
		return
	}
	redirects, err := repository.GetRedirectsByLinkID(strconv.Itoa(link.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get redirects"})
		return
	}
	// Check if headers match any of the redirect rules
	headerRedirect := headerCheck(headers, redirects)
	if headerRedirect != nil {
		c.Redirect(http.StatusFound, headerRedirect.RedirectURL)
		return
	}
	cookieRedirect := cookieCheck(cookies, redirects)
	if cookieRedirect != nil {
		c.Redirect(http.StatusFound, cookieRedirect.RedirectURL)
		return
	}
	// No redirect matched, redirect to the original link
	c.Redirect(http.StatusFound, link.Original)

}

func cookieCheck(cookies []*http.Cookie, redirects []repository.Redirect) *repository.Redirect {
	for _, redirect := range redirects {
		if redirect.TargetType == "cookie" && redirect.TargetName != nil {
			for _, cookie := range cookies {
				if cookie.Name == *redirect.TargetName {
					switch redirect.TargetMethod {
					case "match":
						if redirect.TargetValue != nil && *redirect.TargetValue == cookie.Value {
							return &redirect
						}
					case "regex":
						if redirect.TargetValue != nil && regexMatch(cookie.Value, *redirect.TargetValue) {
							return &redirect
						}
					case "contains":
						if redirect.TargetValue != nil && containsMatch(cookie.Value, *redirect.TargetValue) {
							return &redirect
						}
					case "startsWith":
						if redirect.TargetValue != nil && startsWithMatch(cookie.Value, *redirect.TargetValue) {
							return &redirect
						}
					case "endsWith":
						if redirect.TargetValue != nil && endsWithMatch(cookie.Value, *redirect.TargetValue) {
							return &redirect
						}
					}
				}
			}
		}
	}
	return nil
}

func headerCheck(headers http.Header, redirects []repository.Redirect) *repository.Redirect {
	for _, redirect := range redirects {
		if redirect.TargetType == "header" && redirect.TargetName != nil {
			for headerKey, headerValues := range headers {
				if *redirect.TargetName == headerKey {
					for _, headerValue := range headerValues {
						switch redirect.TargetMethod {
						case "match":
							if redirect.TargetValue != nil && *redirect.TargetValue == headerValue {
								return &redirect
							}
						case "regex":
							if redirect.TargetValue != nil && regexMatch(headerValue, *redirect.TargetValue) {
								return &redirect
							}
						case "contains":
							if redirect.TargetValue != nil && containsMatch(headerValue, *redirect.TargetValue) {
								return &redirect
							}
						case "startsWith":
							if redirect.TargetValue != nil && startsWithMatch(headerValue, *redirect.TargetValue) {
								return &redirect
							}
						case "endsWith":
							if redirect.TargetValue != nil && endsWithMatch(headerValue, *redirect.TargetValue) {
								return &redirect
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func startsWithMatch(headerValue string, targetValue string) bool {
	return strings.HasPrefix(headerValue, targetValue)
}

func endsWithMatch(headerValue string, targetValue string) bool {
	return strings.HasSuffix(headerValue, targetValue)
}

func containsMatch(headerValue string, targetValue string) bool {
	return strings.Contains(headerValue, targetValue)
}

func regexMatch(headerValue string, regex string) bool {
	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return re.MatchString(headerValue)
}
