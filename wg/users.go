package main

import (
	"sync"
	"time"
)

type (
	UserDataAggregator interface {
		GetUserName(*sync.WaitGroup)
		GetFriends(*sync.WaitGroup)
		Read() User
	}

	UserData struct {
		name    chan string
		friends chan []User
	}

	User struct {
		Name    string
		Friends []User
	}
)

func (u *UserData) GetUserName(wg *sync.WaitGroup) {
	defer wg.Done()

	u.name <- getUserData().Name
}

func (u *UserData) GetFriends(wg *sync.WaitGroup) {
	defer wg.Done()

	u.friends <- getUserFriends()
}

func (u UserData) Read() (user User) {
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

func getUserData() User {
	time.Sleep(time.Millisecond * 20)

	return User{
		Name: "Josh",
	}
}

func getUserFriends() []User {
	time.Sleep(time.Millisecond * 50)

	return []User{
		{Name: "Bob"},
		{Name: "Amy"},
	}
}
