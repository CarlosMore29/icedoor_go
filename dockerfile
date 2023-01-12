FROM golang:1.19.0-alpine3.16 AS base
RUN apk add --no-cache tzdata
ENV TZ=America/Mexico_City
RUN cp /usr/share/zoneinfo/America/Mexico_City /etc/localtime
WORKDIR /app
ADD . .
## Add this go mod download command to pull in any dependencies
RUN apk update && apk add build-base git
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main .
FROM scratch AS prod
COPY --from=base /app/main ./main


## Our start command which kicks off
## our newly created binary executable
#CMD ["/app/main"]
CMD ["./main"]
