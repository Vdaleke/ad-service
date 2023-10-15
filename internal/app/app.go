package app

import (
	validator "github.com/Vdaleke/ad-validation"
	"github.com/pkg/errors"
	"homework10/internal/ads"
	"homework10/internal/users"
	"strings"
	"time"
)

type App interface {
	CreateAd(title string, text string, userId int64) (ads.Ad, error)
	ChangeAdStatus(adId int64, userId int64, published bool) (ads.Ad, error)
	UpdateAd(adId int64, userId int64, title string, text string) (ads.Ad, error)
	GetAd(adId int64) (ads.Ad, error)
	DeleteAd(adId int64, userId int64) error
	ListAds(pubFilter bool, userFilter int64, timeFilter time.Time) ([]ads.Ad, error)
	SearchAds(pattern string) ([]ads.Ad, error)

	CreateUser(name string, email string) (users.User, error)
	UpdateUser(userId int64, name string, email string) (users.User, error)
	GetUser(userId int64) (users.User, error)
	DeleteUser(userId int64) error
}

type Repository interface {
	Add(e interface{}) error
	Update(id int64, ad interface{}) error
	Get(id int64) (interface{}, error)
	Delete(id int64) error
	CheckIdExist(id int64) bool
	GetNextId() int64
	GetArray() []interface{}
}

func NewApp(adRepo Repository, userRepo Repository) App {
	return &AdService{ads: adRepo, users: userRepo}
}

type AdService struct {
	ads   Repository
	users Repository
}

var PermissionDenied = errors.New("the user does not have enough permission to edit the ad")
var DefunctUser = errors.New("there is no user with this ID")
var DefunctAd = errors.New("there is no ad with this ID")

func (a *AdService) CreateAd(title string, text string, userId int64) (ads.Ad, error) {
	if !a.users.CheckIdExist(userId) {
		return ads.Ad{}, DefunctUser
	}

	ad := ads.Ad{ID: a.ads.GetNextId(), Title: title, Text: text, AuthorID: userId, Published: false, CreatedAt: time.Now().UTC()}

	err := validator.ValidateAd(title, text)
	if err != nil {
		return ad, err
	}

	return ad, a.ads.Add(ad)
}

func (a *AdService) ChangeAdStatus(adId int64, userId int64, published bool) (ads.Ad, error) {
	if !a.users.CheckIdExist(userId) {
		return ads.Ad{}, DefunctUser
	}

	if !a.ads.CheckIdExist(adId) {
		return ads.Ad{}, DefunctAd
	}

	res, err := a.ads.Get(adId)
	ad := res.(ads.Ad)

	if err != nil {
		return ad, err
	}

	if ad.AuthorID != userId {
		return ad, PermissionDenied
	}

	ad.Published = published

	return ad, a.ads.Update(adId, ad)
}

func (a *AdService) UpdateAd(adId int64, userId int64, title string, text string) (ads.Ad, error) {
	if !a.users.CheckIdExist(userId) {
		return ads.Ad{}, DefunctUser
	}
	if !a.ads.CheckIdExist(adId) {
		return ads.Ad{}, DefunctAd
	}

	res, err := a.ads.Get(adId)
	ad := res.(ads.Ad)

	if err != nil {
		return ad, err
	}

	if ad.AuthorID != userId {
		return ad, PermissionDenied
	}

	err = validator.ValidateAd(title, text)
	if err != nil {
		return ad, err
	}

	ad.Title = title
	ad.Text = text
	ad.UpdatedAt = time.Now().UTC()

	return ad, a.ads.Update(adId, ad)
}

func (a *AdService) GetAd(adId int64) (ads.Ad, error) {
	if !a.ads.CheckIdExist(adId) {
		return ads.Ad{}, DefunctAd
	}

	res, err := a.ads.Get(adId)
	ad := res.(ads.Ad)

	return ad, err
}

func (a *AdService) DeleteAd(adId int64, userId int64) error {
	if !a.users.CheckIdExist(userId) {
		return DefunctUser
	}
	if !a.ads.CheckIdExist(adId) {
		return DefunctAd
	}

	res, err := a.ads.Get(adId)
	ad := res.(ads.Ad)

	if err != nil {
		return err
	}
	if ad.AuthorID != userId {
		return PermissionDenied
	}

	return a.ads.Delete(adId)
}

func (a *AdService) ListAds(pubFilter bool, userFilter int64, timeFilter time.Time) ([]ads.Ad, error) {
	res := a.ads.GetArray()
	adsArray := make([]ads.Ad, 0)

	for _, e := range res {
		ad := e.(ads.Ad)
		if pubFilter == ad.Published &&
			(userFilter == -1 || userFilter == ad.AuthorID) &&
			(timeFilter.IsZero() || timeFilter.Equal(ad.CreatedAt)) {
			adsArray = append(adsArray, ad)
		}
	}

	return adsArray, nil
}

func (a *AdService) SearchAds(pattern string) ([]ads.Ad, error) {
	allAds := a.ads.GetArray()
	filteredAds := make([]ads.Ad, 0)

	for _, e := range allAds {
		ad := e.(ads.Ad)
		if strings.Contains(ad.Title, pattern) {
			filteredAds = append(filteredAds, ad)
		}
	}

	return filteredAds, nil
}

func (a *AdService) CreateUser(name string, email string) (users.User, error) {
	user := users.User{ID: a.users.GetNextId(), Name: name, Email: email}

	return user, a.users.Add(user)
}

func (a *AdService) UpdateUser(userId int64, name string, email string) (users.User, error) {
	if !a.users.CheckIdExist(userId) {
		return users.User{}, DefunctUser
	}

	res, err := a.users.Get(userId)
	user := res.(users.User)

	if err != nil {
		return user, err
	}

	user.Name = name
	user.Email = email

	return user, a.users.Update(userId, user)
}

func (a *AdService) GetUser(userId int64) (users.User, error) {
	if !a.users.CheckIdExist(userId) {
		return users.User{}, DefunctUser
	}

	res, err := a.users.Get(userId)
	user := res.(users.User)

	if err != nil {
		return user, err
	}

	return user, err
}

func (a *AdService) DeleteUser(userId int64) error {
	for _, e := range a.ads.GetArray() {
		ad := e.(ads.Ad)
		if ad.AuthorID == userId {
			err := a.DeleteAd(ad.ID, userId)
			if err != nil {
				return err
			}
		}
	}

	return a.users.Delete(userId)
}
