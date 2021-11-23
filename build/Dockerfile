FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/openeoe-snapshot \
    USER_UID=1001 \
    USER_NAME=openeoe-snapshot

# install operator binary
COPY _output/bin/openeoe-snapshot ${OPERATOR}
COPY template /root/template

COPY bin /usr/local/bin
RUN /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER ${USER_UID}