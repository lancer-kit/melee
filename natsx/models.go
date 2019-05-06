package natsx

import (
	"encoding/json"
	"errors"
)

type Data struct {
	Addressee string `json:"addressee"`
	Subject   string `json:"subject"`
	Text      string `json:"text"`
}

type Listener struct {
	Ch chan []byte
}

// WARNING: does not work with goroutines and load test
func (l *Listener) ParseNATS(i interface{}) error {
	data, ok := <-l.Ch
	if !ok {
		return errors.New("chanel close")
	}
	if data == nil {
		return errors.New("close chanel `listen timeout` ")
	}
	err := json.Unmarshal(data, i)
	if err != nil {
		return err
	}
	return nil
}

func NewListener(topic string) *Listener {
	listener := new(Listener)
	listener.Ch = make(chan []byte)
	go NatsListner(topic, listener.Ch)
	return listener
}
