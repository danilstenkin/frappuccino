package router

import (
	"net/http"
)

func SetupRouter() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
