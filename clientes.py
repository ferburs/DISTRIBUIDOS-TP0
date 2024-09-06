import sys

def generar_compose(nombre_archivo, cantidad_clientes):
    with open(nombre_archivo, 'w') as archivo:
        archivo.write("name: tp0\n")
        archivo.write("services:\n")
        
        # Definir el servicio del servidor
        archivo.write("  server:\n")
        archivo.write("    container_name: server\n")
        archivo.write("    image: server:latest\n")
        archivo.write("    entrypoint: python3 /main.py\n")
        archivo.write("    environment:\n")
        archivo.write("      - PYTHONUNBUFFERED=1\n")
        archivo.write("      - LOGGING_LEVEL=DEBUG\n")
        archivo.write("    networks:\n")
        archivo.write("      - testing_net\n")
        archivo.write("    volumes:\n")
        archivo.write("      - ./server/config.ini:/config.ini\n\n")

        
        # Definir los servicios de los clientes
        for i in range(1, int(cantidad_clientes) + 1):
            archivo.write(f"  client{i}:\n")
            archivo.write(f"    container_name: client{i}\n")
            archivo.write("    image: client:latest\n")
            archivo.write("    entrypoint: /client\n")
            archivo.write("    environment:\n")
            archivo.write(f"      - CLI_ID={i}\n")
            archivo.write("      - CLI_LOG_LEVEL=DEBUG\n")
            archivo.write("    networks:\n")
            archivo.write("      - testing_net\n")
            archivo.write("    depends_on:\n")
            archivo.write("      - server\n")
            archivo.write("    volumes:\n")
            archivo.write("      - ./client/config.yaml:/config.yaml\n\n")
        
        # Definir la red
        archivo.write("networks:\n")
        archivo.write("  testing_net:\n")
        archivo.write("    ipam:\n")
        archivo.write("      driver: default\n")
        archivo.write("      config:\n")
        archivo.write("        - subnet: 172.25.125.0/24\n")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Uso: python3 mi-generador.py <nombre del archivo de salida> <cantidad de clientes>")
        sys.exit(1)

    nombre_archivo = sys.argv[1]
    cantidad_clientes = sys.argv[2]
    
    generar_compose(nombre_archivo, cantidad_clientes)
