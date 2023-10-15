package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/ads"
	"homework10/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	AuthorID  int64     `json:"author_id"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"creation_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type deleteAdRequest struct {
	UserID int64 `json:"user_id"`
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
			CreatedAt: ad.CreatedAt,
			UpdatedAt: ad.UpdatedAt,
		},
		"error": nil,
	}
}

func AdsSuccessResponse(ads *[]ads.Ad) *gin.H {
	var adsResponseData []adResponse
	for i := 0; i < len(*ads); i++ {
		adsResponseData = append(adsResponseData, adResponse{
			ID:        (*ads)[i].ID,
			Title:     (*ads)[i].Title,
			Text:      (*ads)[i].Text,
			AuthorID:  (*ads)[i].AuthorID,
			Published: (*ads)[i].Published,
			CreatedAt: (*ads)[i].CreatedAt,
			UpdatedAt: (*ads)[i].UpdatedAt,
		})
	}

	return &gin.H{
		"data":  adsResponseData,
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		"error": nil,
	}
}

func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
