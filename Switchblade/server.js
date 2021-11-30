class Clients{
        constructor(){
                this.clientList = {};
                this.saveClient = this.saveClient.bind(this);
        }
        saveClient(username, client){
                this.clientList[username] = client;
        }
}

const WebSocket = require('ws')

const clients = new Clients();
const socket = new WebSocket.Server({ port: 8010});
socket.on('connection', (client) => {
        client.on('message', (msg) => {
                const parsedMsg = JSON.parse(msg);
                clients.saveClient(parsedMsg.username, client);
        });
});
