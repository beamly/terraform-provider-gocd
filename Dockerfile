FROM gocd/gocd-server:v19.5.0

ARG UID

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go
