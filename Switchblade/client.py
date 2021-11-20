#!/usr/bin/python3

import socket
import ssl
import time

host_addr = '192.168.202.135'
host_port = 443
server_sni_hostname = 'test.temp'
server_cert = '/etc/nginx/certs/ca.crt'
client_cert = '/etc/nginx/certs/client.crt'
client_key = '/etc/nginx/certs/client.key'

context = ssl.create_default_context(ssl.Purpose.SERVER_AUTH, cafile=server_cert)
context.load_cert_chain(certfile=client_cert, keyfile=client_key)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
conn = context.wrap_socket(s, server_side=False, server_hostname=server_sni_hostname)
conn.connect((host_addr, host_port))
print("SSL established. Peer: {}".format(conn.getpeercert()))

request = b"GET / HTTP/1.1\nHost: test.temp\n\n"

conn.send(request)
result = conn.recv(10000)
while (len(result) > 0):
    print(result)
    result = conn.recv(10000)
