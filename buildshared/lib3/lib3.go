package lib3

import (
	"log"

	"github.com/yunabe/golang-codelab/buildshared/lib2"
)

func init() {
	log.Println("Package lib3 is initialized")
	log.Printf("From lib3: lib2.IntSum(3, 4, 5) == %d", lib2.IntSum(3, 4, 5))
}
