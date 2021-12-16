A switchblade style C2 WIP

This is based on the CIA design below:
![image](https://user-images.githubusercontent.com/18468466/142744338-70845a6f-733b-4847-9432-a68a5d5e8426.png)

The idea is that beacons that reach out will provide a series of certs needed for mutual TLS authentication to port 443. If these are valid, they are forwarded to the back-end server. If they aren't valid, they are delivered to a fake webpage. This keeps your C2 open to the public but only authorized clients can actually interact with it.
![with-without-cert](https://user-images.githubusercontent.com/18468466/142713549-979c1b07-0e3f-480b-98a4-c7c6d816f513.png)

Run single_client_server.js and then run the test-websocket.py to get going.

To generate certs:

openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 365 -key ca.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out ca.crt

openssl req -newkey rsa:2048 -nodes -keyout server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=_*.example.com" -out server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:example.com,DNS:www.example.com") -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

Taken from here https://security.stackexchange.com/questions/74345/provide-subjectaltname-to-openssl-directly-on-the-command-line

Beacon.go now works with the above certs supplied to an mTLS enabled C2 server taken from BlackhatGo
