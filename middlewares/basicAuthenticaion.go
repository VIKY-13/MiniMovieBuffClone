package middlewares

import (
	"net/http"
	"golangmovietask/config"
)

func HostAuthentication(next http.HandlerFunc)http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		requsername,reqpassword,ok := r.BasicAuth()
		if ok{
			username := config.AUTH_USERNAME
			password := config.AUTH_PASSWORD
			if username==requsername && password==reqpassword{
				next.ServeHTTP(w,r)
				return
			}
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}