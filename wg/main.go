package main

import (
	"concurrency-poc/pkg/models"
	"fmt"
	"sync"
	"time"
	"wg/aggregation"
)

func main() {
	// Sync - ~70ms
	SyncWork()

	// Async - ~50ms
	AsyncWork()
}

func SyncWork() {
	start := time.Now()
	user := aggregation.SyncGetUserData()
	friends := aggregation.SyncGetUserFriends()

	fmt.Printf("User: %q has %d friends, determined in %v\n", user.Name, len(friends), time.Since(start))
}

func AsyncWork() {
	start := time.Now()

	// Join mechanisms
	var wg sync.WaitGroup
	var user models.User

	// TODO: Extract and leverage UserDataAggregator interface
	userData := aggregation.NewUserData()
	wg.Add(2)

	// Fork
	go userData.GetUserName(&wg)
	go userData.GetFriends(&wg)

	// Join
	wg.Wait()
	user = userData.Read()

	fmt.Printf("User: %q has %d friends, determined in %v\n", user.Name, len(user.Friends), time.Since(start))
}
