package main

import (
	"io"
	"log"
	"net/http"
	"github.com/lucasb-eyer/go-colorful"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

var indexHtml string = `
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
		io.WriteString(w, indexHtml)
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
	log.Println("Initting GPIO Pin 12...")
	seven := gpioreg.ByName("12")
	if seven == nil {
		log.Fatal("Failed to init GPIO12")
	}
	log.Printf("%s: %s\n", seven, seven.Function())
	http.HandleFunc("/", testFunction)
	port := "6060"

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
