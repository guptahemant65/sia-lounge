// main.go

package main

import (
	"os"
)

func main() {
	a := App{}
	// You need to set your Username and Password here
	a.Initialize("b08c704fbb334f", "0a40bb43", "heroku_18e8f056bc19e0a")

	a.Run(":" + os.Getenv("PORT"))
}
