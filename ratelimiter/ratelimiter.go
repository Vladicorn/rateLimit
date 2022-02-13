package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type RateLmt struct {
	MaxSameTime  int //максимальное количество одновременных задач
	MaxPerMinute int //максимальное количество задач в течении минуты
}

//буфферное значение канала с флагом
type RsTrig struct {
	S int
	R bool
}

func (rateLmt *RateLmt) Ratelimit(wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()
	var wg1 sync.WaitGroup

	i := 0         //счетчик текущих задач
	j := 0         // счетчик в минуту
	countTask := 0 //счетчик неотработаных
	t := time.Now()
	flagOk := false    //Флаг, триггер для послед канала
	var bufChan RsTrig //буфферное значение канала
	bufChan.R = false
	timeMax := 60 * time.Second
	done := make(chan bool) //канал о выполнении задачи

exit:
	for {
		if (time.Since(t) < timeMax) && (j < rateLmt.MaxPerMinute) {
			if i < rateLmt.MaxSameTime {

				if bufChan.R {
					wg1.Add(1)
					go SimplTask(&wg1, done, bufChan.S) //выполнение задачи
					countTask++
					i++
					j++
					bufChan.R = false

				} else {
					c, ok := <-ch
					if !ok {
						break exit
					}
					wg1.Add(1)
					go SimplTask(&wg1, done, c) //выполнение задачи
					countTask++
					i++
					j++
				}

			} else {
				//когда максимальное одновременных каналов
				<-done
				countTask--
				i--
			}
			flagOk = true

		} else {
			if flagOk {
				c, ok := <-ch
				if !ok {
					break exit
				} else {
					bufChan.S = c
					bufChan.R = true
					flagOk = false
				}
			}

			timer := time.NewTimer(timeMax - time.Since(t))
			select {
			case <-done:
				countTask--
				i--
				break
			case <-timer.C:
				j = 0
				t = time.Now()
			}

		}
	}

	for {
		if countTask == 0 {
			break
		} else {
			<-done
			countTask--
		}
	}
	wg1.Wait()
}

//Выполнение простой задачи
func SimplTask(wg1 *sync.WaitGroup, chanTask chan bool, g int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	go func() {
		//r := rand.Intn(5)
		r := 4
		time.Sleep(time.Duration(r) * time.Second)
		done <- true
	}()

	for {
		select {
		case <-done:
			chanTask <- true
			wg1.Done()
			fmt.Println(g, "Имитация бурной деятельности...  ", time.Now())
			return
		}
	}

}
