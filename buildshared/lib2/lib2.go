package lib2

import (
	"log"
)

func init() {
	log.Println("Package lib2 is initialized")
}

func IntSum(args ...int) (sum int) {
	for i := 0; i < len(args); i++ {
		sum += args[i]
	}
	return
}
