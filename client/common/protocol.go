package common

import (
	"net"
	"bufio"
	"io"
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
    reader := bufio.NewReader(p.conn)

    for {
        part, err := reader.ReadString('\n')
        msg += part

        if err == io.EOF {
            return "", err
        }

        if err != nil {
            log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", ID, err)
            return "", err
        }

        // Verificar si los últimos dos caracteres son '\n\n'
        if len(msg) >= 2 && msg[len(msg)-2:] == "\n\n" {
            break
        }
    }

    return msg, nil
}