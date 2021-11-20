A switchblade style C2 WIP

The idea is that beacons that reach out will provide a series of certs needed for mutual TLS authentication to port 443. If these are valid, they are forwarded to the python server. If they aren't valid, they are delivered to a fake webpage. This keeps your C2 open to the public but only authorized clients can actually interact with it.
