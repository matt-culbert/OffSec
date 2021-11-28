import asyncio
import pathlib
import ssl
import websockets

ssl_context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
#localhost_pem = pathlib.Path(__file__).with_name("/etc/nginx/certs/key.crt")
ssl_context.load_verify_locations("/etc/nginx/certs/server.crt")
SSLContext.load_cert_chain("/etc/nginx/certs/ca.crt", keyfile=None, password=None)

async def hello():
    uri = "wss://test.temp:443"
    async with websockets.connect(
        uri, ssl=ssl_context
    ) as websocket:
        name = input("What's your name? ")

        await websocket.send(name)
        print(f"> {name}")

        greeting = await websocket.recv()
        print(f"< {greeting}")

asyncio.get_event_loop().run_until_complete(hello())
