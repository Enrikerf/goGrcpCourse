# GRCP course in golang


[reference](https://www.udemy.com/course/grpc-golang/)

### install protobuf


    $ sudo apt-get install protobuf-compiler
    $ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}




### install evans:

    $ brew tap ktr0731/evans
    $ brew install evans

Run evans cli

    $ evans -p :portNumber -r 

Show the available tools:

    show package
    show services
    desc :requestTypeName or :responseType

Select a package

    package calculator

Select a service

    service calculator

Call to a function to run as a command

    call Sum

if you call a function biStream control+D to exit

### Mongo

setup in /docker/mongo/


### Grcp gateway for expose restApi

To expose the services as a rest api too you need install

* [grpc-gatewaty](https://github.com/grpc-ecosystem/grpc-gateway)

# Info

* [fork of grcp](https://github.com/gogo/protobuf)


[Mind the notes](/blog/README.md)