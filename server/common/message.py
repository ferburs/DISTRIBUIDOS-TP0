from .utils import Bet

class Message:
    def __init__ (self, message):
        self.msg = message 


    def deserialize(self):
        """
        Deserialize the message
        """
        betstring = self.msg.split('\n')
        bets = []

        for msg in betstring:
            betArg = msg.split('#')
            bets.append(Bet(betArg[0], betArg[1], betArg[2], betArg[3],betArg[4], betArg[5]))
        return bets
    
    def deserializeRequestWinners(self):

        id = self.msg.split('#')[0]
        return id
