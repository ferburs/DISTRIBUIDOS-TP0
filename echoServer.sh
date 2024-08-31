#!/bin/bash

# Definir el nombre del contenedor y el puerto
SERVER_CONTAINER="server"
PORT="12345"
MESSAGE="Hello World!"

# Obtener la dirección IP del contenedor
SERVER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$SERVER_CONTAINER")

# Enviar el mensaje al servidor usando netcat y guardar la respuesta
RESPONSE=$(echo "$MESSAGE" | nc "$SERVER_IP" "$PORT")

# Imprimir la respuesta recibida
echo "$RESPONSE"

# Verificar si la respuesta coincide con el mensaje enviado
if [ "$RESPONSE" = "$MESSAGE" ]; then
  echo "OK"
else
  echo "ERROR"
  echo "Expected: $MESSAGE"
  echo "Received: $RESPONSE"
fi
