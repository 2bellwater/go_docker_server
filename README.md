#docker file for golang server

##prepare redis docker

1. docker install
2. redis image pull
3. run redis docker
$docker run --name jongsoo-redis -d -p 6379:6379 redis

##run golang server

1.git clone https://github.com/mekingme/go_docker_server.git
2.cd go_docker_server
3.run 
$./go_docker_server


##example
$ curl http://localhost:8080/hello
hello GET

$ curl http://localhost:8080/redis
GET ID:default idGET PW:default pw
