# Trabajo practico 0 Sistemas Distribuidos - Fernando Bursztyn



## Parte 2 repaso de comunicaciones 

Para esta seccion del trabajo  y la comunicacion entre cliente y servidor, se implemento un protocolo de texto dividos por un caracter especial entre los distintos datos y el uso de `\n` o `\n\n` para indicar el fin del de la apuesta o el fin del mensaje.

Para realizar este trabajo escogi un protocolo de texto, el cual envia mensajes compuesto por campos de texto, los cuales se separan por un
caracter especial, un salto de linea. Los mensajes culminan con un doble salto de linea, y comienzan con un primer parametro que indica el tipo de mensaje.


### Detalles de implemetacion

El cliente envia mensajes del tipo MESSAGE en donde contiene todo los campos de la apuesta. Estos son: nombre, apellido, documento, fecha de nacimiento, y numero de la apuesta. Ademas de a que agencia correspode. El envio del mensaje se hace con una concatenacion de strings de la forma: PARAM1#PARAM2#PARAM3#PARAM4#PARAM5#PARAM6\n, siendo # el caracter especial encargador de separar cada parametro de la apuesta, finalizando con salto de linea que da inicio a la siguiete apuesta teniendo un limite de N apuestas por segmeto (batch). Ademas, al final del la tira de string se agrega un `\n` adicional para indicar el fin del mensaje. 


El server a medida que va recibiendo la informacion, con su propio protocolo, arma la tira de apuestas provenietes de la tira de string del batch. Y por cada apuesta la agrega a su base de datos. Cada cliente procesa su archivo de apuestas por partes, despues de cada mensaje que envia se desconecta dejandole lugar al otro cliente.


Una vez que el cliente envia todos los datos del archivo envia un mensaje que ya termino, seguido despues de un mensaje donde pide los resultados. El server una vez que registra que todos los clietes terminaro con sus respectivos archivos empieza a enviar los resultados.

Siedo del siguiente formato: 

%v#NOTIFY_DONE\n\n

%v#REQUEST_WINNERS\n\n

En caso que el server no tenga la confirmacion de todos los clientes corta su comunicacion con dicho cliente, haciendo que espere (dormido) y vuelva a solicitarlo despues de un lapso de tiempo.


## Mecanismo de sincronizacion - parte 3.

Se hizo uso de dos herramientas, barries y locks. Ademas de la libreria de multiproccessing de python.


### Locks

Hay un recurso compartido que es la base de datos de apuestas en donde se guardan las nuevas y se consulta para obtener los ganadores.

### Barrier

Cuando una agencia termina de enviar sus datos, le avisa al server que ya termino y este queda esperando la barrera ubicada para sincronizar con el resto de procesos. Una vez que se encuetran todos se dan los datos de los ganadores.