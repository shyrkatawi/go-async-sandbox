package sync_cond

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var pokemonList = [5]string{"Pikachu", "Charmander", "Squirtle", "Bulbasaur", "Jigglypuff"}
var currentPokemon = ""

var waitDuration = 11 * time.Millisecond

func getRandomPokemon() string {
	return pokemonList[rand.Intn(len(pokemonList))]
}

func waitWithCh() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	stopProducingCh := make(chan struct{})
	pokemonCn := make(chan string)

	pokemonProducer := func() {
		defer wg.Done()
		defer close(pokemonCn)

		for i := 0; i < 5; i++ {
			time.Sleep(waitDuration)
			newPokemon := getRandomPokemon()
			fmt.Println("newPokemon", newPokemon)

			select {
			case <-stopProducingCh:
				return
			case pokemonCn <- newPokemon:
			}
		}
	}

	pokemonConsumer := func() {
		defer wg.Done()

		for pokemon := range pokemonCn {
			fmt.Println("get", pokemon)
			if pokemon == pokemonList[0] {
				fmt.Println("finally")
				close(stopProducingCh)
				return
			}
		}
	}

	go pokemonProducer()
	go pokemonConsumer()

	wg.Wait()
}

func waitWithSyncCond() {
	pokemonsToProduceNumber := 10
	pokemonsToWaitNumber := 5

	wg := sync.WaitGroup{}
	wg.Add(pokemonsToWaitNumber + 1)

	isProducing := true

	cond := sync.NewCond(&sync.Mutex{})

	pokemonProducer := func() {
		defer wg.Done()

		for i := 0; i < pokemonsToProduceNumber; i++ {
			time.Sleep(waitDuration)

			cond.L.Lock()
			newPokemon := getRandomPokemon()
			fmt.Println("newPokemon", newPokemon)
			currentPokemon = newPokemon
			cond.L.Unlock()

			cond.Signal()
		}

		isProducing = false

		cond.Broadcast()
	}

	pokemonConsumer := func(consumerNumber int, pokemonToWait string) {
		defer wg.Done()

		cond.L.Lock()
		defer cond.L.Unlock()

		for currentPokemon != pokemonToWait {
			if isProducing == false {
				fmt.Println(consumerNumber, "Stop waiting for", pokemonToWait)
				return
			}
			fmt.Println(consumerNumber, "get", currentPokemon)
			fmt.Println(consumerNumber, "waiting more...")
			cond.Wait()
		}

		fmt.Println("FINALLY GOT", currentPokemon)
	}

	go pokemonProducer()
	for i := 0; i < pokemonsToWaitNumber; i++ {
		go pokemonConsumer(i, getRandomPokemon())
	}

	wg.Wait()
}

func RunSyncCond() {
	//waitWithSyncCond()
	waitWithCh()
}
