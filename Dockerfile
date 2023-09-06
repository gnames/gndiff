FROM alpine:3.14

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2021-12-27

WORKDIR /bin

COPY ./gndiff /bin

ENTRYPOINT [ "gndiff" ]
