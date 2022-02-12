package test

//go test
import (
	"fmt"
	"rateLimit/ratelimiter"
	"sync"
	"testing"
	"time"
)

func TestRate(t *testing.T) {
	var wg sync.WaitGroup

	//инициализация
	rate := ratelimiter.RateLmt{
		MaxSameTime:  7,  //максимальное количество одновременных задач
		MaxPerMinute: 15, //максимальное количество задач в течении минуты
	}

	wg.Add(1)
	ch := make(chan int)
	go rate.Ratelimit(&wg, ch)

	for i := 1; i <= 40; i++ {
		ch <- i
	}

	close(ch)
	wg.Wait()
	fmt.Println("END", time.Now())
}
