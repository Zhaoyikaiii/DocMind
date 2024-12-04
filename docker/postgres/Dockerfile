FROM postgres:14.0

COPY ./init.sql /docker-entrypoint-initdb.d/

ENV TZ=Asia/Shanghai

ENV LANG en_US.utf8