FROM golang:1.10

# create a working directory
WORKDIR /go/src/github.com/ofonimefrancis/onefootball/

# add source code
ADD ./ /go/src/github.com/ofonimefrancis/onefootball/


RUN echo $GOPATH

RUN go get -d -v

# build main.go
RUN go build main.go

CMD ["./main"]