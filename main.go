package main

import (
	"fmt"
	"rateLimit/ratelimiter"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	//инициализация
	rate := ratelimiter.RateLmt{
		MaxSameTime:  5,  //максимальное количество одновременных задач
		MaxPerMinute: 30, //максимальное количество задач в течении минуты
	}

	wg.Add(1)
	ch := make(chan int)
	go rate.Ratelimit(&wg, ch)

	for i := 1; i < 100; i++ {
		ch <- i
	}

	close(ch)
	wg.Wait()
	fmt.Println("END")
}
