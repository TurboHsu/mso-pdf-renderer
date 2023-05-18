package main

import (
	"mso-pdf-renderer/server"
)

func main() {
	server.ListenAndServe(":8080")
}
