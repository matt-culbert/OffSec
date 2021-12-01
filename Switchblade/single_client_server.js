console.log("Server started");
var Msg = '';

const readline = require('readline').createInterface({
  input: process.stdin,
  output: process.stdout
});

var WebSocketServer = require('ws').Server
    , wss = new WebSocketServer({port: 8010});
    wss.on('connection', function(ws) {
        ws.on('message', function(message) {
        console.log('Received from client: %s', message);
        readline.question('>',query =>{ws.send(query);});;
    });
 });
