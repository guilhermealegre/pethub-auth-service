FROM 176300518568.dkr.ecr.eu-west-1.amazonaws.com/baseimages:golang1.20.4-alpine3.17 as builder

COPY . /go/src/be-user-service
WORKDIR /go/src/be-user-service

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o be_user cmd/server/main.go

RUN chown root:root be_user
RUN chown 755 be_user

FROM 176300518568.dkr.ecr.eu-west-1.amazonaws.com/baseimages:alpine3.14

COPY --from=builder --chown=root:root /go/src/be-user-service/be_user .

RUN apk --no-cache add ca-certificates \
    curl \
    bash

RUN apk --update --no-cache add python3 py3-pip && \
    pip3 install awscli

RUN mkdir -p migrations
COPY migrations migrations
COPY internal/swagger/docs internal/swagger/docs
COPY build/bash-multi.entrypoint.sh /entrypoint.sh


RUN chmod +x /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]
CMD ["./be_user"]
