package main

import (
	"concurrency-poc/pkg/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"wg/pkg/aggregation"
)

const EXAMPLE_NAME = "wg/examples/data-aggregation"

func main() {
	log.Println("Running example:", EXAMPLE_NAME)
	// Sync - ~70ms
	SyncWork()

	// Async - ~50ms
	AsyncWork()
	log.Println("Finished example:", EXAMPLE_NAME)
}

func SyncWork() {
	start := time.Now()
	user := aggregation.SyncGetUserData()
	friends := aggregation.SyncGetUserFriends()

	fmt.Printf("[SYNC] - User: %q has %d friends, determined in %v\n", user.Name, len(friends), time.Since(start))
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
	go userData.GetUserName(context.TODO(), &wg)
	go userData.GetFriends(context.TODO(), &wg)

	// Join
	wg.Wait()
	user = userData.Read()

	fmt.Printf("[ASYNC] - User: %q has %d friends, determined in %v\n", user.Name, len(user.Friends), time.Since(start))
}
