package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	limit           = 1000
	concurrencySize = 5
)

func producer(limit int) chan int {
	inputs := make(chan int)
	go func() {
		for i := 1; 1 <= limit; i++ {
			inputs <- i
		}
	}()
	return inputs
}

func processor(i int) (int, error) {
	if i == 10 {
		return 0, errors.New("i hate 5")
	}
	time.Sleep(5 * time.Second)
	return i * i, nil
}

func terminator(results chan int) {
	for item := range results {
		fmt.Println(item)
	}
}

func consumer(concurrencySize int, inputs chan int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errors := make(chan error)
	results := make(chan int)
	go terminator(results)
	for i := 0; i < concurrencySize; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case item := <-inputs:
					result, err := processor(item)
					if err != nil {
						errors <- err
					} else {
						results <- result
					}
				}
			}
		}()
	}
	err := <-errors
	cancel()
	return err
}

func main() {
	inputs := producer(limit)
	err := consumer(concurrencySize, inputs)
	fmt.Println(err)
}
