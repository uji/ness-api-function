FROM golang:1.15

COPY . /repo
WORKDIR /repo

ENV PORT 3000
ENV GITHUB_USER_NAME ""
ENV GITHUB_TOKEN ""
ENV GITHUB_EMAIL ""
ENV DB_HOST db
ENV DB_PORT 5432
ENV DB_NAME postgres
ENV DB_USER postgres
ENV DB_PASS password
