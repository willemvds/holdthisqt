package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"runtime"
	"encoding/json"

	"github.com/willemvds/holdthisqt"
)

const SLOTCOUNT = 8

type HttpLockRequest struct {
	Slot int
	Value string
}

func main() {
	runtime.GOMAXPROCS(2)

	var lockList = holdthisqt.NewLockList(SLOTCOUNT)
	log.Println(lockList)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("req read error:", err)
			return
		}
		var lr HttpLockRequest
		err = json.Unmarshal(bytes, &lr)
		if err != nil {
			log.Println("json error:", err)
			return
		}
		lf, err := lockList.GetLockFunc(lr.Slot)
		if err != nil {
			return
		}
		result := lf([]byte(lr.Value))
		if result == holdthisqt.RESPONSE_OK {
			res.Write([]byte("OK"))
		} else {
			res.Write([]byte("ER"))
		}
	})

	err := http.ListenAndServe(":40080", nil)
	log.Println(err)
}
