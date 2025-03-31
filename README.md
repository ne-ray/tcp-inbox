# tcp-inbox

# Requirements
Design and implement “Word of Wisdom” tcp server.
 • TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
 • The choice of the POW algorithm should be explained.
 • After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
 • Docker file should be provided both for the server and for the client that solves the POW challenge


# Build & Run server

## Build Docker container
```bash
docker build -f server.Dockerfile ./ -t nti_server:latest
```

## Run server over docker
```bash
docker run -d -p 19777:19777/tcp --name nti_server nti_server:latest
```
