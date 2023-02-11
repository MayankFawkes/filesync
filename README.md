# filesync ⚠️ Not Production ready
Sync all your SSL and other small files around the servers with docker-compose


## Features

- Light Weight & Easy to setup
- Docker image avaiable 
- One command installation


## Links
- [https://hub.docker.com/r/mayankfawkes/httptoy](https://hub.docker.com/r/mayankfawkes/httptoy "https://hub.docker.com/r/mayankfawkes/httptoy")

## Setup 

I have 2 servers Master `172.31.37.244` and Worker `172.31.37.245`, I am connecting them with local ip since i have both servers at same datacenter, you can use public ip if you want make sure the port is accessible. 

### Master

```
version: "3.7"

services:
  web:
    image: mayankfawkes/filesync:latest
    ports:
      - 8000:8000
    environment:
      NODE: MASTER
      WATCH_PATH: './data'
      PORT: 8000
    volumes:
      - ./data:/data:ro
```

### Worker

```
version: "3.7"

services:
  web:
    image: mayankfawkes/filesync:latest
    ports:
      - 8000:8000
    environment:
      NODE: WORKER
      WATCH_PATH: './data'
      PORT: 8000
      MASTER_IP: '172.31.37.244'
      MASTER_PORT: 8000
    volumes:
      - ./data:/data

```


All the file of Master node at location `./data` will be live sync with worker at `./data`



## Conclusion

I had alot of problem related to letsencrypt ssl and other local environment files to be updated in all of my server so i made this lite solution.


## Use Cases

I have everything manual i use multiple servers loadbalancer and i dont wanna update ssl on all servers after every 3 months. i use docker-compose with my web services and nginx i wanted a solution from which i can edit server config files in one server and it will get modefied in all servers. you can also check [Rsync](https://en.wikipedia.org/wiki/Rsync "https://en.wikipedia.org/wiki/Rsync") 
