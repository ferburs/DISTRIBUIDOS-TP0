

class Protocol:
    def __init__(self, sock):
        self.sock = sock

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

    def winnerToAgency(self,listWinners):
        """
        Send the winner count to the client
        """
        msg = "\n".join(listWinners) + "\n\n"
        self.send_all(msg.encode('utf-8'))