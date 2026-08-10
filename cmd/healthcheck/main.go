package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	url := os.Getenv("HEALTHCHECK_URL")
	if url == "" {
		url = "http://127.0.0.1:8000/health"
	}

	client := http.Client{Timeout: 2 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		fmt.Fprintf(os.Stderr, "unexpected health status: %d\n", res.StatusCode)
		os.Exit(1)
	}

	fmt.Println("ok")
}
