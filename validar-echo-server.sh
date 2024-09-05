#!/bin/bash

# Definir el nombre del contenedor y el puerto
SERVER_CONTAINER="server"
PORT="12345"
MESSAGE="Hello World!"

# Obtener la dirección IP del contenedor
RESPONSE=$(docker run --rm --network tp0_testing_net busybox sh -c "echo \"$MESSAGE\" | nc \"$SERVER_CONTAINER\" \"$PORT\"")

# Verificar si la respuesta coincide con el mensaje enviado
if [ "$RESPONSE" = "$MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
