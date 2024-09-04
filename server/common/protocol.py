import json

class Protocol:
    def __init__(self, sock):
        self.sock = sock

    def recv_message(self):
        """
        Receive and deserialize a JSON message from the socket
        """
        data = self.__recv_all().decode('utf-8')
        return json.loads(data)

    def send_message(self, message):
        """
        Serialize and send a JSON message through the socket
        """
        data = json.dumps(message).encode('utf-8')
        self.__send_all(data)

    def __recv_all(self):
        """
        Receive all data from the socket
        """
        data = b''
        while True:
            part = self.sock.recv(1024)
            data += part
            if len(part) < 1024:
                break
        return data

    def __send_all(self, data):
        """
        Send all data to the socket
        """
        total_sent = 0
        while total_sent < len(data):
            sent = self.sock.send(data[total_sent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            total_sent += sent