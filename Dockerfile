FROM golang:alpine AS build
RUN apk add git
ADD . /go/src/PCCISO/
WORKDIR /go/src/PCCISO
RUN go build .

FROM alpine
WORKDIR /work
RUN mkdir /work/data
COPY ./web /work/web
COPY --from=build /go/src/PCCISO/PCCISO /bin/PCCISO
ENTRYPOINT [ "PCCISO" ]