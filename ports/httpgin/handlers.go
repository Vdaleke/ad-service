package httpgin

import (
	"errors"
	validator "github.com/Vdaleke/ad-validation"
	"github.com/gin-gonic/gin"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/users"
	"net/http"
	"strconv"
	"time"
)

func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if errors.Is(err, validator.ValidationError) || errors.Is(err, app.DefunctUser) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.ChangeAdStatus(int64(adID), reqBody.UserID, reqBody.Published)

		if errors.Is(err, app.PermissionDenied) {
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		} else if errors.Is(err, validator.ValidationError) ||
			errors.Is(err, app.DefunctUser) || errors.Is(err, app.DefunctAd) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if errors.Is(err, app.PermissionDenied) {
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		} else if errors.Is(err, validator.ValidationError) ||
			errors.Is(err, app.DefunctUser) || errors.Is(err, app.DefunctAd) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.GetAd(int64(adID))

		if errors.Is(err, app.DefunctAd) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		err = a.DeleteAd(int64(adID), reqBody.UserID)
		if errors.Is(err, app.PermissionDenied) {
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		} else if errors.Is(err, validator.ValidationError) || errors.Is(err, app.DefunctUser) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ads.Ad{}))
	}
}

func listAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		pubFilter := true
		if c.Query("published") == "false" {
			pubFilter = false
		}

		userFilter, atoiErr := strconv.Atoi(c.Query("user_id"))
		if atoiErr != nil {
			userFilter = -1
		}

		timeFilter, _ := time.Parse(time.RFC3339, c.Query("creation_time"))

		ads, err := a.ListAds(pubFilter, int64(userFilter), timeFilter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(&ads))
	}
}

func searchAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		pattern := c.Param("pattern")

		ads, err := a.SearchAds(pattern)

		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(&ads))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		user, err := a.CreateUser(reqBody.Name, reqBody.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&user))
	}
}

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		user, err := a.UpdateUser(int64(userID), reqBody.Name, reqBody.Email)

		if errors.Is(err, app.DefunctUser) {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&user))
	}
}

func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		user, err := a.GetUser(int64(userID))

		if errors.Is(err, app.DefunctUser) {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&user))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		err = a.DeleteUser(int64(userID))

		if errors.Is(err, app.DefunctUser) {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&users.User{}))
	}
}
