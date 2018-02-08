FROM golang:1.9.2-alpine

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin
ENV WORKDIR /go/src/github.com/ld100/goblet
ENV BUILDDIR /app
ENV PATH="${BUILDDIR}:${PATH}"

# Install dependencies
RUN apk --update add git

# Install Golang's dependency manager
RUN go-wrapper download -u github.com/golang/dep/cmd/dep \
    && go-wrapper install github.com/golang/dep/cmd/dep
    # && rm -rf /usr/lib/go /go/src /go/pkg /var/cache/*


# Set project directories
RUN mkdir -p $BUILDDIR
RUN mkdir -p $WORKDIR
ADD . $WORKDIR/
WORKDIR $WORKDIR

# Copies the Gopkg.toml and Gopkg.lock to WORKDIR
COPY Gopkg.toml Gopkg.lock ./

# Get dependencies
# install the dependencies without checking for go code
RUN dep ensure -vendor-only
#RUN dep ensure

EXPOSE 8080

#CMD ["go-wrapper", "run"]

# Build my app
#RUN go build -o $BUILDDIR/goblet .
RUN go build -o $BUILDDIR/goblet main.go
CMD ["$BUILDDIR/goblet", "serve"]
