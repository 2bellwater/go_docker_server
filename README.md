# docker file for golang server

## prepare redis docker

### 1. docker install
### 2. run redis docker
$ docker run --name jongsoo-redis -d -p 6379:6379 redis
### 3. run mysql docker 
$ docker run --name js-mysql -e MYSQL_ROOT_PASSWORD=jspassword -d -p 3306:3306  mysql
### 4. create default database
$ docker exec -it js-mysql mysql -u root -p  
Enter password: jspassword     
mysql> create database jsdatabase;  
mysql> use jsdatabase;  
mysql> create table jstable( uuid int not null auto_increment, userid varchar(30), userpw varchar(30), studytime int, primary key(uuid) );  

## run golang server

1.git clone https://github.com/mekingme/go_docker_server.git
2.cd go_docker_server
3.run 
$./go_docker_server

## example
$ curl http://localhost:8080/hello
hello GET

$ curl http://localhost:8080/redis
GET ID:default idGET PW:default pw
