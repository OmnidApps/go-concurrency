package aggregation

import (
	"concurrency-poc/pkg/models"
	"context"
	"reflect"
	"sync"
	"time"
)

type (
	UserDataAggregator interface {
		GetUserName(*sync.WaitGroup)
		GetFriends(*sync.WaitGroup)
		Read() models.User
		Done()
	}

	UserData struct {
		Name    chan string
		Friends chan []models.User
	}
)

func NewUserData() UserData {
	return UserData{
		Name:    make(chan string, 1),
		Friends: make(chan []models.User, 1),
	}
}

func (u *UserData) GetUserName(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	u.Name <- getUserData().Name
}

func (u *UserData) GetFriends(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	u.Friends <- getUserFriends()
}

// Done allows you to specify when you are finished aggregating.
// It will close any channels which don't yet have data, ensuring
// subsequent Read does not block.
func (u *UserData) Done() {
	v := reflect.ValueOf(*u)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()
		possibleChan := reflect.ValueOf(field)
		var isChannel = possibleChan.Kind() == reflect.Chan
		if !isChannel {
			continue
		}

		if possibleChan.Len() == 0 {
			possibleChan.Close()
		}
	}
}

// Read aggregates all data that was prepared for consumption.
func (u UserData) Read() (user models.User) {
	// FIXME: This has two main caveats currently: the complexity scales with
	// the amount of data in the struct, and it only supports reading from
	// channels of one element
	for {
		select {
		case friends, ok := <-u.Friends:
			if !ok {
				u.Friends = nil
			} else {
				user.Friends = friends
				close(u.Friends)
			}
		case name, ok := <-u.Name:
			if !ok {
				u.Name = nil
			} else {
				user.Name = name
				close(u.Name)
			}
		default:
			continue
		}

		if u.Name == nil && u.Friends == nil {
			break
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
