class Message:
    def init (self, message):
        self.bets = message 


    def deserialize(self):
        """
        Deserialize the message
        """
        betstring = self.bets.split('\n')
        bets = []

        for bet in betstring:
            betArg = bet.split('-')
            bets.append(Bet(betArg[0], betArg[1], betArg[2], betArg[3],betArg[4],betArg[5]))
        return bets
