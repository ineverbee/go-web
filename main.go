package main

import (
	"log"
	"net/http"
    "fmt"
    "strings"
    "os"

    "github.com/gorilla/csrf"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

    // load .env file
    err := godotenv.Load(".env")

    if err != nil {
      log.Fatalf("Error loading .env file")
    }

    return os.Getenv(key)
}

func getPort() string {
  
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	return port
  
}

var dbUsername = goDotEnvVariable("DB_USERNAME")
var dbPassword = goDotEnvVariable("DB_PASSWORD")
//URI used in mongo.go
var URI = fmt.Sprintf("mongodb+srv://%s:%s@cluster0.z2tkb.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", dbUsername, dbPassword)

func main() {
	router := mux.NewRouter()

    // Create a file server which serves files out of the "./static" directory.
    // Note that the path given to the http.Dir function is relative to the project
    // directory root.
    fileServer := http.FileServer(http.Dir("./static/"))

    // Use the router.Handle() function to register the file server as the handler for
    // all URL paths that start with "/static/". For matching paths, we strip the
    // "/static" prefix before the request reaches the file server.
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
    router.HandleFunc("/", homePage)
    router.HandleFunc("/ws", wsConnections)
    router.HandleFunc("/auth", authPage)
    router.HandleFunc("/sign_out", func(w http.ResponseWriter, r *http.Request) {
        c := http.Cookie{
                Name:   "user",
                MaxAge: -1}
        http.SetCookie(w, &c)
        http.Redirect(w, r, "/", http.StatusMovedPermanently)
    })
    router.HandleFunc("/delete_post", func(w http.ResponseWriter, r *http.Request) {
        postID := r.URL.Query().Get("post_id")
        postID = strings.TrimLeft(strings.TrimRight(postID,`")`),`ObjectID("`)
        deletePost(postID)
        imagePath := r.URL.Query().Get("delete_image")
        if imagePath != "" {
            err := os.Remove(imagePath)

            if err != nil { log.Fatal(err) }
        }
        http.Redirect(w, r, "/", http.StatusMovedPermanently)
    })

    go wsMessages()

    csrf.Secure(true)
	CSRF := csrf.Protect([]byte(goDotEnvVariable("CSRF_KEY")))
	
	port := getPort()

    log.Println("Starting server on " + fmt.Sprintf(":%s", port))
    err := http.ListenAndServe(fmt.Sprintf(":%s", port), CSRF(router))
    log.Fatal(err)
}