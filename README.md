# tcp-inbox

# Requirements
Design and implement “Word of Wisdom” tcp server.
 • TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
 • The choice of the POW algorithm should be explained.
 • After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
 • Docker file should be provided both for the server and for the client that solves the POW challenge

# Choosed PoW
## Fiat Shamir
In first time I choose this algorithm, why? 
It have two way communication with multiple round request/response, but this algorithm have same difficile for server and client

## Hashcash
This algorithm have all options for Proof of Work protection will work
1. One request - one work (for work server generate randow string and wait client calculate PoW)
2. Algorithm have simpe work for server and difficile for client
And other argumend for choose this
1. Algorithm have implementation for Golang

# Build & Run server

## Build Docker container
```bash
docker build -f server.Dockerfile ./ -t nti_server:latest
```

## Run server over docker
```bash
docker run -d -p 19777:19777/tcp --name nti_server nti_server:latest
```

# Build & Run client

## Build Docker container
```bash
docker build -f client.Dockerfile ./ -t nti_client:latest
```

## Run client over docker
```bash
docker run --rm -i nti_client:latest ./client post -s host.docker.internal -l %line book of "Word of Wisdom"%  -c %chapter book of "Word of Wisdom"%
```

for more information send data and flag use command:
```bash
docker run --rm nti_client:latest ./client post --help
```