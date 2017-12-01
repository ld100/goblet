//package main
//
//import (
//	"fmt"
//	"os"
//
//	redis "gopkg.in/redis.v4"
//)
//
//func main() {
//	client := redis.NewClient(&redis.Options{
//		Addr:     os.Getenv("REDIS_URL"),
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//
//	pong, err := client.Ping().Result()
//	fmt.Println(pong, err)
//}

package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
)

 func main() {
 	fmt.Println("Hello World")
 	r := chi.NewRouter()
 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
 		w.Write([]byte("welcome"))
 	})
 	http.ListenAndServe(":8080", r)
 }
