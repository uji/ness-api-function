FROM golang:1.15

COPY . /repo
WORKDIR /repo

ENV PORT 3000
ENV GITHUB_USER_NAME ""
ENV GITHUB_TOKEN ""
ENV GITHUB_EMAIL ""
ENV AWS_ACCESS_KEY_ID dammy
ENV AWS_SECRET_ACCESS_KEY dammy
ENV AWS_REGION us-east-1
ENV DB_ENDPOINT http://db-with-gui:8000

RUN go get github.com/golang/mock/mockgen@v1.4.4
RUN go get github.com/99designs/gqlgen
