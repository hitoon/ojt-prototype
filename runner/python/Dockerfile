FROM python:3.6 
# TODO: change base language

RUN apt-get update && apt-get install -y wget

RUN wget -q https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz

RUN tar -xf go1.10.3.linux-amd64.tar.gz

RUN mv go /usr/local

RUN rm go1.10.3.linux-amd64.tar.gz

ENV GOROOT /usr/local/go

ENV GOPATH /goworks

ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

RUN go get -u github.com/gin-gonic/gin

RUN mkdir $GOPATH/src/github.com/ojt 

WORKDIR $GOPATH/src/github.com/ojt

COPY . $GOPATH/src/github.com/ojt

EXPOSE 7070
# TODO: change

# Run app.py when the container launches
CMD ["go", "run", "runServer.go"]
