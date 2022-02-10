package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type rateLmt struct {
	maxSameTime  int
	maxPerMinute int
}

func main() {
	var wg sync.WaitGroup
	//инициализация
	rate := rateLmt{
		maxSameTime:  10, //максимальное количество одновременных задач
		maxPerMinute: 50, //максимальное количество задач в течении минуты
	}

	wg.Add(1)
	ch := make(chan int)
	go rate.ratelimit(&wg, ch)

	for i := 1; i < 100; i++ {
		ch <- i
	}

	close(ch)
	wg.Wait()
	fmt.Println("END")
}

func (rateLmt *rateLmt) ratelimit(wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()

	i := 1
	j := 1
	t := time.Now()
	timeMax := 60 * time.Second

	for {
		if (time.Since(t) < timeMax) && (j <= rateLmt.maxPerMinute) {
			if i <= rateLmt.maxSameTime {
				c, ok := <-ch
				go SimplTask(c)
				if !ok {
					break
				}

				i++
				j++
			} else {
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()
				done := make(chan bool)

				go func() {
					time.Sleep(1 * time.Second)
					done <- true
				}()

				for {
					select {
					case <-done:
						i = 0
						break
						return
					}
					if i == 0 {
						break
					}
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
func SimplTask(g int) {
	fmt.Println(g, time.Now())
}
