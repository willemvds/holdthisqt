package main

import (
	"log"
	"fmt"
	"net"
	"os"
	"runtime"
	"os/signal"

	"github.com/willemvds/holdthisqt"
)

const TCPCOUNT = 10
const UNIXCOUNT = 10

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listeners := make([]net.Listener, 0, TCPCOUNT + UNIXCOUNT)

	for i := 0; i < TCPCOUNT; i++ {
		tcpl, err := net.Listen("tcp", fmt.Sprintf(":420%02d", i))
		if err != nil {
			log.Fatalf("tcp listen error", err)
		}
		listeners = append(listeners, tcpl)
	}

	for i := 0; i < UNIXCOUNT; i++ {
		sockpath := fmt.Sprintf("/tmp/lock%02d.sock", i)
		os.Remove(sockpath)
		unixl, err := net.Listen("unix", sockpath)
		if err != nil {
			log.Fatalf("unix listen error", err)
		}
		listeners = append(listeners, unixl)
	}

	server := holdthisqt.NewLockServer(listeners)
	fmt.Println(server)
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
