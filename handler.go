package main

import (
	"net/http"
	"context"
	"log"
	"errors"
	"syscall"
	"time"
	"io/ioutil"
	"html/template"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/csrf"
	"github.com/go-vk-api/vk"
	"golang.org/x/oauth2"
	vkAuth "golang.org/x/oauth2/vk"
)

// Hash keys should be at least 32 bytes long
var hashKey = []byte(goDotEnvVariable("HASH_KEY"))
// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte(goDotEnvVariable("BLOCK_KEY"))
var s = securecookie.New(hashKey, blockKey)

var conf = &oauth2.Config{
    ClientID:     goDotEnvVariable("CLIENT_ID"),
    ClientSecret: goDotEnvVariable("CLIENT_SECRET"),
    RedirectURL:  "http://localhost:8080/auth",
    Scopes:       []string{},
    Endpoint:     vkAuth.Endpoint,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/home.html", "templates/base.html")
	
	if err != nil {
		log.Fatal(err)
	}

	searchText := r.URL.Query().Get("search_text")

	if cookie, err := r.Cookie("user"); err != nil {
    	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

    	posts := findPosts(searchText)

		var users []User

		for _, v := range posts {
			usr := User{ID: v.User}
			users = append(users, usr.findUser())
		}

		err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
	        csrf.TemplateTag: csrf.TemplateField(r),
	        "URL": url,
	        "Users": users,
	        "Posts": posts,
	        "NoPosts": "There are no posts yet!",
	    })

		if errors.Is(err, syscall.EPIPE) {
            // just ignore.
            return
        } else if err != nil {
            log.Fatal(err)
        }
	} else {
		
		var value int64
		if err = s.Decode("cookie-user-auth", cookie.Value, &value); err != nil {
			log.Fatal(err)
		}

		user := User{VK_ID: value}
		user = user.findUser()

		if r.Method == "POST" {
			var img string = ""
			// Parse our multipart form, 10 << 20 specifies a maximum
		    // upload of 10 MB files.
		    r.ParseMultipartForm(10 << 20)
		    // FormFile returns the first file for the given key
		    // it also returns the FileHeader so we can get the Filename,
		    // the Header and the size of the file
		    file, _, err := r.FormFile("image")
		    if err == nil {
		        // Create a temporary file within our temp-images directory that follows
			    // a particular naming pattern
			    tempFile, err := ioutil.TempFile("static/temp-images", "upload-*.png")
			    if err != nil {
			        log.Fatal(err)
			    }
			    defer tempFile.Close()

			    // read all of the contents of our uploaded file into a
			    // byte array
			    fileBytes, err := ioutil.ReadAll(file)
			    if err != nil {
			        log.Fatal(err)
			    }
			    // write this byte array to our temporary file
			    tempFile.Write(fileBytes)
			    img = tempFile.Name()
			    defer file.Close()
		    }
		    
		    
		    post := Post{
		    	User: user.ID,
		    	Text: r.FormValue("text"),
		    	Image: img,
		    	Link: r.FormValue("link"),
		    }
		    
		    post.createPost()
		    http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}

		posts := findPosts(searchText)

		var users []User

		for _, v := range posts {
			usr := User{ID: v.User}
			users = append(users, usr.findUser())
		}

		err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
	        csrf.TemplateTag: csrf.TemplateField(r),
	        "U": user,
	        "Users": users,
	        "Posts": posts,
	        "NoPosts": "There are no posts yet!",
	    })

		if errors.Is(err, syscall.EPIPE) {
            // just ignore.
            return
        } else if err != nil {
            log.Fatal(err)
        }
	}
}

func authPage(w http.ResponseWriter, r *http.Request) {

	// получаем код от API VK из квери стринга
    authCode := r.URL.Query()["code"]

    if authCode != nil {
    	// меняем код на access токен
	    tok, err := conf.Exchange(context.Background(), authCode[0])

	    if err != nil {
	        log.Fatal(err)
	    }
	    // создаем клиент для получения данных из API VK
	    client, err := vk.NewClientWithOptions(vk.WithToken(tok.AccessToken))

	    if err != nil {
	        log.Fatal(err)
	    }

	    user := getCurrentUser(client)

	    user.createUser()

	    r.URL.Query().Del("code")

	    expiration := time.Now().Add(365 * 24 * time.Hour)
		value := user.VK_ID
		encoded, err := s.Encode("cookie-user-auth", value)

		if err != nil {
			log.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:  "user",
			Value: encoded,
			HttpOnly: true,
			Expires: expiration,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func getCurrentUser(api *vk.Client) User {
    var users []User
  
    api.CallMethod("users.get", vk.RequestParams{
       "fields": "id,first_name,last_name,photo_400_orig,city",
    }, &users)
  
    return users[0]
}
