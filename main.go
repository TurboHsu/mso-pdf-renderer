package main

import (
	"flag"
	"github.com/TurboHsu/mso-pdf-renderer/manager"
	"github.com/TurboHsu/mso-pdf-renderer/server"
)

func main() {
	// Read address from flags
	var addr string
	var lifeCycle, lifeCycleCheckInterval int64
	flag.StringVar(&addr, "addr", ":8080", "address to listen on (Default value is :8080)")
	flag.Int64Var(&lifeCycle, "lifecycle", 300, "life cycle of cache in seconds (Default value is 300 seconds.)")
	flag.Int64Var(&lifeCycleCheckInterval, "lifecycle-check-interval", 60, "life cycle check interval in seconds (Default value is 60 seconds.)")
	flag.Parse()

	// Start lifecycle manager
	go manager.CacheLifeCycleRoutine(lifeCycle, lifeCycleCheckInterval)
	server.ListenAndServe(addr)
}
