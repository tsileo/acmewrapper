# Example

This is a very basic server which shows how to use acmewrapper

```
go build
```

will give you the executable

## Running

The server will serve the directory you specify as its last argument:

```
./example -accept -host=:8443 mywebsite.com www.mywebsite.com ./www
```

The above will run a server at port 8443 (which 443 is assumed to forward to), with certs for example.com and www.example.com.

It will serve the www directory.

Note that when testing, you should use the `-test` flag to make sure you are not generating valid certificates (you might run into the limits accidentally)
