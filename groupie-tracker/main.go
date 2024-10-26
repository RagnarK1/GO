package main

import (
	print2 "server/models"
	server "server/views"
)

func main() {
	print2.HandleJson()
	server.Main()
}
