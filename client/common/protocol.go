package common

import (
	"net"
	"bufio"

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
func (p *Protocol) WriteData(data string) error {
	totalSent := 0
	for totalSent < len(data) {
		n, err := p.conn.Write([]byte(data[totalSent:]))
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}

// ReadAll reads all data from the connection until EOF or error.
func (p *Protocol) ReadAll(ID string) (string, error) {
	var msg string
	readBuffer := bufio.NewReader(p.conn)
	for {
		line, err := readBuffer.ReadString('\n')
		if err != nil {
			log.Errorf("action: read_all | result: fail | client_id: %s | error: %v", ID, err)
			return "", err
		}
		msg += line
	}
	log.Infof("action: read_all | result: success | client_id: %s | message: %s", ID, msg)	
	return msg, nil
}

// SerializeBet serializes a bet map to JSON format.
// func SerializeBet(bet map[string]string) ([]byte, error) {
// 	return json.Marshal(bet)
// }

// // SendBet serializes and sends a bet over the connection.
// func (p *Protocol) SendBet(bet map[string]string) error {
// 	betData, err := SerializeBet(bet)
// 	if err != nil {
// 		log.Errorf("action: serialize_bet | result: fail | error: %v", err)
// 		return err
// 	}

// 	err = p.WriteData(append(betData, '\n'))
// 	if err != nil {
// 		log.Errorf("action: send_bet | result: fail | error: %v", err)
// 		return err
// 	}

// 	return nil
// }

// ReceiveResponse reads the response from the server.
// func (p *Protocol) ReceiveResponse() ([]byte, error) {
// 	return p.ReadAll()
// }
