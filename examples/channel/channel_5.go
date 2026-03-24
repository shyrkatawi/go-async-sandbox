package channel

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
implement processV1 fn which accepts slice of user ids, must return slice with users
if any error occurs in fetchUser, other "parallel" fetch functions must abort

type user struct {
	id   uint
	name string
}

var users = map[uint]*user{
	1: {id: 1, name: "user_1"},
	2: {id: 2, name: "user_2"},
	3: {id: 3, name: "user_3"},
}

func fetchUser(id uint) (*user, error) {
	time.Sleep(time.Second)
	usr, ok := users[id]
	if !ok {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return usr, nil
}

func processV1(ids []string) ([]user, error){}

*/

type user struct {
	id   uint
	name string
}

var usersStore = map[uint]*user{
	1: {id: 1, name: "user_1"},
	2: {id: 2, name: "user_2"},
	3: {id: 3, name: "user_3"},
}

func fetchUser(ctx context.Context, id uint) (*user, error) {
	type userResponse struct {
		user *user
		err  error
	}

	userResponseCh := make(chan *userResponse)
	emulateRequest := func() {
		time.Sleep(time.Duration(rand.Intn(10)) * 100 * time.Millisecond)
		usr, ok := usersStore[id]
		if !ok {
			userResponseCh <- &userResponse{err: fmt.Errorf("user with id %d not found", id)}
		} else {
			userResponseCh <- &userResponse{user: usr}
		}
	}

	go emulateRequest()

	select {
	case <-ctx.Done():
		fmt.Println("fetch was canceled")
		return nil, ctx.Err()
	case usrResp := <-userResponseCh:
		if usrResp.err != nil {
			return nil, usrResp.err
		}
		return usrResp.user, nil
	}
}

func processV1(ctx context.Context, ids []uint) ([]*user, error) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	users := make([]*user, len(ids))
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	var returnErr error

	for i, id := range ids {
		wg.Add(1)

		go func() {
			defer wg.Done()
			usr, err := fetchUser(ctxWithCancel, id)
			if err != nil {
				mu.Lock()
				if returnErr == nil {
					returnErr = err
					cancel()
				}
				mu.Unlock()
			} else {
				users[i] = usr
			}
		}()
	}

	wg.Wait()

	if returnErr != nil {
		return nil, returnErr
	}
	return users, nil
}

func processV2(ctx context.Context, ids []uint) ([]*user, error) {
	errGrp, errGrpCxt := errgroup.WithContext(ctx)

	users := make([]*user, len(ids))

	for i, id := range ids {
		errGrp.Go(
			func() error {
				usr, err := fetchUser(errGrpCxt, id)
				if err != nil {
					return err
				}
				users[i] = usr
				return nil
			},
		)
	}

	err := errGrp.Wait()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func RunChannel5() {
	showResult := func(tag string, fn func() ([]*user, error)) {
		time.Sleep(time.Second)
		users, err := fn()
		fmt.Println(tag)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Success: %+v\n", users)
		}
		fmt.Println("———————————————————")
	}

	//showResult("withoutError v1", func() ([]*user, error) {
	//	return processV1(context.Background(), []uint{1, 2, 3})
	//})
	//showResult("withoutError v2", func() ([]*user, error) {
	//	return processV2(context.Background(), []uint{1, 2, 3})
	//})
	//
	//showResult("withTimeout v1", func() ([]*user, error) {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	go func() {
	//		time.Sleep(100 * time.Millisecond)
	//		cancel()
	//	}()
	//	return processV1(ctx, []uint{1, 2, 3})
	//})
	//showResult("withTimeout v2", func() ([]*user, error) {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	go func() {
	//		time.Sleep(100 * time.Millisecond)
	//		cancel()
	//	}()
	//	return processV2(ctx, []uint{1, 2, 3})
	//})

	showResult("withRequestError v1", func() ([]*user, error) {
		return processV1(context.Background(), []uint{1, 2, 3, 4, 5, 6, 7})
	})
	showResult("withRequestError v2", func() ([]*user, error) {
		return processV2(context.Background(), []uint{1, 2, 3, 4, 5, 6, 7})
	})
}
