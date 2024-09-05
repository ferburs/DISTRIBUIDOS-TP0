package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"encoding/csv"

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

	//maxTamanio := 8 * 1024
	kbsize := 1024
	endofFile := false

	for !endofFile {
		c.createClientSocket()
		batch := ""
		for i := 0; i < c.config.BatchMaxSize; i++ {
			lineRead, err := c.reader.Read()
			//log.Infof("LA LINEA ES line: %v", lineRead)
			//log.Infof("nombre: %v", lineRead[0])
			//log.Infof("apellido: %v", lineRead[1])
			if lineRead == nil {
				endofFile = true
				break
			}
			if err != nil {
				log.Errorf("action: read_line | result: fail | client_id: %v | error: %v", c.config.ID, err)
				endofFile = true
				break
			}

			betMessage := NewMessage(lineRead, c.config.ID)
			messageSerialized := betMessage.Serialize()
			//log.Infof("mensaje serializadooooooooooo: %v", messageSerialized)
			batch += messageSerialized
			}
		batch += "\n"
		//log.Infof("batch: %v", batch)

		err := c.protocol.WriteData(batch)
		if err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
		} else {
			log.Infof("action: send_batch | result: success | client_id: %v | batch_size: %v", c.config.ID, float64(len(batch)) / float64(kbsize))
		time.Sleep(c.config.LoopPeriod)	
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	}
}	 

// func createMessage(line []string, agency string) map[string]string {
// 	return NewMessage(agency, line[0], line[1], line[2], line[3], line[4])
// }

func verifySize(message string, maxTamanio int) bool {
	return len(message) <= maxTamanio
}