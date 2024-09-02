import socket
import logging
import signal
from common import utils

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

        Function that closes the server socket and exits the program
        """
        logging.info("action: graceful_shutdown | result: in_progress")
        self._running = False
        self._server_socket.close()
        logging.info("action: graceful_shutdown | result: success")
        #sys.exit(0)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError: 
                logging.error("action: server_run | result: stopped")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # # TODO: Modify the receive to avoid short-reads
            # msg = client_sock.recv(1024).rstrip().decode('utf-8')
            # addr = client_sock.getpeername()
            # logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # # TODO: Modify the send to avoid short-writes
            # client_sock.send("{}\n".format(msg).encode('utf-8'))

            msg = self.__recv_all(client_sock).decode('utf-8')
            bet = json.loads(msg)

            newBet = utils.Bet(bet['agency'], bet['NOMBRE'], bet['APELLIDO'], bet['DOCUMENTO'], bet['NACIMIENTO'], bet['NUMERO'])
            utils.store_bets([newBet])

            logging.info(f'action: bet accepted | result: success | bet: {bet}')


            response = json.dumps({"status": "success"})
            self.__send_all(client_sock, response.encode('utf-8'))

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __recv_all(self, sock):
        """
        Receive all data from a specific socket

        Function receives all data from a specific socket and returns it
        """
        data = b''
        while True:
            part = sock.recv(1024)
            data += part
            if len(part) < 1024:
                break
        return data
    

    def __send_all(self, sock, data):
        """
        Send all data to a specific socket

        Function sends all data to a specific socket
        """
        total_sent = 0
        while total_sent < len(data):
            sent = sock.send(data[total_sent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            total_sent += sent