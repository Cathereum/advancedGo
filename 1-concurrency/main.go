package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	fmt.Println(getSq())

}

func getSq() []int {
	ch := make(chan []int)
	var wg sync.WaitGroup
	res := []int{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		numbers := make([]int, 10)
		for i := range numbers {
			numbers[i] = rand.Intn(100) + 1
		}
		ch <- numbers
		close(ch)
	}()

	go func(ch chan []int) {
		defer wg.Done()
		for nums := range ch {
			for _, n := range nums {
				res = append(res, n*n)
			}
		}
	}(ch)

	wg.Wait()
	return res

}
