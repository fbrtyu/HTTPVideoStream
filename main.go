package main

import (
	"HTTPVideoStream/packages/webserver"
	"fmt"
)

func main() {
	fmt.Println("Server starting...")
	webserver.Webserver()
}
