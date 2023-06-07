package main

import (
	"concurrency-poc/pkg/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"wg/aggregation"
)

func main() {
	start := time.Now()

	ctx := context.Background()

	user, err := fetchUserData(ctx)

	fmt.Printf("got user %v, with error %v\n", user, err)
	fmt.Printf("execution finished in %v\n", time.Since(start))
}

func fetchUserData(ctx context.Context) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
	defer cancel()

	wg := &sync.WaitGroup{}
	userData := aggregation.NewUserData()

	done := make(chan bool)
	wg.Add(2)
	go userData.GetFriends(ctx, wg)
	go userData.GetUserName(ctx, wg)
	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			return models.User{}, fmt.Errorf("context closed while fetching user data, %v", ctx.Err())
		case <-done:
			return userData.Read(), nil
		default:
			log.Println("Just keep waiting...")
			time.Sleep(10 * time.Millisecond)
		}
	}

}
