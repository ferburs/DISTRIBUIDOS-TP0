#!/bin/bash

# Verifica si se proporcionaron los dos parámetros necesarios
if [ $# -ne 2 ]; then
    echo "Uso: $0 <nombre_archivo_salida> <cantidad_clientes>"
    exit 1
fi

# Asigna los parámetros a variables
nombre_archivo_salida=$1
cantidad_clientes=$2

echo "Nombre del archivo de salida: $nombre_archivo_salida"
echo "Cantidad de clientes: $cantidad_clientes"

# Llama al subscript de Python para generar el archivo Docker Compose
python3 clientes.py $nombre_archivo_salida $cantidad_clientes

echo "Archivo Docker Compose generado: $nombre_archivo_salida"
