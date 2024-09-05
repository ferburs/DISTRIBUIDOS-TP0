import socket
import logging
import signal
from .utils import store_bets, load_bets, has_won
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
        self._agency_notifications = 0
        self.winners = {}
        self.clientsReady = False

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

    
    
    def __accept_new_connection(self):
        """
        Accept a new connection
        """
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    

    
    def __handle_client_connection(self, client_sock):

        """
        Handle communication with a client
        """
        protocol = Protocol(client_sock)
        msg = protocol.recv_all()
        newMessage = Message(msg)

        #logging.info(f"\n\n mensaje: {msg} \n\n ")
        try:
            if msg.__contains__("NOTIFY_DONE"):
                logging.info(f"msg: {msg}")
                self._agency_notifications += 1
                logging.info(f"action: notify_received | result: success | notified_agencies: {self._agency_notifications}/2")

                if self._agency_notifications == 5:
                    #logging.info("ARRANCA O NO ARRANCA EL SORTEO. SIEMPRE ARRACA CON BUJIAS JECHER")
                    logging.info("action: sorteo | result: success")
                    self.clientsReady = True

            elif msg.__contains__("REQUEST_WINNERS"):
                #logging.info(f"msg: {msg}")
                ID = newMessage.deserializeRequestWinners()
                #print(f"ID: {ID}")
                if not self.clientsReady:
                    #print("No se puede enviar los ganadores")
                    client_sock.close()
                else:
                    #print("Se puede enviar los ganadores")
                    all_bets = load_bets()
                    agency_bets_count = sum(1 for bet in all_bets if bet.agency == int(ID) and has_won(bet))
                    #print(f"agency_bets_count: {agency_bets_count}")
                    protocol.winnerToAgency(agency_bets_count)

            else:
                #newMessage = Message(msg)
                bets = newMessage.deserialize()

                store_bets(bets)

                total_bets = 0 
                for bet in load_bets():
                    total_bets += 1
                
                logging.info(f"action: apuestas_almacenadas | result: success | cantidad de apuestas: {len(bets)} | cant totales: {total_bets}")
        except OSError as e:
                    logging.error("action: receive_message | result: fail | error: {e}")
        finally:
                client_sock.close()

