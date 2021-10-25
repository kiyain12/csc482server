# #get an image from docker 
# FROM golang:1.17.1

# ENV LOGGLY_TOKEN=b8f11a91-c1ab-4e66-ad6d-ae4eb5651dd9
# #create a new folder called app inside the image above
# RUN mkdir /app 
# #add whatever is in this folder then take it to the app folder inside the image
# ADD . /app
# #app is the working directory
# WORKDIR /app
# #build an executable from the main.go file 
# RUN go build -o main .
# #go inside the app directory and run the main excutable file 
# CMD ["/app/main"]



# Start the Go app build
FROM golang:latest AS build

# Copy source
WORKDIR /app
COPY . .

# Get required modules (assumes packages have been added to ./vendor)
RUN go mod download

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /

# Copy files from previous build container
COPY --from=build /app/main ./

# Add environment variables
# ENV ...
ENV LOGGLY_TOKEN=b8f11a91-c1ab-4e66-ad6d-ae4eb5651dd9

# Check results
RUN env && pwd && find .

# Start the application
CMD ["./main"]
