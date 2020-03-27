# go-SPAYD

Go implementation of Short PAYment Descriptor - a payment QRcode generator.

`go get github.com/katomaso/go-spayd`

## Usage

`go-spayd` is distributed as a single binary and is intended to be used as a
micro-service. It listens on configurable address and port, waiting for POST
requests with JSON describing intended payment by default on port `8484`.

Expected JSON does not use SPAYD keys but uses whole words instead. Is it a good
decision? Maybe.

```JSON
spayd = {
	"account":      "string", // spayd:"ACC" max_len:"46" format:"IBAN" mandatory:"true"`
	"amount":       "float",  // spayd:"AM" max_len:"10" precision:"2" mandatory:"true"`
	"currency":     "string", // spayd:"CC" len:"3"`
	"ref":          "string", // spayd:"RF" max_len:"16"`
	"name":         "string", // spayd:"RN" max_len:"35"`
	"date":         "string", // spayd:"DT" format:"date:YYYYMMDD", len:"8"`
	"paymentType":  "string", // spayd:"PT" len:"3"`
	"message":      "string", // spayd:"MSG" max_len:"60"`
	"notify":       "string", // spayd:"NT" len:"1"`
	"notifyAddress":"string", // spayd:"NTA" max_len:"320"`
	"url":          "string", // spayd:"X-URL" max_len:"140"`
	"vs":           "string", // spayd:"X-VS" max_len:"10" format:"numeric"`
}
```

You can simply test generator wirh `curl`

```sh
curl  -d '{"account":"123456", "amount":12.50}' -o qrcode.png localhost:8484
```

Embedding into your web page could be done as simply as

```html
<div id="qrcode"></div>

<script type="text/javascript" async defer>
  fetch("/spayd", {
    body: JSON.stringify({account: 12345678, amount: 150, message: "{{.Person.Email}}"})
  }).then(
    function(respose) {
      if (response.ok) {
        var img = document.createElement("img");
        img.src="data:image/png;base64,"+btoa.encode(response.Body.arrayBuffer())
        document.queryElement("div#qrcode").appendChild(img);
      } else {
        console.log(response.body.text())
      }
    }
  );
</script>
```

## TODO

* Implement all keys - currently only the most used are implemented.
* Implement format checkers - currently only `len`, `max_len` and `mandatory` constraints are implemented