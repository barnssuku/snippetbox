package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the 
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a dafault value of 
	// ":4000" and some short help text explaining what the flag controls. The 
	// value of the flag will be stored in the addr variable at runtine.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always containthe defualt value of ":4000". If any errors are
	// encountered during parseing the application will be terminated.
	flag.Parse()

	// Now log.New() to create a logger for writing information messages. This takes 
	// three paremeters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional infomation to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error maessages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies.
	app := &application {
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Swap the route declarations to use the application strut's methods as the 
	// handler functions.
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files wout of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handler() function to register the file server as the handler for 
	// all URL paths that start with "/static/". For matching path, we strip the 
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Initialize a new http.Server struct. We se the Addr and Hnadler fields
	// so that the server use the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in 
	// the event of any problems.
	srv := &http.Server {
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %v", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}