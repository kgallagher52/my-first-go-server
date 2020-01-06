# my-first-go-server
This is a simple Golang server with JWT auth. 


## Resources

[lectures](https://www.udemy.com/course/build-jwt-authenticated-restful-apis-with-golang/learn/lecture/)

[HTTP-ROUTER](https://github.com/gorilla/mux) - install: "go get -u github.com/gorilla/mux"

[JWT](https://github.com/dgrijalva/jwt-go) - install: "go get github.com/dgrijalva/jwt-go"

[pq](https://github.com/lib/pq) - install: "go get github.com/lib/pq"

[postman](https://www.getpostman.com/) - for testing endpoints

[golang-install](https://golang.org/doc/install)

[JWT-Offical-site](https://jwt.io)

[JWT-Info](https://tools.ietf.org/html/rfc7519)

### What are JWTs?

    - JWT stands for JSON Web Token

    - JWT is a means of exchanging information between two parties

    - Digitally signed

    - Structure of a JWT: {Base64 encoded Header}.{Base64 encoded Payload}.{Signature}

      - **header**: Algorithm & Token Type

      - **payload**: Carry claims 

        - Contains data such as User information token expiry etc..

        - Three types of claims: Registered, Public, and Private
      
      - **signature**: Computed from the Header, Payload and a Secret

        - An algorithm to generate the signature

        - Digitally signed using a secret string only known to the developer