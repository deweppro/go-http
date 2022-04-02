package main

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/deweppro/go-http/pkg/httputil/enc"
	"github.com/deweppro/go-http/pkg/routes"
	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-http/servers/web"
	"github.com/deweppro/go-logger"
)

type Person struct {
	XMLName   xml.Name `xml:"person"`
	Id        int      `xml:"id,attr"`
	FirstName string   `xml:"name>first"`
	LastName  string   `xml:"name>last"`
	Age       int      `xml:"age"`
	Height    float32  `xml:"height,omitempty"`
	Married   bool
	Comment   string `xml:",comment"`
}

func main() {
	logger.Default().SetLevel(logger.LevelDebug)

	conf := servers.Config{Addr: ":8080"}
	route := routes.NewRouter()
	serv := web.New(conf, route, logger.Default())

	route.Route("/", func(w http.ResponseWriter, r *http.Request) {
		enc.Raw(w, []byte("Hello"))
	}, http.MethodGet, http.MethodPost)

	route.Route("/xml", func(w http.ResponseWriter, r *http.Request) {
		v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
		v.Comment = "Need more details."
		enc.XML(w, v)
	}, http.MethodGet, http.MethodPost)

	route.Route("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("*")
	}, http.MethodGet, http.MethodPost)

	route.Global(
		routes.RecoveryMiddleware(logger.Default()),
		routes.ThrottlingMiddleware(1000),
		routes.CORSMiddleware(routes.CORSConfig{
			Age:     100,
			Origin:  []string{"localhost"},
			Methods: []string{http.MethodGet, http.MethodPost},
			Headers: []string{"X-Token"},
		}),
		func(ctrlFunc routes.CtrlFunc) routes.CtrlFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Trace", "0000")
				ctrlFunc(w, r)
			}
		},
	)

	if err := serv.Up(); err != nil {
		panic(err)
	}

	<-time.After(60 * time.Second)

	if err := serv.Down(); err != nil {
		panic(err)
	}

	logger.Close()
}
