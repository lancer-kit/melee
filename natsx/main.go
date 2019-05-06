package natsx

import (
	"time"

	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/melee/config"
	"github.com/nats-io/go-nats"
)

func GetConn() (*nats.Conn, error) {
	var natsConn *nats.Conn
	natsConn, err := nats.Connect(
		config.Config().NATS,
	)
	if err != nil {
		return nil, err
	}
	return natsConn, nil
}

func NatsListner(subj string, ch chan<- []byte) {
	connNats, err := GetConn()
	if err != nil {
		log.Default.Errorln(err)
	}

	defer connNats.Close()
	defer close(ch)

	var b []byte
	_, err = connNats.Subscribe(subj, func(msg *nats.Msg) {
		if msg != nil {
			b = msg.Data
		}
	})
	if err != nil {
		log.Default.Errorln(err)
	}

	timerChan := time.NewTimer(nats.DefaultDrainTimeout / 3).C

	for {
		select {
		case <-timerChan:
			return
		default:
			if b != nil {
				ch <- b
				return
			}
		}
	}
}
