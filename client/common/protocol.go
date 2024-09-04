package common

import (
	"encoding/json"
	"net"
	//"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// Protocol handles communication over a network connection.
type Protocol struct {
	conn net.Conn
}

// NewProtocol creates a new Protocol instance with the provided connection.
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		conn: conn,
	}
}

// WriteData sends data over the connection, ensuring that all data is written.
func (p *Protocol) WriteData(data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		n, err := p.conn.Write(data[totalSent:])
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}

// ReadAll reads all data from the connection until EOF or error.
func (p *Protocol) ReadAll() ([]byte, error) {
	var data []byte
	for {
		buf := make([]byte, 1024)
		n, err := p.conn.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
	}
	return data, nil
}

// SerializeBet serializes a bet map to JSON format.
func SerializeBet(bet map[string]string) ([]byte, error) {
	return json.Marshal(bet)
}

// SendBet serializes and sends a bet over the connection.
func (p *Protocol) SendBet(bet map[string]string) error {
	betData, err := SerializeBet(bet)
	if err != nil {
		log.Errorf("action: serialize_bet | result: fail | error: %v", err)
		return err
	}

	err = p.WriteData(append(betData, '\n'))
	if err != nil {
		log.Errorf("action: send_bet | result: fail | error: %v", err)
		return err
	}

	return nil
}

// ReceiveResponse reads the response from the server.
func (p *Protocol) ReceiveResponse() ([]byte, error) {
	return p.ReadAll()
}
