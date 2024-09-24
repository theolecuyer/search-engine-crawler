package main

import (
	"io"
	"net/http"
)

func download(url string) (string, error) {
	if rsp, err := http.Get(url); err != nil {
		return "", err
	} else {
		defer rsp.Body.Close()
		if bts, err := io.ReadAll(rsp.Body); err != nil {
			return "", err
		} else {
			return string(bts), err
		}
	}
}
