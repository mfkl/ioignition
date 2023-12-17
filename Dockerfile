FROM alpine:3.14

# Set the Current Working Directory inside the container
WORKDIR /app

ARG DB_URL
ARG JWT_SECRET
ARG REDIS_URL
ARG IPINFO_TOKEN

# Copy binary & assets folder
COPY main .
COPY public public

RUN echo "DB_URL=${DB_URL}" >> .env
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
RUN echo "REDIS_URL=${REDIS_URL}" >> .env
RUN echo "IPINFO_TOKEN=${IPINFO_TOKEN}" >> .env

ENV HOST=0.0.0.0

EXPOSE 8080

CMD ["./main"]
