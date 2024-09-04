package common

import (
	//"encoding/json"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"github.com/op/go-logging"
)

//var log = logging.MustGetLogger("log")

type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

type Client struct {
	config  ClientConfig
	conn    net.Conn
	done    chan bool
	protocol *Protocol
}

func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		done:   make(chan bool, 1),
	}
	return client
}

func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	c.conn = conn
	c.protocol = NewProtocol(conn)
	return nil
}

func (c *Client) StartClientLoop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Infof("action: graceful_shutdown | result: success | client_id: %v | signal: %v", c.config.ID, sig)
		if c.conn != nil {
			c.conn.Close()
			log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
		}
		c.done <- true
	}()

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		err := c.createClientSocket()
		if err != nil {
			return
		}

		bet := map[string]string{
			"agency":        c.config.ID,
			"NOMBRE":    os.Getenv("NOMBRE"),
			"APELLIDO":   os.Getenv("APELLIDO"),
			"DOCUMENTO": os.Getenv("DOCUMENTO"),
			"NACIMIENTO": os.Getenv("NACIMIENTO"),
			"NUMERO":   os.Getenv("NUMERO"),
		}

		err = c.protocol.SendBet(bet)
		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		_, err = c.protocol.ReceiveResponse()
		if err != nil {
			log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		c.conn.Close()

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet["DOCUMENTO"], bet["NUMERO"])

		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
