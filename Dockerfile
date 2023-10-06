# Use a minimal Linux-based image as the base image
FROM golang:1.21-alpine

# Create a directory to store your Go executable
RUN mkdir /app

# Copy your Go executable into the container
COPY iacsignarsrv_linux_amd64.sh /app/

# Copy your Go configuration file into the container
COPY signalRconfig.json /app/
# Set the working directory inside the container
WORKDIR /app

# Define the command to run your Go application
CMD ["./iacsignarsrv_linux_amd64.sh"]
