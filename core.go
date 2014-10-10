package holdthisqt

import (
	"bytes"
	"errors"
)

// A lock value is simply a byte sequence
type lock []byte

type lockResult int8

const (
	RESPONSE_OK lockResult = iota
	RESPONSE_NO
)

type lockRequest struct {
	slot       int
	value      []byte
	resultChan chan lockResult
}

type lockList struct {
	slots       []lock
	requestChan chan lockRequest
}

func NewLockList(slotCount int) lockList {
	locklist := lockList{}
	locklist.slots = make([]lock, slotCount)
	locklist.requestChan = make(chan lockRequest)
	var request lockRequest
	go func() {
		for {
			request = <-locklist.requestChan
			request.resultChan <- locklist.lock(request)
		}
	}()

	return locklist
}

func (locklist lockList) lock(req lockRequest) lockResult {
	for i := range locklist.slots {
		if bytes.Equal(req.value, locklist.slots[i]) {
			locklist.slots[req.slot] = nil
			return RESPONSE_NO
		}
	}

	locklist.slots[req.slot] = req.value
	return RESPONSE_OK
}

// just having some fun i guess
type lockFunc func([]byte) lockResult

func (locklist lockList) GetLockFunc(slot int) (lockFunc, error) {
	if slot > len(locklist.slots) {
		return nil, errors.New("Fail")
	}
	return func(value []byte) lockResult {
		request := lockRequest{slot: slot, value: value}
		request.resultChan = make(chan lockResult)
		locklist.requestChan <- request
		return <-request.resultChan
	}, nil
}
