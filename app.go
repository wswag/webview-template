package main

import (
	"net/http"
	"errors"
	"github.com/zserge/webview"
	"github.com/gorilla/mux"
	"controllers"
	"fmt"
)

var mainViewController *controllers.MainViewController

// ConfigureWindow should return an initialized window.Settings object to initialize
// the webview window
func ConfigureWindow() webview.Settings {
	return webview.Settings{
		Width:  800,
		Height: 600,
		Resizable: true,
		Title:  "Hello WebView!",
		Debug: true,
	}
}

// SetupServer should initialize server routes etc.
// note that there will also be a route to a http.FileServer be added
// to '/' to provide asset content
// an example usage is adding an api route to provide local hdd file content
// As the view has not yet been initialized, you may not call any
// code injection funcs here
func SetupServer(router *mux.Router) {
	
	mainViewController = &controllers.MainViewController {
		Title: "Hello Go Controller!",
		Answer: 42,
	}

	// add some custom routes providing queryable data beside js bindings
	router.PathPrefix("/local/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(404) })

	// do some stuff maybe
	fmt.Println("Server controllers initialized")
}

// SetupView might inject additional code on startup and/or setup the controllers.
// Note that code injection is performed after the page has already been loaded
func SetupView() {
	
}

// QueryService should return a service interface identified by the given name which
// then will be injected into javascript and be available as a js object through this name.
// Note that all properties of this struct will be accessible under name.data in js!
// an injection is queried from HTML content by adding a <script> tag pointing to
// the inject API route, eg:
//
//		<script src="/inject/model.js"></script>
//
func QueryService(name string) (interface{}, error) {
	// map your services here
	switch (name) {
		case "model": return mainViewController, nil
		default: return nil, errors.New("Service not available")
	}
}