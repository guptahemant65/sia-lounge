// main.go

package main

import (
	"os"
)

func main() {
	a := App{}
	// You need to set your Username and Password here
	a.Initialize("root", "abcd5520", "sia_lounge")

	a.Run(":" + os.Getenv("PORT"))
}
