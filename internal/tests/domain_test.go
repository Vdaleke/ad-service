package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := GetTestClient()

	createdUser1, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	resp, err := client.CreateAd(createdUser1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	createdUser2, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(createdUser2.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := GetTestClient()

	createdUser1, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	resp, err := client.CreateAd(createdUser1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	createdUser2, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	_, err = client.updateAd(createdUser2.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := GetTestClient()

	createdUser1, err := client.CreateUser("Test User", "test@testing.ru")
	assert.NoError(t, err)

	resp, err := client.CreateAd(createdUser1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.CreateAd(createdUser1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.CreateAd(createdUser1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}
