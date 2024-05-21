package main

import (
	"fmt"
	"time"
)

func main() {

	// Using for loop
	start := time.Now()
	sum := 0
	for i := 0; i < 1000000000; i++ {
		sum += i
	}
	fmt.Println("For loop time:", time.Since(start))

	// Using goto
	start = time.Now()
	sum = 0
	i := 0
loop:
	if i < 1000000000 {
		sum += i
		i++
		goto loop
	}
	fmt.Println("Goto time:", time.Since(start))

	go func() {
		// Using goto
		start := time.Now()
		sum := 0
		i := 0
	loop:
		if i < 1000000000 {
			sum += i
			i++
			goto loop
		}
		fmt.Println("Goto time gorutine:", time.Since(start))
	}()
	time.Sleep(1 * time.Second)
	go func() {

		// Using for loop
		start := time.Now()
		sum := 0
		for i := 0; i < 1000000000; i++ {
			sum += i
		}
		fmt.Println("For loop time gorutine:", time.Since(start))

	}()

	time.Sleep(1 * time.Second)
}
