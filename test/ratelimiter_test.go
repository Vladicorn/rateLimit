package test

import (
	"fmt"
	"sync"
	"testing"
)

func TestRate(t *testing.T) {
	var wg sync.WaitGroup
	//инициализация

	rate := Ratelimiter.rateLmt{
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
