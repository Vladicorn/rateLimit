package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type rateLmt struct {
	maxSameTime  int
	maxPerMinute int
}

func (rateLmt *rateLmt) ratelimit(wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()

	i := 1
	j := 1
	t := time.Now()
	timeMax := 60 * time.Second
	done := make(chan bool)

	for {
		if (time.Since(t) < timeMax) && (j <= rateLmt.maxPerMinute) {
			if i <= rateLmt.maxSameTime {

				c, ok := <-ch
				go SimplTask(done, c)
				if !ok {
					break
				}

				i++
				j++
			} else {

				select {
				case <-done:
					i--
					break
					return
				}
			}
		} else {
			select {
			case <-time.After(timeMax - time.Since(t)):
				j = 0
				t = time.Now()
				break
			}
		}
	}
}

//Выполнение простой задачи
func SimplTask(chanTask chan bool, g int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		r := rand.Intn(5)
		//r := 3
		time.Sleep(time.Duration(r) * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println(g, "Имитация бурной деятельности...  ", time.Now())
			chanTask <- true
			return
		}
	}

}
