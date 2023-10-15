package repo

import (
	"github.com/pkg/errors"
	"homework10/internal/app"
	"sync"
)

func New() app.Repository {
	return &Repo{storage: make(map[int64]interface{}), nextNum: 0}
}

type Repo struct {
	storage map[int64]interface{}
	nextNum int64
	mu      sync.Mutex
}

var DefunctEntity = errors.New("there is no entity with this id")

func (a *Repo) Add(e interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.storage[a.nextNum] = e
	a.nextNum++
	return nil
}

func (a *Repo) Update(id int64, e interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.storage[id]

	if !exists {
		return DefunctEntity
	}

	a.storage[id] = e
	return nil
}

func (a *Repo) Get(id int64) (interface{}, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.storage[id]

	if !exists {
		return nil, DefunctEntity
	}

	return a.storage[id], nil
}

func (a *Repo) Delete(id int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.storage[id]

	if !exists {
		return DefunctEntity
	}

	delete(a.storage, id)

	return nil
}

func (a *Repo) CheckIdExist(id int64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, exists := a.storage[id]

	return exists
}

func (a *Repo) GetNextId() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.nextNum
}

func (a *Repo) GetArray() []interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	arr := make([]interface{}, 0)

	for _, e := range a.storage {
		arr = append(arr, e)
	}

	return arr
}
