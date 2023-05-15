FROM docker.io/golang:1.20 as build

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY ./main.go ./main.go
COPY ./internal ./internal
COPY ./routers ./routers
COPY ./templates ./templates
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/idd4

FROM alpine as system
RUN addgroup idd4; \
    adduser idd4 -G idd4 -D  -h /home/b -s /bin/nologin;

FROM scratch

COPY --from=system /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=system /etc/nsswitch.conf /etc/nsswitch.conf
COPY --from=system /etc/passwd /etc/passwd

COPY --from=build /app/bin/idd4 /usr/bin/idd4
COPY --from=build /app/templates /templates

USER idd4

CMD ["idd4"]