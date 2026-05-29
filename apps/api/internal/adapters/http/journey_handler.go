package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/dto/openapi"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

const adminSessionCookieName = "admin_session"
const adminSessionTTL = time.Hour

type memoryJourneyRepository struct {
	journeys []domain.Journey
}

func newMemoryJourneyRepository() *memoryJourneyRepository {
	thumbnail := "https://picsum.photos/seed/tsr-remote-desk-2020/800/600"
	return &memoryJourneyRepository{
		journeys: []domain.Journey{
			{
				ID:        1,
				Name:      "Remote desk, first production habits",
				Timestamp: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Location:  "Bangkok",
				Thumbnail: &thumbnail,
			},
		},
	}
}

func (r *memoryJourneyRepository) List(context.Context) ([]domain.Journey, error) {
	return r.journeys, nil
}

func (r *memoryJourneyRepository) Get(_ context.Context, id int64) (domain.JourneyContent, error) {
	for _, journey := range r.journeys {
		if journey.ID == id {
			return domain.JourneyContent{
				Name:      journey.Name,
				Timestamp: journey.Timestamp,
				Content:   "# Remote setup\n\nSmall production habits started from a compact desk in Bangkok.",
			}, nil
		}
	}

	return domain.JourneyContent{}, service.ErrJourneyNotFound
}

func (r *memoryJourneyRepository) Create(_ context.Context, journey domain.CreateJourney) (int64, error) {
	id := int64(len(r.journeys) + 1)
	r.journeys = append(r.journeys, domain.Journey{
		ID:        id,
		Name:      journey.Name,
		Timestamp: time.Now().UTC(),
		Location:  journey.Location,
		Thumbnail: journey.Thumbnail,
	})

	return id, nil
}

func registerJourneyRoutes(router gin.IRoutes, repository service.JourneyRepository, adminToken string) {
	journeyService := service.NewJourneyService(repository)

	router.GET("/journeys", func(c *gin.Context) {
		journeys, err := journeyService.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list journeys"})
			return
		}

		c.JSON(http.StatusOK, toJourneysResponse(journeys))
	})

	router.GET("/journey/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid journey id"})
			return
		}

		journey, err := journeyService.Get(c.Request.Context(), id)
		if errors.Is(err, service.ErrJourneyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "journey not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get journey"})
			return
		}

		c.JSON(http.StatusOK, toJourneyContentResponse(journey))
	})

	router.POST("/journey", requireAdminSession(adminToken), func(c *gin.Context) {
		var request openapi.CreateJourneyRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid journey request"})
			return
		}

		id, err := journeyService.Create(c.Request.Context(), domain.CreateJourney{
			Name:      request.Name,
			Location:  request.Location,
			Thumbnail: request.Thumbnail,
			Content:   request.Content,
			Highlight: request.Highlight,
		})
		if errors.Is(err, service.ErrInvalidJourney) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "journey name, location, and content are required"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create journey"})
			return
		}

		c.JSON(http.StatusCreated, openapi.CreateJourneyResponse{Id: id})
	})
}

func registerAdminRoutes(router gin.IRoutes, adminToken string) {
	router.GET("/admin-verify", requireAdminSession(adminToken), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/admin-verify", func(c *gin.Context) {
		var request openapi.AdminVerifyRequest
		if err := c.ShouldBindJSON(&request); err != nil || !isValidAdminToken(request.Token, adminToken) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid admin token"})
			return
		}

		setAdminSessionCookie(c, adminToken)
		c.Status(http.StatusOK)
	})
}

func requireAdminSession(adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(adminSessionCookieName)
		if err != nil || !isValidAdminSession(cookie, adminToken) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid admin token"})
			return
		}

		c.Next()
	}
}

func isValidAdminToken(token string, adminToken string) bool {
	if adminToken == "" {
		adminToken = os.Getenv("ADMIN_TOKEN")
	}
	token = strings.TrimSpace(token)
	adminToken = strings.TrimSpace(adminToken)
	if token == "" || adminToken == "" || len(token) != len(adminToken) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(token), []byte(adminToken)) == 1
}

func setAdminSessionCookie(c *gin.Context, adminToken string) {
	expiresAt := time.Now().UTC().Add(adminSessionTTL).Unix()
	value := signAdminSession(expiresAt, adminToken)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		adminSessionCookieName,
		value,
		int(adminSessionTTL.Seconds()),
		"/",
		"",
		isSecureRequest(c),
		true,
	)
}

func isValidAdminSession(value string, adminToken string) bool {
	if adminToken == "" {
		adminToken = os.Getenv("ADMIN_TOKEN")
	}
	adminToken = strings.TrimSpace(adminToken)
	if value == "" || adminToken == "" {
		return false
	}

	var expiresAt int64
	var signature string
	if _, err := fmt.Sscanf(value, "%d.%s", &expiresAt, &signature); err != nil {
		return false
	}
	if time.Now().UTC().Unix() > expiresAt {
		return false
	}

	expectedValue := signAdminSession(expiresAt, adminToken)
	return subtle.ConstantTimeCompare([]byte(value), []byte(expectedValue)) == 1
}

func signAdminSession(expiresAt int64, adminToken string) string {
	message := strconv.FormatInt(expiresAt, 10)
	mac := hmac.New(sha256.New, []byte(adminToken))
	mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return message + "." + signature
}

func isSecureRequest(c *gin.Context) bool {
	return c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
}

func toJourneysResponse(journeys []domain.Journey) openapi.JourneysResponse {
	response := make(openapi.JourneysResponse, 0, len(journeys))

	for _, journey := range journeys {
		response = append(response, openapi.Journey{
			Id:        journey.ID,
			Name:      journey.Name,
			Timestamp: journey.Timestamp.Format(time.DateOnly),
			Location:  journey.Location,
			Thumbnail: journey.Thumbnail,
		})
	}

	return response
}

func toJourneyContentResponse(journey domain.JourneyContent) openapi.JourneyContent {
	return openapi.JourneyContent{
		Name:      journey.Name,
		Timestamp: journey.Timestamp.Format(time.DateOnly),
		Content:   journey.Content,
	}
}
