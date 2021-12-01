
class Clients{
        constructor(){
                this.clientList = {};
                this.saveClient = this.saveClient.bind(this);
        }
        saveClient(username, client){
                this.clientList[username] = client;
        }
}

const readline = require('readline').createInterface({
  input: process.stdin,
  output: process.stdout
});

var WebSocketServer = require('ws').Server, wss = new WebSocketServer({port: 8010});

const clients = new Clients();

const run = async() =>{
        wss.on('connection', (client) => {
                client.on('message', (msg) => {
                        const parsedMsg = JSON.parse(msg);
                        clients.saveClient(parsedMsg.username, client);
                        //clients.clientList[parsedMsg.username].send("'Check in'");
                        readline.question('>',query =>{clients.clientList[parsedMsg.username].send(query);});
                });
        });
}

run();
