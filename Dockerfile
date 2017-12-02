FROM golang:1.9.2-alpine

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin
ENV WORKDIR /go/src/app
ENV BUILDDIR /app

# build directories
RUN mkdir $BUILDDIR
RUN mkdir $WORKDIR
ADD . $WORKDIR
WORKDIR $WORKDIR
# COPY Gopkg.toml Gopkg.lock ./

# Install dependencies
RUN apk --update add git

# Install dependency manager
RUN go-wrapper download -u github.com/golang/dep/cmd/dep \
    && go-wrapper install github.com/golang/dep/cmd/dep
    # && rm -rf /usr/lib/go /go/src /go/pkg /var/cache/*

# Get dependencies
RUN dep ensure -v
#RUN dep ensure

EXPOSE 8080

CMD ["go-wrapper", "run"]

# Build my app
# RUN go build -o $BUILDDIR/main .
# CMD ["$BUILDDIR/main"]