package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/SergJeepee/go-minions/minions"
)

func main() {
	f := func(s string) string {
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		return fmt.Sprintf("Banana? %v!", s)
	}

	instructedMinions, input, output := minions.ListenHere[string, string](f, 5, 10, 10)
	done := instructedMinions.Go()

	go func() {
		for i := 0; i < 10; i++ {
			input <- "Banana " + strconv.Itoa(i+1)
		}
		close(input)
	}()

	for res := range output {
		log.Println(res)
	}

	<-done
}
