package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
)

const TCPCOUNT = 10
const UNIXCOUNT = 10

var lockValues = [][]byte{
	[]byte("foo.bar"),
	[]byte("bar.one"),
	[]byte("one.time"),
	[]byte("time.time"),
	[]byte("time.machine"),
	[]byte("machine.man"),
	[]byte("man.poo"),
	[]byte("poo.party"),
}

var errNo = errors.New("NO NO NO... wait... NO")
var errEr = errors.New("Remote Error")
var errWhat = errors.New("Eh wha?")

func getLock(conn net.Conn, value []byte) error {
	binary.Write(conn, binary.LittleEndian, int8(len(value)))
	conn.Write(value)
	buf := make([]byte, 2, 2)
	_, err := io.ReadFull(conn, buf)
	if err == nil {
		if bytes.Equal(buf, []byte("OK")) {
			return nil
		} else if bytes.Equal(buf, []byte("NO")) {
			return errNo
		} else if bytes.Equal(buf, []byte("ER")) {
			return errEr
		}
		return errWhat
	}
	return err
}

func worker(id string, conn net.Conn) {
	okCount := 0
	noCount := 0
	erCount := 0
	otherCount := 0
	for {
		err := getLock(conn, lockValues[rand.Intn(len(lockValues))])
		if err != nil {
			if err == errNo {
				noCount++
			} else if err == errEr {
				erCount++
			} else {
				otherCount++
			}
		} else {
			okCount++
		}
		if (okCount+noCount+erCount+otherCount)%50000 == 0 {
			atomic.AddUint64(&total, 50000)
			fmt.Printf("%s: %d requests, %d OK, %d NO, %d ER, %d Other\n", id, okCount+noCount+erCount+otherCount, okCount, noCount, erCount, otherCount)
		}
	}
}

var total uint64

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 0; i < UNIXCOUNT; i++ {
		conn, err := net.Dial("unix", fmt.Sprintf("/tmp/lock%02d.sock", i))
		if err != nil {
			log.Fatalf(err.Error())
		}
		go worker(fmt.Sprintf("Unix Worker %d", i), conn)
	}

	for i := 0; i < TCPCOUNT; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:420%02d", i))
		if err != nil {
			log.Fatalf(err.Error())
		}
		go worker(fmt.Sprintf("TCP Worker %d", i), conn)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Printf("Total Requests: %d", total)
}
