ARG GO_IMAGE=golang:1.21-alpine

# Set the Base Image using the argument set above
FROM $GO_IMAGE AS base
RUN apk --no-cache add dumb-init \
    && addgroup -g 1000 go \
    && adduser -D -u 1000 -G go go
# Install curl, bash, node, npm
RUN apk add --no-cache bash curl nodejs npm

RUN mkdir -p /home/go/app 
RUN chown go:go /home/go/app

WORKDIR /home/go/app
USER go

# For these to work you need to do two things on the dokku host:
# 1. dokku config:set <app-name> KEY=VALUE
# 2. dokku docker-options:add <app-name> build "--build-arg ARG_NAME=${KEY}
# Where key is the one you've set in step 1.
# https://dokku.com/docs/deployment/builders/dockerfiles/?h=dockerfile#build-time-configuration-variables
ARG DATABASE_URL
ARG PORT
ARG JWT_SECRET

# Copy the Go mod and sum files (these describe the app's dependencies) and download dependencies
FROM base AS dependencies
COPY --chown=go:go go.mod ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the entire source code from the current directory to the Working Directory inside the container
COPY --chown=go:go . .

# Build tailwind css
COPY --chown=go:go package.json ./
COPY --chown=go:go tailwind.config.js ./
RUN npm install

FROM dependencies AS build
RUN npm run build-css

# install templ
RUN npm run build-templ

# build js script
RUN npm run build-script


# Build server
RUN scripts/buildprod.sh

FROM base AS production

ENV HOST=0.0.0.0
ENV DATABASE_URL=${DATABASE_URL}
ENV PORT=${PORT}
ENV JWT_SECRET=${JWT_SECRET}

COPY --chown=go:go --from=build /home/go/app/main .

RUN mkdir -p public
RUN mkdir -p view

COPY --chown=go:go --from=build /home/go/app/public ./public/
COPY --chown=go:go --from=build /home/go/app/view/*_templ.go ./view/

EXPOSE ${PORT}

# Command to run when starting the container
CMD ["./main"]
