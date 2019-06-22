package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type routerWrapper struct {
	router *mux.Router
}

func newRouterWrapper(router *mux.Router) *routerWrapper {
	return &routerWrapper{router}
}

func (r *routerWrapper) RegisterJson(path string, handler func(HttpContext) (*JsonResponse, error)) {
	r.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		context := HttpContext{
			Args: mux.Vars(r),
			Body: r.Body,
		}

		response, err := handler(context)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		b, err := json.Marshal(response.Model)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		w.Write(b)
	})
}

func (r *routerWrapper) RegisterData(path string, handler func(HttpContext) (*DataResponse, error)) {
	r.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		context := HttpContext{
			Args: mux.Vars(r),
			Body: r.Body,
		}

		response, err := handler(context)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Println(response.Error)
			w.WriteHeader(response.StatusCode)
			return
		}
		defer response.Data.Close()

		w.Header().Add("Content-Type", "application/octet-stream")
		w.WriteHeader(response.StatusCode)

		wb := bufio.NewWriter(w)
		rb := bufio.NewReader(response.Data)
		defer wb.Flush()

		buf := make([]byte, 1024)
		for {
			// read a chunk
			n, err := rb.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			if n == 0 {
				break
			}

			// write a chunk
			if _, err := w.Write(buf[:n]); err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
		}
	})
}
