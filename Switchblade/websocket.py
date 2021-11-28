import socket
import ssl
import websockets
import asyncio

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

sslSettings.load_cert_chain(certfile="/etc/nginx/certs/server.crt", keyfile="/etc/nginx/certs/server.key")

# Create a stream based client socket

async def hello():
    uri = "wss://192.168.202.135:443"
    async with websockets.connect(
        uri, ssl=sslSettings
    ) as websocket:
        name = input("What's your name? ")

        await websocket.send(name)
        print(f"> {name}")

        greeting = await websocket.recv()
        print(f"< {greeting}")

asyncio.get_event_loop().run_until_complete(hello())
