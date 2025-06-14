package middleware

import(
	"net/http"
	"log"
)


//for log every request | maybe better to store in file !!!!!!!!!!!!! 
func LoggerMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("---> Paths[ %s ] | MethodType[ %s ]", r.URL.Path, r.Method)
        next.ServeHTTP(w, r) //Go to the next handler
    })
}