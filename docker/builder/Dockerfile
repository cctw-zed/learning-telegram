# This Dockerfile creates a custom builder image with all the necessary
# C build tools pre-installed to speed up the main application build.
FROM golang:1.23-alpine

# Install the C build tools needed for CGo.
# This is the slow step we want to run only once when building this base image.
RUN apk add --no-cache build-base

# Set a default Go proxy for convenience, so it doesn't need to be in the app's Dockerfile.
ENV GOPROXY=https://goproxy.cn,direct 