package common

import (
	"bufio"
	"encoding/json"
	//"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	done   chan bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		done:   make(chan bool, 1),
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

// ReadAll reads all data from the connection until EOF or error
func (c *Client) readAll() ([]byte, error) {
	var data []byte
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		data = append(data, line...)
	}
	return data, nil
}

// StartClientLoop Send messages to the client until some time threshold is met
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
		// Crear el socket de conexión con el servidor
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

		// Serializar la apuesta a formato JSON
		betData, err := json.Marshal(bet)
		if err != nil {
			log.Errorf("action: serialize_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		_, err = c.conn.Write(append(betData, '\n'))
		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		// Leer la respuesta del servidor
		_, err = c.readAll()
		if err != nil {
			log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		c.conn.Close()

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			bet["DOCUMENTO"],
			bet["NUMERO"],
		)

		// Esperar un tiempo antes de enviar el siguiente mensaje
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
