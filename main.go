package main

import (
	"flag"
	"mso-pdf-renderer/server"
)

func main() {
	// Read address from flags
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "address to listen on")
	flag.Parse()
	server.ListenAndServe(addr)
}
