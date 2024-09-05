
type message struct {
	AGENCIA string
	NOMBRE string
	APELLIDO string
	DOCUMENTO string
	NACIMIENTO string
	NUMERO string
}

func NewMessage(AGENCIA, NOMBRE, APELLIDO, DOCUMENTO, NACIMIENTO, NUMERO string) *message {
	return &message{
		AGENCIA: AGENCIA,
		NOMBRE: NOMBRE,
		APELLIDO: APELLIDO,
		DOCUMENTO: DOCUMENTO,
		NACIMIENTO: NACIMIENTO,
		NUMERO: NUMERO,
	}
}

func (message *message) Serialize() string{
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s-%s\n",
		message.AGENCIA,
		message.NOMBRE,
		message.APELLIDO,
		message.DOCUMENTO,
		message.NACIMIENTO,
		message.NUMERO,
	)	
}
