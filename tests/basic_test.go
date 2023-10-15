package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAd(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(createdUser.Data.ID))
	assert.False(t, response.Data.Published)
	assert.False(t, response.Data.CreatedAt.IsZero())
	assert.True(t, response.Data.UpdatedAt.IsZero())
}

func TestChangeAdStatus(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.False(t, response.Data.CreatedAt.IsZero())
	assert.True(t, response.Data.UpdatedAt.IsZero())

	response, err = client.updateAd(createdUser.Data.ID, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
	assert.False(t, response.Data.CreatedAt.IsZero())
	assert.False(t, response.Data.UpdatedAt.IsZero())
}

func TestListAds(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.CreateAd(createdUser.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestFilterListAds(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	filteredAd, err := client.CreateAd(createdUser.Data.ID, "hello2", "world2")
	assert.NoError(t, err)

	createdUser2, err := client.CreateUser("Test User2", "test2@testing.ru")
	assert.NoError(t, err)

	response, err = client.CreateAd(createdUser2.Data.ID, "hello3", "world3")
	assert.NoError(t, err)

	ads, err := client.filterListAds("?published=false&user_id=0")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, filteredAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, filteredAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, filteredAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, filteredAd.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestSearchAds(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	ads, err := client.searchAds("ell")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, response.Data.ID)
	assert.Equal(t, ads.Data[0].Title, response.Data.Title)
	assert.Equal(t, ads.Data[0].Text, response.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, response.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestGetAd(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	ad, err := client.getAd(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, publishedAd.Data.ID)
	assert.Equal(t, ad.Data.Title, publishedAd.Data.Title)
	assert.Equal(t, ad.Data.Text, publishedAd.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ad.Data.Published)
}

func TestDeleteAd(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(createdUser.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	err = client.deleteAd(publishedAd.Data.ID, createdUser.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(publishedAd.Data.ID)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser(t *testing.T) {
	client := GetTestClient()

	response, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Name, "Test User")
	assert.Equal(t, response.Data.Email, "test@testing.ru")
}

func TestUpdateUser(t *testing.T) {
	client := GetTestClient()

	response, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.ID, "Test User 2", "test2@testing.ru")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Name, "Test User 2")
	assert.Equal(t, response.Data.Email, "test2@testing.ru")
}

func TestGetUser(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.updateUser(createdUser.Data.ID, "Test User 2", "test2@testing.ru")
	assert.NoError(t, err)

	user, err := client.getUser(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Data.ID, response.Data.ID)
	assert.Equal(t, user.Data.Name, response.Data.Name)
	assert.Equal(t, user.Data.Email, response.Data.Email)
}

func TestDeleteUser(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	response, err := client.updateUser(createdUser.Data.ID, "Test User 2", "test2@testing.ru")
	assert.NoError(t, err)

	err = client.deleteUser(response.Data.ID)
	assert.NoError(t, err)

	_, err = client.getUser(response.Data.ID)
	assert.ErrorIs(t, err, ErrBadRequest)
}
