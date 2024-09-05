import socket
import logging
import signal
import multiprocessing
from .utils import store_bets, load_bets, has_won
from .protocol import Protocol
from .message import Message

import json

class Server:
    def __init__(self, port, listen_backlog, total_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._agency_notifications = 0
        self.winners = {}
        self.clientsReady = False
        self.connected_clients = []
        self.total_clients = int(total_clients)

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

        lock_bets = multiprocessing.Lock()
        barrier = multiprocessing.Barrier(self.total_clients)

        while self._running:
            try:
                client_sock = self.__accept_new_connection()

                client_process = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock, lock_bets, barrier))
                self.connected_clients.append(client_process)

                client_process.start()
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
    

    
    def __handle_client_connection(self, client_sock, betLock, barrier):

        """
        Handle communication with a client
        """
        process_client_running = True

        while process_client_running and self._running:

            protocol = Protocol(client_sock)
            msg = protocol.recv_all()
            newMessage = Message(msg)


            #while process_client_running and self._running:
            try:
                if msg.__contains__("NOTIFY_DONE"):
                    barrier.wait()
                    logging.info(f"msg: {msg}")
                    #self._agency_notifications += 1
                    #logging.info(f"action: notify_received | result: success | notified_agencies: {self._agency_notifications}/2")

                    # if self._agency_notifications == 5:
                    #     logging.info("action: sorteo | result: success")
                    #     self.clientsReady = True

                elif msg.__contains__("REQUEST_WINNERS"):
                    print("server entra a request winners")
                
                    ID = newMessage.deserializeRequestWinners()
                    #print(f"ID: {ID}")
                    # if not self.clientsReady:
                    #     #print("No se puede enviar los ganadores")
                    #     client_sock.close()
                    # else:
                    #     #print("Se puede enviar los ganadores")
                    winners = []
                    with betLock:
                        all_bets = load_bets()
                    for bet in all_bets:
                        if bet.agency == int(ID) and has_won(bet):
                            winners.append(bet.document)
                    #print(f"agency_bets_count: {agency_bets_count}")
                    protocol.winnerToAgency(winners)
                    process_client_running = False

                else:
                    #newMessage = Message(msg)
                    bets = newMessage.deserialize()
                    #print(f"bets: {bets}")

                    with betLock:
                        store_bets(bets)

                    total_bets = 0 
                    for bet in load_bets():
                        total_bets += 1
                    
                    logging.info(f"action: apuestas_almacenadas | result: success | cantidad de apuestas: {len(bets)} | cant totales: {total_bets}")
            except OSError as e:
                        logging.error("action: receive_message | result: fail | error: {e}")
            #finally:
        client_sock.close()
        

