package aggregation

import (
	"concurrency-poc/pkg/models"
	"context"
	"sync"
	"time"
)

type (
	UserDataAggregator interface {
		GetUserName(*sync.WaitGroup)
		GetFriends(*sync.WaitGroup)
		Read() models.User
	}

	UserData struct {
		name    chan string
		friends chan []models.User
	}
)

func NewUserData() UserData {
	return UserData{
		name:    make(chan string, 1),
		friends: make(chan []models.User, 1),
	}
}

func (u *UserData) GetUserName(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	u.name <- getUserData().Name
}

func (u *UserData) GetFriends(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	u.friends <- getUserFriends()
}

func (u UserData) Read() (user models.User) {
	var received int = 0
	for received != 2 {
		select {
		case friends := <-u.friends:
			user.Friends = friends
			received++
		case name := <-u.name:
			user.Name = name
			received++
		}
	}

	return user
}

func SyncGetUserData() models.User {
	return getUserData()
}

func SyncGetUserFriends() []models.User {
	return getUserFriends()
}

func getUserData() models.User {
	time.Sleep(time.Millisecond * 20)

	return models.User{
		Name: "Josh",
	}
}

func getUserFriends() []models.User {
	time.Sleep(time.Millisecond * 50)

	return []models.User{
		{Name: "Bob"},
		{Name: "Amy"},
	}
}
