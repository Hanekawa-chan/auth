# get golang image for build as workspace
FROM golang:1.19 AS build

ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
ENV GOPRIVATE="github.com/Hanekawa-chan"
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