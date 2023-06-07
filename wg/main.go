package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Sync - ~70ms
	SyncWork()

	// Async - ~50ms
	AsyncWork()
}

func SyncWork() {
	start := time.Now()
	user := getUserData()
	friends := getUserFriends()

	fmt.Printf("User: %q has %d friends, determined in %v\n", user.Name, len(friends), time.Since(start))
}

func AsyncWork() {
	start := time.Now()

	// Join mechanisms
	var wg sync.WaitGroup
	var userData UserDataAggregator = &UserData{
		name:    make(chan string, 1),
		friends: make(chan []User, 1),
	}

	var user User
	wg.Add(2)

	// Fork
	go userData.GetUserName(&wg)
	go userData.GetFriends(&wg)

	// Join
	wg.Wait()
	user = userData.Read()

	fmt.Printf("User: %q has %d friends, determined in %v\n", user.Name, len(user.Friends), time.Since(start))
}
