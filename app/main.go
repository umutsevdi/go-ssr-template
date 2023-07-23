package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/umutsevdi/webwatch/app/api"
	"github.com/umutsevdi/webwatch/app/pages"
	"github.com/umutsevdi/webwatch/app/util"
)

func main() {
	go pages.Serve()
	go api.Serve()
	str := ":" + strconv.Itoa(int(*util.C.Port))
	log.Fatal(http.ListenAndServe(str, nil))
}
