#!/bin/bash

git config --global user.name ${GITHUB_USER_NAME}
git config --global user.email ${GITHUB_EMAIL}
git config --global url."https://${GITHUB_USER_NAME}:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
git remote set-url origin https://github.com/uji/ness-api-function.git

go get github.com/99designs/gqlgen
go get -v github.com/rubenv/sql-migrate/...
go get github.com/golang/mock/mockgen@v1.4.4
