package holdthisqt

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

func getResultBytes(r lockResult) []byte {
	if r == RESPONSE_OK {
		return []byte("OK")
	}
	return []byte("NO")
}

type lockServer struct {
	listeners []net.Listener
	locklist  lockList
}

// TODO: more cowbell
func NewLockServer(listeners []net.Listener) lockServer {
	server := lockServer{}
	server.listeners = listeners
	listenerCount := len(listeners)
	server.locklist = NewLockList(listenerCount)

	for i := 0; i < listenerCount; i++ {
		listener := &server.listeners[i]
		lockfunc, err := server.locklist.GetLockFunc(i)
		if err != nil {
			panic("Not cool man")
		}
		go func() {
			for {
				conn, err := (*listener).Accept()
				if err != nil {
					print(err)
					return
				}
				go HandleConnection(conn, lockfunc)
			}
		}()
	}
	return server
}

func HandleConnection(conn net.Conn, lockfunc lockFunc) {
	var datalen int8
	for {
		lenbuf := make([]byte, 1, 1)
		_, err := io.ReadFull(conn, lenbuf)
		if err == io.EOF {
			return
		}
		lenreader := bytes.NewReader(lenbuf)
		err = binary.Read(lenreader, binary.LittleEndian, &datalen)
		buf := make([]byte, datalen, datalen)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			conn.Write([]byte("ER"))
			return
		}
		result := lockfunc(buf)
		conn.Write(getResultBytes(result))
	}
}
