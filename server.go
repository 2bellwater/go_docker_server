package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Resource interface {
	Uri() string
	Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Put(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Delete(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
}

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (PostNotSupported) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (PutNotSupported) Put(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (DeleteNotSupported) Delete(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func HttpResponse(rw http.ResponseWriter, req *http.Request, res Response) {
	content, err := json.Marshal(res)

	if err != nil {
		abort(rw, 500)
	}

	rw.WriteHeader(res.Code)
	rw.Write(content)
}

func AddResource(router *httprouter.Router, resource Resource) {
	fmt.Println("\"" + resource.Uri() + "\" api is registerd")

	router.GET(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Get(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.POST(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Post(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.PUT(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Put(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.DELETE(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Delete(rw, r, ps)
		HttpResponse(rw, r, res)
	})
}

type HelloResource struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

func (HelloResource) Uri() string {
	return "/hello"
}

func (HelloResource) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	fmt.Println("hello GET")
	return Response{200, "message", map[string]interface{}{
		"key1": "value1",
	}}
}

func RadisGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	getID, err := client.Get("id").Result()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "GET ID:%s", getID)

	getPW, err := client.Get("pw").Result()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "GET PW:%s", getPW)

}

func RadisPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	pw := ps.ByName("pw")

	fmt.Fprintf(w, "RadisPost! ID:%s PW:%s\n", id, pw)

	err := client.Set("ID", id, 0).Err()
	if err != nil {
		panic(err)
	}

	err = client.Set("PW", pw, 0).Err()
	if err != nil {
		panic(err)
	}

}

var client *redis.Client

func main() {

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	e := client.Set("id", "default id", 0).Err()
	if e != nil {
		panic(e)
	}
	e = client.Set("pw", "default pw", 0).Err()
	if e != nil {
		panic(e)
	}

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	router := httprouter.New()
	fmt.Print("httpServerStarted 8080 port")

	router.POST("/redis/input:id,pw", RadisPost)
	router.GET("/redis", RadisGet)
	AddResource(router, new(HelloResource))
	log.Fatal(http.ListenAndServe(":8080", router))

}
