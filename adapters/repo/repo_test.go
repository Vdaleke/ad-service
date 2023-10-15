package repo

import (
	"fmt"
	"homework10/internal/app"
	"testing"
)

func TestRepo_Add(t *testing.T) {
	type Test struct {
		Name   string
		Item   interface{}
		Expect error
	}

	tests := [...]Test{
		{"Add 7", 7, nil},
		{"Add 9", 9, nil},
	}

	repo := New()

	for _, test := range tests {
		got := repo.Add(test.Item)
		if got != test.Expect {
			t.Fatalf(`test %q: expect %v got %v`, test.Name, test.Expect, got)
		}
	}
}

func FuzzRepo_Get(f *testing.F) {
	repo := New()

	for i := 1; i <= 100; i++ {
		_ = repo.Add(i)
	}

	f.Fuzz(func(t *testing.T, id int64) {
		_, err := repo.Get(id)
		var expectErr error
		if repo.CheckIdExist(id) {
			expectErr = nil
		} else {
			expectErr = DefunctEntity
		}

		if err != expectErr {
			t.Errorf("For (%d) Expect: %s, but got: %s", id, expectErr, err)
		}
	})
}

func TestRepo_Update(t *testing.T) {
	teardown := func() {
		fmt.Println("End testing update")
	}

	var repo app.Repository

	setup := func(t *testing.T) {
		t.Cleanup(teardown)
		repo = New()
		_ = repo.Add(1)
		_ = repo.Add(2)
		fmt.Println("Set up repo")
	}

	type Test struct {
		Name   string
		Pos    int64
		Item   interface{}
		Expect error
	}

	tests := [...]Test{
		{"Update item at position 1", 0, 2, nil},
		{"Update item at position 2", 1, 3, nil},
		{"Update non-existent item at position 3", 2, 4, DefunctEntity},
	}

	t.Run("with Cleanup", func(t *testing.T) {
		setup(t)
		for _, test := range tests {
			err := repo.Update(test.Pos, test.Item)
			if err != test.Expect {
				t.Fatalf(`test %q: expect %v got %v`, test.Name, test.Expect, err)
			}
		}

		for _, test := range tests {
			item, err := repo.Get(test.Pos)
			if err == nil && item != test.Item {
				t.Fatalf(`test %q: expect %v at poition %d got %v`, test.Name, test.Item, test.Pos, item)
			}
		}
	})
}
