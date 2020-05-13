FROM golang:1.14-alpine AS build

RUN apk update \
    && apk --no-cache add ca-certificates
WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/trumail

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/trumail /bin/trumail
EXPOSE 8080
ENTRYPOINT ["/bin/trumail"]