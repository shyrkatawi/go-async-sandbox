package to_refactor

import (
	"errors"
	"fmt"
	"sync"
)

/*
explain and refactor if needed :D

type Stack struct {
	mutex sync.Mutex
	data  []string
}

func NewStack() Stack {
	return Stack{}
}

func (b Stack) Push(value string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.data = append(b.data, value)
}

func (b Stack) Pop() {
	if len(b.data) < 0 {
		panic("pop: stack is empty")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.data = b.data[:len(b.data)-1]
}

func (b Stack) Top() string {
	if len(b.data) < 0 {
		panic("top: stack is empty")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.data[len(b.data)-1]
}

var stack Stack

func producer() {
	for i := 0; i < 1000; i++ {
		stack.Push("message")
	}
}

func consumer() {
	for i := 0; i < 10; i++ {
		_ = stack.Top()
		stack.Pop()
	}
}

func main() {
	producer()

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			consumer()
		}()
	}

	wg.Wait()
}
*/

type myStack struct {
	mutex sync.RWMutex
	data  []string
}

func newStack() *myStack {
	return &myStack{}
}

func (stack *myStack) Push(value string) {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	stack.data = append(stack.data, value)
}

func (stack *myStack) Pop() (string, error) {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	if len(stack.data) == 0 {
		return "", errors.New("pop: stack is empty")
	}

	index := len(stack.data) - 1
	v := stack.data[index]
	stack.data = stack.data[:index]

	return v, nil
}

func (stack *myStack) Top() (string, error) {
	stack.mutex.RLock()
	defer stack.mutex.RUnlock()

	if len(stack.data) == 0 {
		return "", errors.New("top: stack is empty")
	}
	return stack.data[len(stack.data)-1], nil
}

func main() {
	stack := newStack()

	producer := func() {
		for i := 0; i < 1000; i++ {
			stack.Push(fmt.Sprintf("message %d", i))
		}
	}

	consumer := func() {
		for i := 0; i < 10; i++ {
			go func() {
				topV, err := stack.Top()
				if err != nil {
					fmt.Println("Top err:", err)
				} else {
					fmt.Println("Top v:", topV)
				}
			}()
			go func() {
				popV, err := stack.Pop()
				if err != nil {
					fmt.Println("Pop err:", err)
				} else {
					fmt.Println("Pop v:", popV)
				}
			}()
		}
	}

	producer()

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer()
		}()
	}

	wg.Wait()
	fmt.Println("done")
	fmt.Println(len(stack.data))
}

func ToRefactor1() {
	main()
}
