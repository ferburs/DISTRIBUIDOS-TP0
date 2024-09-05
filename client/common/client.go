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
	BatchMaxSize     int
}

type Client struct {
	config  ClientConfig
	conn    net.Conn
	done    chan bool
	protocol *Protocol
	reader *csv.Reader
}

func NewClient(config ClientConfig, reader *csv.Reader) *Client {
	client := &Client{
		config: config,
		done:   make(chan bool, 1),
		reader: reader,
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

	maxTamanio = 8 * 1024
	endofFile := false

	for !endofFile {
		c.createClientSocket()
		batch = ""
		for i := 0; i <= c.config.BathMaxmount; i++ {
			lineRead, err := c.reader.Read()
			if lineRead == nil {
				endofFile = true
				break
			}
			betMessage := createMessage(lineRead, c.config.ID)
			messageSerialized := betMessage.Serialize()
			}
		batch += "\n"

		err := c.protocol.WriteData(batch)
		if err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
		} else {
			log.Infof("action: send_batch | result: success | client_id: %v | batch_size: %v", c.config.ID, len(batch) / maxTamanio)
		time.Sleep(c.config.LoopPeriod)	
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	}
}	 

func createMessage(line []string, agency string) map[string]string {
	return NewMessage(agency, line[0], line[1], line[2], line[3], line[4])
}

func verifySize(message string, maxTamanio int) bool {
	return len(message) <= maxTamanio
}