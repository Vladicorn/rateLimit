package test

import (
	"fmt"
	"sync"
	"testing"
	"https://github.com/Vladicorn/rateLimit/blob/main/ratelimiter/ratelimiter.go"
	
)

func TestRate(t *testing.T) {
	var wg sync.WaitGroup
	//инициализация

	rate := ratelimiter.rateLmt{
		maxSameTime:  5,  //максимальное количество одновременных задач
		maxPerMinute: 30, //максимальное количество задач в течении минуты
	}

	wg.Add(1)
	ch := make(chan int)
	go rate.ratelimiterratelimit(&wg, ch)

	for i := 1; i < 100; i++ {
		ch <- i
	}

	close(ch)
	wg.Wait()
	fmt.Println("END")
}
