package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	resp, err := http.Post("http://127.0.0.1:9999/helloworld", "text/html; charset=utf-8", strings.NewReader(`{"name":"zhou"}`))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}
