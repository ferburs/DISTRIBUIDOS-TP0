package common

import (
	"net"
	"bufio"
	"fmt"
	//"io"

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
		log.Infof("La linea es %s", line)
		// if err == io.EOF && line != "" {
		// 	print("entro al end of file")
		// 	msg += line
		// 	return msg, err
		// }
		if err != nil {
			log.Errorf("action: read_all | result: fail | client_id: %s | error: %v", ID, err)
			return "", err
		}
		msg += line
		break
	}
	log.Infof("action: read_all | result: success | client_id: %s | message: %s", ID, msg)	
	return msg, nil
}


func (p *Protocol) NotifyDone(ID string) error {
	msg := fmt.Sprintf("%v#NOTIFY_DONE\n\n", ID)
	err := p.WriteData(msg)
	if err != nil {
		return err
	}
	return nil
}

func (p *Protocol) RequestWinners() error {
	err := p.WriteData("REQUEST_WINNERS\n\n")
	if err != nil {
		return err
	}
	return nil
}