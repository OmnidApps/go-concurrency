package aggregation_test

import (
	"concurrency-poc/pkg/models"
	"context"
	"sync"
	"testing"
	"wg/pkg/aggregation"
)

// Checks if User slice s contains a user with name e
func contains(s []models.User, e string) bool {
	for _, a := range s {
		if a.Name == e {
			return true
		}
	}
	return false
}

func TestUserData(t *testing.T) {
	t.Run("GetUserName", func(t *testing.T) {
		u := aggregation.NewUserData()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go u.GetUserName(context.TODO(), wg)
		wg.Wait()
		u.Done()

		// TODO: Move away from static data
		expected := "Josh"
		got := u.Read().Name
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("GetFriends", func(t *testing.T) {
		u := aggregation.NewUserData()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go u.GetFriends(context.TODO(), wg)
		wg.Wait()
		u.Done()

		// TODO: Move away from static data
		expected := []string{
			"Bob",
			"Amy",
		}
		got := u.Read().Friends
		for _, item := range expected {
			if !contains(got, item) {
				t.Errorf("expected %q, got %q", expected, got)
			}
		}
	})

}
