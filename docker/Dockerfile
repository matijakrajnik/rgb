FROM node:16 AS frontendBuilder

# set app work dir
WORKDIR /rgb

# copy assets files to the container
COPY assets/ .

# set assets/ as work dir to build frontend static files
WORKDIR /rgb/assets
RUN npm install
RUN npm run build

FROM golang:1.16.3 AS backendBuilder

# set app work dir
WORKDIR /go/src/rgb

# copy all files to the container
COPY . .

# build app executable
RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/rgb/rgb cmd/rgb/main.go

# build migrations executable
RUN CGO_ENABLED=0 GOOS=linux go build -o migrations/migrations migrations/*.go

FROM alpine:3.14

# Create a group and user deploy
RUN addgroup -S deploy && adduser -S deploy -G deploy

ARG ROOT_DIR=/home/deploy/rgb

WORKDIR ${ROOT_DIR}

RUN chown deploy:deploy ${ROOT_DIR}

# copy static assets file from frontend build
COPY --from=frontendBuilder --chown=deploy:deploy /rgb/build ./assets/build

# copy app and migrations executables from backend builder
COPY --from=backendBuilder --chown=deploy:deploy /go/src/rgb/migrations/migrations ./migrations/
COPY --from=backendBuilder --chown=deploy:deploy /go/src/rgb/cmd/rgb/rgb .

# set user deploy as current user
USER deploy

# start app
CMD [ "./rgb", "-env", "prod" ]
