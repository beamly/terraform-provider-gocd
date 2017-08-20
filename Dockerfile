FROM gocd/gocd-server:v17.7.0

COPY ./scripts/shutdown-clean-goserver.sh /shutdown.sh

RUN chmod +x /shutdown.sh