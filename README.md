# Go-Web

First golang web project. It's just a wall where you can add whatever you want. But first you should log in with your vk account to be able to create blog posts. Also on this page there is a search to find content you like. Simple general chat will let you communicate with online users.

# Stack

## Backend

- HTTP library [http](https://pkg.go.dev/net/http)
- URL router [mux](https://github.com/gorilla/mux)
- Cross Site Request Forgery middleware [CSRF](https://github.com/gorilla/csrf)
- Load ENV variables from .env file [godotenv](https://github.com/joho/godotenv)
- Secure cookie [securecookie](https://github.com/gorilla/securecookie)
- WebSocket protocol defined in RFC 6455 [websocket](https://github.com/gorilla/websocket)
- OAuth2.0 library [OAuth](https://golang.org/x/oauth2)
- Database [MongoDB](https://www.mongodb.com/)
- VK API library [vkapi](https://github.com/go-vk-api/vk)

## Frontend

- Server side templating [Go Templates](https://golang.org/pkg/text/template/)
- Bootstrap framework [Bootstrap](https://getbootstrap.com/)
- Javascript UI library [JQuery](https://jquery.com/)

## Project structure

![structure](structure.png)

## Run application

This project requires go v1.16+ to run the required services.

1. Get app from github:

```
git clone https://github.com/ineverbee/go-web.git
```

2. To run application and start server

```
go run .
```

3. Navigate to [page](http://localhost:8080/)

## App on Heroku

![wall-go.herokuapp.com](https://wall-go.herokuapp.com/)

