package main

import (
	"controllers"
	"fmt"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/wswag/webview"
)

var browserContext webview.WebView

func main() {
	settings := ConfigureWindow()
	settings.URL = startServer()
	w := webview.New(settings)
	defer w.Exit()

	browserContext = w

	SetupView()

	w.Run()
}

func startServer() string {
	// start tcp listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	// run http listener
	go func() {
		defer ln.Close()

		router := mux.NewRouter()
		router.Use(LogRequestsMiddleware)

		// default routes
		// add the inject api
		router.PathPrefix("/inject/").HandlerFunc(injectServiceHandler)

		// add user defined api routes etc.
		SetupServer(router)

		// add the static file server
		router.PathPrefix("/").Handler(http.FileServer(assetFS()))

		http.Handle("/", router)
		log.Fatal(http.Serve(ln, nil))
	}()

	// return URL to newly created app server
	return "http://" + ln.Addr().String()
}

// injectServiceHandler: enriched bind interface
// of webview to create the binding funcs return the js code instead
// of evaluating it directly -> done
func injectServiceHandler(w http.ResponseWriter, r *http.Request) {
	name := filepath.Base(r.URL.Path)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	fmt.Println("Injecting a service: " + name)
	w.Header().Add("Content-Type", mime.TypeByExtension(ext))

	// query the service and bind it in query mode
	service, err := QueryService(name)
	if err != nil {
		log.Fatalln("Service-Query of " + name + " threw an error: " + err.Error())
	}
	sync, js, err := browserContext.BindQuery(name, service)

	// check synchronizer interface and attach sync function
	syncer, ok := service.(controllers.Synchronizer)
	if ok {
		syncer.SetSyncFunction(func() { browserContext.Dispatch(sync) })
	}

	if err != nil {
		fmt.Println("Injection failed: " + err.Error())
	} else {
		w.Write([]byte(js))
		fmt.Println("Injection succeeded.")
	}
}

// InjectJs injects javascript code into the browser context
func InjectJs(code ...string) {
	log.Println("Injecting Js Code..")
	browserContext.Dispatch(func() {
		for _, c := range code {
			browserContext.Eval(c)
		}
	})
}

// InjectInterface binds a go object to the current browser context
func InjectInterface(objname string, itf interface{}) (err error) {
	browserContext.Dispatch(func() {
		_, err = browserContext.Bind(objname, itf)
	})
	return
}

// LogRequestsMiddleware logs every passed request to console by 'log' package
var LogRequestsMiddleware = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// measure time
		t0 := time.Now()
		next.ServeHTTP(w, r) //proceed in the middleware chain!
		t1 := time.Now()

		// log url path for now
		log.Print(r.Method + " request to: " + r.URL.Path + " from " + r.RemoteAddr + " (handled in " + fmt.Sprint(t1.Sub(t0)) + ")")
	})
}
