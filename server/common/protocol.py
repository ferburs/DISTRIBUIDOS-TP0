

class Protocol:
    def __init__(self, sock):
        self.sock = sock

    # def recv_message(self):
    #     """
    #     Receive and deserialize a JSON message from the socket
    #     """
    #     data = self.__recv_all().decode('utf-8')
    #     return json.loads(data)

    # def send_message(self, message):
    #     """
    #     Serialize and send a JSON message through the socket
    #     """
    #     data = json.dumps(message).encode('utf-8')
    #     self.__send_all(data)

    def recv_all(self):
        """
        Receive all data from the socket
        """
        data = b''
        while True:
            part = self.sock.recv(1024)
            if not part:
                break
            data += part
            if b'\n\n' in data:
                break
                
        return data[:data.find(b'\n\n')].strip().decode('utf-8')

    def send_all(self, data):
        """
        Send all data to the socket
        """
        total_sent = 0
        while total_sent < len(data):
            sent = self.sock.send(data[total_sent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            total_sent += sent

    def winnerToAgency(self,agency_bets_count):
        """
        Send the winner count to the client
        """
        msg = str(agency_bets_count) + "\n"
        #print (f"msg: {msg}")
        self.send_all(msg.encode('utf-8'))
        #self.sock.send(str(agency_bets_count).encode('utf-8') + "\n".encode('utf-8'))