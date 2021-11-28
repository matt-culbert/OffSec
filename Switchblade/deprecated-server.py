#!/usr/bin/python3

import socket
from socket import AF_INET, SOCK_STREAM, SO_REUSEADDR, SOL_SOCKET, SHUT_RDWR
import ssl

listen_addr = '192.168.202.135'
listen_port = 8082
#server_cert = '/etc/nginx/certs/server.crt'
#server_key = '/etc/nginx/certs/server.key'
#client_certs = '/etc/nginx/certs/ca.crt'

#context = ssl.create_default_context(ssl.Purpose.CLIENT_AUTH)
#context.verify_mode = ssl.CERT_REQUIRED
#context.load_cert_chain(certfile=server_cert, keyfile=server_key)
#context.load_verify_locations(cafile=client_certs)

bindsocket = socket.socket()
bindsocket.bind((listen_addr, listen_port))
bindsocket.listen(5)

while True:
    print("Waiting for client")
    newsocket, fromaddr = bindsocket.accept()
    print("Client connected: {}:{}".format(fromaddr[0], fromaddr[1]))
