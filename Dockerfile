FROM gocd/gocd-server:v18.7.0

ARG UID

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go
