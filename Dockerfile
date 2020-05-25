FROM golang:1.13

WORKDIR /go/src/app
COPY . .

ENV DEBUG_PATH="/var/log/"
ENV PAGE_TITLE="My logs page"

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3333

CMD ["app"]