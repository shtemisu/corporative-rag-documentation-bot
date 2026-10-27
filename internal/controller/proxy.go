package controller

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func NewProxy(proxyURL string) *http.Client {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		log.Fatalf("failed to parse url: %v", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedURL),
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}
	return httpClient
}
