package frameworks

import (
	"fmt"
	"net/http"
)

func StartServer(router *Router) {
	address := ":8080"
	fmt.Println("Server started on", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
