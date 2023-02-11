# get golang image for build as workspace
FROM golang:1.19 AS build

ARG GITHUB_TOKEN
RUN git config --global url.https://hanekawa_san:${GITHUB_TOKEN}@github.com/.insteadOf https://github.com/
RUN go env -w GOPRIVATE="github.com/kanji-team"
ENV PROJECT="auth"
# make build dir
RUN mkdir /${PROJECT}
WORKDIR /${PROJECT}
COPY go.mod go.sum ./

# download dependencies if go.sum changed
RUN go mod download
COPY . .

RUN make build

# create image with new binary
FROM multiarch/ubuntu-core:arm64-bionic AS deploy

ENV PROJECT="auth"
COPY --from=build /${PROJECT}/migrations /migrations
COPY --from=build /${PROJECT}/bin/${PROJECT} /${PROJECT}

CMD ["./auth"]