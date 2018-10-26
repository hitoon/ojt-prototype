FROM golang:1.9

ENV GOROOT /usr/local/go

ENV GOPATH /goworks

ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

RUN go get -u github.com/gin-gonic/gin

RUN mkdir $GOPATH/src/github.com/ojt 

# Set the working directory to /app
WORKDIR $GOPATH/src/github.com/ojt

# Copy the current directory contents into the container at /app
COPY . $GOPATH/src/github.com/ojt


# Make port 80 available to the world outside this container
EXPOSE 8080

# Run app.py when the container launches
CMD ["go", "run", "webServer.go"]
