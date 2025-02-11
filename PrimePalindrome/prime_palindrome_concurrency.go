package main

import (
	"fmt"
	"sync"
)

func main() {
	N := 10
	sum := FindPrimePalindromes(N)
	fmt.Println(sum)
}

func FindPrimePalindromes(N int) int {
	var wg sync.WaitGroup
	ch := make(chan int, N)
	sumCh := make(chan int)
	var count, sum int
	var mu sync.Mutex

	go func() {
		for p := range ch {
			mu.Lock()
			sum += p
			count++
			mu.Unlock()
			if count == N {
				close(ch)
			}
		}
		sumCh <- sum
	}()

	num := 2
	for {
		if count >= N {
			break
		}
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if isPrime(n) && isPalindrome(n) {
				ch <- n
			}
		}(num)
		num++
	}

	wg.Wait()
	close(sumCh)
	return <-sumCh
}

func isPrime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func isPalindrome(num int) bool {
	original, reversed := num, 0
	for num > 0 {
		reversed = reversed*10 + num%10
		num /= 10
	}
	return original == reversed
}