FROM gocd/gocd-server:v17.8.0

ARG UID

COPY ./scripts/shutdown-clean-goserver.sh /shutdown.sh

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go && \
    chmod +x /shutdown.sh