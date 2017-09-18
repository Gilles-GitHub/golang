# My Golang Rest API

Simple program to train myself to use Golang : https://golang.org/

 - Go 1.9
 - Spring Boot 1.5.7 (SOAP 1.1)

Documents
-------------

It handles a JSON request then query a Springboot SOAP WS to return a country informations as JSON.

> **Example:**

> GET http://localhost:12345/country/Spain will return
> 
```javascript
{
    "name": "Spain",
    "population": "46704314",
    "capital": "Madrid",
    "currency": "EUR"
}
```
----------

How to start
-------------

Run the Springboot jar
```bash
java -jar gs-producing-web-service-0.1.0.jar
```

Run the main go
```bash
go get sessiontechniquegolang.go
go run sessiontechniquegolang.go
```

Launch a request (with Postman for example) and you're good to go !!
