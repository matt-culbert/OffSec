A switchblade style C2 WIP

This is based on the CIA design below:
![image](https://user-images.githubusercontent.com/18468466/142744338-70845a6f-733b-4847-9432-a68a5d5e8426.png)

The idea is that beacons that reach out will provide a series of certs needed for mutual TLS authentication to port 443. If these are valid, they are forwarded to the python server. If they aren't valid, they are delivered to a fake webpage. This keeps your C2 open to the public but only authorized clients can actually interact with it.
![with-without-cert](https://user-images.githubusercontent.com/18468466/142713549-979c1b07-0e3f-480b-98a4-c7c6d816f513.png)

To compile client.go, run gccgo client.go -o client
