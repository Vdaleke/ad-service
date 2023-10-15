package app

import (
	"github.com/stretchr/testify/mock"
	"homework10/internal/ads"
	"homework10/internal/mocks"
	"testing"
)

func TestAdService_CreateAd(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("Add", mock.Anything).
		Return(nil)
	repo.On("CheckIdExist", mock.Anything).
		Return(true)

	nextId := int64(0)
	repo.On("GetNextId", mock.Anything).
		Return(func() int64 { defer func() { nextId++ }(); return nextId })

	app := NewApp(repo, repo)

	user, _ := app.CreateUser("test user", "test@email")

	type Test struct {
		Name      string
		ad        ads.Ad
		ExpectErr error
	}

	tests := [...]Test{
		{"Add first ad", ads.Ad{
			Title:    "test1",
			Text:     "text test1",
			AuthorID: user.ID,
		}, nil},
		{"Add second ad", ads.Ad{
			Title:    "test2",
			Text:     "text test2",
			AuthorID: user.ID,
		}, nil},
	}

	for _, test := range tests {
		_, err := app.CreateAd(test.ad.Title, test.ad.Text, test.ad.AuthorID)
		if err != test.ExpectErr {
			t.Fatalf(`test %q: expect %v got %v`, test.Name, test.ExpectErr, err)
		}
	}
}
