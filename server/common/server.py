import socket
import logging
import signal
from .utils import store_bets, load_bets
from .protocol import Protocol
from .message import Message

import json

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True

        signal.signal(signal.SIGINT, self.graceful_shutdown)
        signal.signal(signal.SIGTERM, self.graceful_shutdown)

    def graceful_shutdown(self, signum, frame):
        """
        Graceful shutdown of the server
        """
        logging.info("action: graceful_shutdown | result: in_progress")
        self._running = False
        self._server_socket.close()
        logging.info("action: graceful_shutdown | result: success")

    def run(self):
        """
        Server loop to accept new connections and handle communication
        """
        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError:
                logging.error("action: server_run | result: stopped")

    def __handle_client_connection(self, client_sock):

        """
        Handle communication with a client
        """
        protocol = Protocol(client_sock)
        msg = protocol.recv_all()
        newMessage = Message(msg)
        bets = newMessage.deserialize()

        store_bets(bets)

        total_bets = 0 
        for bet in load_bets():
            total_bets += 1
        
  
        logging.info(f"action: apuestas_almacenadas | result: success | cantidad de apuestas: {len(bets)} | cant totales: {total_bets}")

        

        

    def __accept_new_connection(self):
        """
        Accept a new connection
        """
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c