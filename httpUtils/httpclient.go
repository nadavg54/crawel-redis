package httpclient

import (
	"io"
	"net/http"
)

//rertrive html content
func getContent(url string) (*io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil,err
	}
	return &(resp.Body),nil
}
