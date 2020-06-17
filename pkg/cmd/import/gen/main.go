package main

import (
	"fmt"
	"math/rand"
)

var areas = []string{"shanghai", "beijing", "chendu", "suzhou", "riben", "feizhou"}

func main() {
	for i := 0; i < 1000; i++ {
		fmt.Printf("%v, %v, \"%s\", \"[images/1.jpg]\"\n", i%2, rand.Intn(100), areas[rand.Intn(5)])
	}
}
