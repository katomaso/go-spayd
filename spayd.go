package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"net/http"
	"reflect"
	"strconv"
)

type Spayd struct {
	Account string `spayd:"ACC" max_len:"46" format:"IBAN" mandatory:"true"`
	// AlternativeAccounts []string `spayd:"ALT-ACC" format:"IBAN"`  // dunno how to implement it yet
	Amount        float64 `spayd:"AM" max_len:"10" mandatory:"true"`
	Currency      string  `spayd:"CC" len:"3"`
	Ref           string  `spayd:"RF" max_len:"16"`
	Name          string  `spayd:"RN" max_len:"35"`
	Date          string  `spayd:"DT" format:"date:YYYYMMDD", len:"8"`
	PaymentType   string  `spayd:"PT" len:"3"`
	Message       string  `spayd:"MSG" max_len:"60"`
	Notify        string  `spayd:"NT" len:"1"`
	NotifyAddress string  `spayd:"NTA" max_len:"320"`
	Url           string  `spayd:"X-URL" max_len:"140"`
	Vs            string  `spayd:"X-VS" max_len:"10" format:"numeric"`
}

func (spayd Spayd) Encode() ([]byte, error) {
	var (
		buffer = bytes.NewBuffer(make([]byte, 0, 512))
		str    string
	)
	buffer.WriteString("SPD*1.0") // SPAYD header
	rt := reflect.TypeOf(spayd)
	rv := reflect.ValueOf(spayd)
	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i)
		fv := rv.Field(i)

		if _, mandatory := ft.Tag.Lookup("mandatory"); mandatory {
			if fv.IsZero() {
				return nil, fmt.Errorf("Key %s is mandatory", ft.Name)
			}
		}
		// get string value for serialization
		switch fv.Kind() {
		case reflect.String:
			str = fv.String()
		case reflect.Float32, reflect.Float64:
			str = fmt.Sprintf("%.2f", fv.Float())
		default:
			panic(fmt.Sprintf("SPAYD encoding does not support %s type\n", ft.Type.String()))
		}
		// check constraints
		// check exact length
		if tag, defined := ft.Tag.Lookup("len"); defined {
			l, _ := strconv.Atoi(tag)
			if len(str) > 0 && len(str) != l {
				return nil, fmt.Errorf("%s has wrong length. Required length is %d yours is %d", ft.Name, l, len(str))
			}
		}
		// check maximal length
		if tag, defined := ft.Tag.Lookup("max_len"); defined {
			max_len, _ := strconv.Atoi(tag)
			if len(str) > max_len {
				return nil, fmt.Errorf("%s is too long. Maximal length is %d but yours is %d", ft.Name, max_len, len(str))
			}
		}
		// check format
		// if tag, defined := ft.Tag.Lookup("format"); defined {
		// 	// TODO
		// }

		if len(str) != 0 {
			buffer.WriteString("*")
			buffer.WriteString(ft.Tag.Get("spayd"))
			buffer.WriteString(":")
			buffer.WriteString(str)
		}
	}
	return buffer.Bytes(), nil
}

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
