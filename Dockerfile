# get golang image for build as workspace
FROM golang:1.19 AS build

# make build dir
RUN mkdir /kanji-auth
WORKDIR /kanji-auth
COPY go.mod go.sum ./

# download dependencies if go.sum changed
RUN go mod download
COPY . .

RUN make build

# create image with new binary
FROM scratch AS deploy

COPY --from=build /kanji-auth/bin/kanji-auth /kanji-auth

CMD ["./kanji-auth"]