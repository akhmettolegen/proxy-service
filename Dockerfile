#
# Контейнер сборки
#
FROM golang:latest as builder

ENV CGO_ENABLED=0

COPY . /go/src/github.com/akhmettolegen/proxy-service
WORKDIR /go/src/github.com/akhmettolegen/proxy-service
RUN \
    version=git describe --abbrev=6 --always --tag; \
    echo "version=$version" && \
    cd cmd/app && \
    go build -a -tags proxy-service -installsuffix proxy-service -ldflags "-X main.version=${version} -s -w" -o /go/bin/proxy-service

#
# Контейнер рантайма
#
FROM scratch
COPY --from=builder /go/bin/proxy-service /bin/proxy-service

# копируем пользователя и группу из alpine
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group

ENTRYPOINT ["/bin/proxy-service"]