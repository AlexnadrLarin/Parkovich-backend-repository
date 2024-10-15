package internal

import (
    "net/http"
	
	"github.com/swaggo/http-swagger"
)

func SetupRouter()  {
    http.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
        httpSwagger.WrapHandler.ServeHTTP(w, r)
    })
}
