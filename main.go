package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lucasb-eyer/go-colorful"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

type appContext struct {
	c *gpio.RgbLedDriver
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.h(ah.appContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func color(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	log.Printf("%s %s\n", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, indexHTML)
		return 200, nil
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
				err = a.c.SetRGB(byte(c.R*255), byte(c.G*255), byte(c.B*255))
				check(err)
				w.WriteHeader(http.StatusAccepted)
				return 202, nil
			}
			w.WriteHeader(http.StatusBadRequest)
			return 400, nil
		}
		fmt.Fprintf(w, indexHTML)
		return 200, nil
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, indexHTML)
		return 200, nil
	}
}

func main() {
	a := raspi.NewAdaptor()
	// These are the PHYSICAL pins here...
	// physical pin 12 == GPIO 4
	// physical pin 11 == GPIO 17
	// physical pin 7 == GPIO 18
	roy := gpio.NewRgbLedDriver(a, "12", "11", "7")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := roy.SetRGB(0, 0, 0)
		check(err)
		err = roy.Off()
		check(err)
		os.Exit(0)
	}()

	context := &appContext{roy}

	http.Handle("/", appHandler{context, color})
	port := "6060"

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
