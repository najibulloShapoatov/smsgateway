package routes

import (
	"net/http"
	"smsc/mini/handler"
	"smsc/mini/middleware"
	"smsc/pkg/log"

	"github.com/gorilla/mux"
)

func Init(path string) *mux.Router {
	var router = mux.NewRouter()
	router.Use(CORS)
	router.Use(LoggingMiddleware)
	router.Use(Recovery)
	router.Use(middleware.AuthAPI)

	router.HandleFunc("/send", handler.Sent)

	return router
}

//Recovery ....
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				log.Error("Recover", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

		}()

		next.ServeHTTP(w, r)

	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here

		log.Println(r.RequestURI)
		log.Println(r.Body)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*r).Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)

	})
}
