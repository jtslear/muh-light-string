package main

import (
	"io"
	"log"
	"net/http"

	"github.com/lucasb-eyer/go-colorful"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

var indexHTML = `
<html>
  <body>
    <form action="/" method="post">
      <input type="color" id="colorChoice" name="colorChoice" value="#aabbcc">
      <input type="submit">
    </form>
    <script>
      var color = document.getElementById("colorChoice").value;
      console.log("Color" + color);
    </script>
  </body>
</html>
`

func testFunction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w, indexHTML)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			if val, ok := r.PostForm["colorChoice"]; ok {
				log.Printf("Input: %s\n", val[0])
				c, err := colorful.Hex(val[0])
				if err != nil {
					log.Printf("Unable to grasp the hex color: %v\n", err)
				}
				log.Printf("Color: R:%d G:%d B:%d\n", int(c.R*255), int(c.G*255), int(c.B*255))
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
		io.WriteString(w, indexHtml)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, indexHtml)
	}
	log.Printf("%s %s\n", r.Method, r.URL.Path)
}

func main() {
	a := raspi.NewAdaptor()
	// These are the PHYSICAL pins here...
	roy := gpio.NewRgbLedDriver(a, "12", "11", "7")
	if roy.State() {
		log.Printf("It would appear on\n")
	} else {
		log.Printf("It is not on\n")
	}
	err := roy.On()
	if err != nil {
		log.Printf("Problem turning on the LED's: %v\n", err)
		panic(err)
	}
	if roy.State() {
		log.Printf("It would appear on\n")
	} else {
		log.Printf("It is not on\n")
	}
	http.HandleFunc("/", testFunction)
	port := "6060"

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
