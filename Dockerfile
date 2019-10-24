FROM gocd/gocd-server:v19.9.0

ARG UID

USER root

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go

USER go
