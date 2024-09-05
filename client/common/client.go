package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"encoding/csv"
	"fmt"
	"strings"
	//"io"

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
	c.createClientSocket()

	for !endofFile {

		batch := ""
		for i := 0; i < c.config.BatchMaxSize; i++ {
			lineRead, err := c.reader.Read()
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
	}
	c.notifyClientDone(c.config.ID)
	c.waitWinners(c.config.ID)
}	 

func (c *Client) notifyClientDone(ID string){
	err := c.protocol.NotifyDone(ID)
	if err != nil {
		log.Errorf("action: notify_done | result: fail | client_id: %v | error: %v",ID, err)
	} 

}

func (c *Client) waitWinners(ID string){

	println("cliet envia esperando ganadores")

	msg := fmt.Sprintf("%v#REQUEST_WINNERS\n\n", ID)
	err := c.protocol.WriteData(msg)

	if err != nil {
		log.Infof("action: request_winners | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.conn.Close()
		time.Sleep(6 * time.Second)
	}
	log.Infof("action: request_winners | result: success | client_id: %v", c.config.ID)

	response, err := c.protocol.ReadAll(c.config.ID)
	if err != nil {
		log.Errorf("action: read_all | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.conn.Close()
		time.Sleep(6 * time.Second)
	}

	//Saco los \n del final
	winners := response[:len(response)-2]
	//println(winners)
	var totalWinners int
	if len(winners) == 0 {
		totalWinners = 0
	} else {
		totalWinners = len(strings.Split(winners, "\n"))
	}
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", totalWinners)
}