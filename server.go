package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"net/http"
)

func main() {
	var (
		listen = "localhost:8484"
		size   int
	)
	flag.IntVar(&size, "size", 256, "QRcode size in pixels")
	flag.Parse()
	if len(flag.Args()) > 0 {
		listen = flag.Arg(0)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			spayd = Spayd{}
		)
		fmt.Print(r.Body)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&spayd)
		if err != nil {
			http.Error(w, "JSON decoding failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		query, err := spayd.Encode()
		if err != nil {
			http.Error(w, "SPAYD encoding failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		png, err := qrcode.Encode(string(query), qrcode.Medium, size)
		if err != nil {
			http.Error(w, "QR generation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "image/png")
		w.Write(png)
	})
	fmt.Println("Listening at " + listen)
	http.ListenAndServe(listen, nil)
}
