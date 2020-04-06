# go-SPAYD

Go implementation of Short PAYment Descriptor - a payment QRcode generator.

`go get github.com/katomaso/go-spayd`

## Usage

`go-spayd` is distributed as a single binary and is intended to be used as a
micro-service. It listens on configurable address and port, waiting for POST
requests with JSON describing intended payment by default on port `8484`.

JSON uses whole words instead of SPAYD keys. Is it a good decision? Maybe. Here
is the expected keys with their expected types.

```
spayd = {
  Account       string  `spayd:"ACC" max_len:"46" format:"IBAN" mandatory:"true"`
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
  KS            string  `spayd:"X-KS" max_len:"10" format:"numeric"`
  SS            string  `spayd:"X-SS" max_len:"10" format:"numeric"`
  VS            string  `spayd:"X-VS" max_len:"10" format:"numeric"`
}
```

You can simply test generator wirh `curl`

```sh
curl  -d '{"account":"CZ7562106701002206308683", "amount":125.50}' -o qrcode.png localhost:8484
```

![result](https://raw.githubusercontent.com/katomaso/go-spayd/master/qrcode.png "Resulting QRcode")
Embedding into your web page could be done as simply as

```html
<div id="qrcode"></div>

<script type="text/javascript" async defer>
  fetch("/spayd", {
    method: 'POST',
    body: JSON.stringify({account: "CZ7562106701002206308683", amount: 150, message: "{{.Person.Email}}"})
  }).then(
    function(response) {
      if (response.ok) {
        return response.blob();
      } else {
        throw response.text();
      }
    }
  ).then(function(blob) {
    // easy converting blob to base64 in javascript
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
        reader.onloadend = () => resolve(reader.result);
        reader.onerror = reject;
      reader.readAsDataURL(blob);
    });
  })
  .then(function(encoded) {
    var img = document.createElement("img");
    img.src=encoded;
    document.querySelector("div#qrcode").appendChild(img);
  })
  .catch(function(reason) {
    console.log(reason);
  });
</script>
```

## TODO

* Implement all keys - currently only the most used are implemented.