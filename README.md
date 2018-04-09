# go-jwt-rest-example
A simple implementation showing REST API secured with JWT in Golang
This is basic example of REST based API server in Golang

The example is very basic
It is an API which allows to perform CRUD operations on a `Product` database which is an SQLITE databse

The program also encorporates JWT authentication mechanism to secure the routes

The program uses ECHO library by Labstack - https://github.com/labstack/echo
The ECHO library is used for creating simple REST endpoints and securing them with JWT

The program also uses GORM - https://github.com/jinzhu/gorm
GORM is a very simple easy to wrapper(not fully ORM) for SQL database

The reason for using GORM is very simple. It allows to use mulitple databases at any given time without changing the existing code
It provides wrappers for popular SQL databases

The program runs on PORT :8080


To use it first use /login to get a JWT token
And the in subsequent requests use the generated token in Authorization header

Also please install required library befor compiling the program-

Use these commands to get all libraries-
go get -u github.com/labstack/echo/...
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/mattn/go-sqlite3
go get -u github.com/jinzhu/gorm

Also make sure that GCC is in your PATH variable if your using Windows
