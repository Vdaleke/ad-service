package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd_EmptyTitle(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	_, err = client.CreateAd(createdUser.Data.ID, "", "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongTitle(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	title := strings.Repeat("a", 101)

	_, err = client.CreateAd(createdUser.Data.ID, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_EmptyText(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	_, err = client.CreateAd(createdUser.Data.ID, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongText(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	text := strings.Repeat("a", 501)

	_, err = client.CreateAd(createdUser.Data.ID, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyTitle(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	resp, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(createdUser.Data.ID, resp.Data.ID, "", "new_world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongTitle(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	resp, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	title := strings.Repeat("a", 101)

	_, err = client.updateAd(createdUser.Data.ID, resp.Data.ID, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyText(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	resp, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(createdUser.Data.ID, resp.Data.ID, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongText(t *testing.T) {
	client := GetTestClient()

	createdUser, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	text := strings.Repeat("a", 501)

	resp, err := client.CreateAd(createdUser.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(createdUser.Data.ID, resp.Data.ID, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}
