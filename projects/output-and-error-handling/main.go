package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	url := "http://localhost:8080/"
	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Errorf("Oops, something went wrong. Please try again later\n Details: %w\n", err).Error())
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			htmlData, err := io.ReadAll(resp.Body) //<--- here!
			if err != nil {
				fmt.Fprintf(os.Stderr, fmt.Errorf("Body oops: %w\n", err).Error())
				os.Exit(2)
			}
			fmt.Println(string(htmlData))
			break
		} else if resp.StatusCode == 429 {
			delayHeader := resp.Header["Retry-After"][0]
			i, err := strconv.Atoi(delayHeader)
			if err == nil {
				if i < 6 {
					if i > 1 {
						fmt.Fprintln(os.Stderr, "Server busy ... retrying")
					}
					time.Sleep(time.Duration(i) * time.Second)
				} else {
					fmt.Fprintln(os.Stderr, "Server too busy")
					break
				}
			} else {
				fmt.Fprintln(os.Stderr, "Server busy ... retrying in 2 seconds")
				time.Sleep(2 * time.Second)
			}
		} else {
			break
		}
	}
}
