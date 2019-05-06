package melee

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nats-io/go-nats"

	"github.com/lancer-kit/armory/api/httpx"
	"github.com/lancer-kit/melee/config"
	"github.com/pkg/errors"
)

func StatusNotOK(resp *http.Response) error {
	body := resp.Body
	b, _ := ioutil.ReadAll(body)
	return errors.Errorf(
		"Status code for request isn`t equal 200. Code: %d. Reason: %s",
		resp.StatusCode,
		string(b),
	)
}

func GetXClient() httpx.Client {
	netCfg := config.Config().Net.Get()

	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(netCfg.Timeout) * time.Second,
			KeepAlive: time.Duration(netCfg.KeepAlive) * time.Second,
			DualStack: true,
		}).DialContext,

		TLSHandshakeTimeout:   5 * time.Second,
		IdleConnTimeout:       time.Duration(netCfg.IdleConnTimeout) * time.Second,
		MaxIdleConns:          netCfg.MaxIdleConns,
		MaxConnsPerHost:       netCfg.MaxIdleConns,
		MaxIdleConnsPerHost:   netCfg.MaxConnsPerHost,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := httpx.NewXClient()
	client.Transport = tr
	client.Timeout = time.Duration(netCfg.Timeout) * time.Second

	return client

}

func wait4Signal(natsAddress string) bool {
	nc, err := nats.Connect(natsAddress)

	if err != nil {
		log.Fatal("Can not connect to nats: ", err)
	}

	//waiting start command from nats connection
	bus := make(chan *nats.Msg, 5)
	sub, err := nc.ChanSubscribe("highload", bus)
	if err != nil {
		log.Fatal(err)
	}
	msg := <-bus
	defer sub.Unsubscribe()

	return string(msg.Data) == "start"

}
