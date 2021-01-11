# Image resizer

Image resizer is a test project for developers hwo would like ot work in the Mad Devs company on a junior Golang developer position. 

## Instalation 
To install this program on your local machine create a new direcotory and type in it:
```

git init

git clone https://github.com/Padrition/go-resize-rest-api.git

```

## Usage 
First start the server for beeing able to use the program.
Go to the go-resize-rest-api directory and type :
```

go run main.go

```
This will start the server on 8080 prot.

The server have three endpoints. One is "/" which returns an html where you can call others endpoints and sortoff play with the program.

The second end point is "/upload" which will parse a given image of jpeg, png or gif formats and upload them to your local machine, if a file you uploaded is vaolationg the format.

The final endpint is "/resize"  which will parse a given jpeg, png or gif file, resize it and send it back to user. Althou it doesn't take any size arguments in url you can parse them in HTTP POST request body. Or just use html file by going to:

```
localhost:8080/
```

When you will be on your beloved localhost:8080 you will have two options:
1)Upload a file to the server. Just upload your image and press upload button.
2)Resize an image with given width and height. Just upload your image set wannable dimensions and press resize button.