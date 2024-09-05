package common

import (
	"fmt"
)

type Message struct {
	AGENCIA string
	NOMBRE string
	APELLIDO string
	DOCUMENTO string
	NACIMIENTO string
	NUMERO string
}

func NewMessage(args []string, agency string) *Message {

	message := &Message{
		AGENCIA: agency,
		NOMBRE: args[0],
		APELLIDO:args[1],
		DOCUMENTO: args[2],
		NACIMIENTO: args[3],
		NUMERO: args[4],
	}
	return message
}

func (message *Message) Serialize() string{
	return fmt.Sprintf(
		"%s#%s#%s#%s#%s#%s\n",
		message.AGENCIA,
		message.NOMBRE,
		message.APELLIDO,
		message.DOCUMENTO,
		message.NACIMIENTO,
		message.NUMERO,
	)	
}