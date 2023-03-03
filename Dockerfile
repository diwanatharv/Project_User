
#pull the latest image of golang
# syntax=docker/dockerfile:1

FROM golang:1.16-alpine
#creating the build directory
#ENTRYPOINT ["C:\Users\atharv.d\GolandProjects\golangapp\"]

 #cd C:\\User\\atharv.d\\GolandProjects
WORKDIR /app

EXPOSE 8000

COPY ../project_User%201/go.mod ./
COPY ../project_User%201/go.sum ./
RUN go mod download
#so we can pull any version of package from github
COPY *.go ./
#we require all these packages from github

#build ke ander jake ye load(clone) karlena
RUN go build -o /docker-gs-ping
CMD [ "/docker-gs-ping" ]
ENTRYPOINT ["/main.go"]
#same as jis port pe apis run ho rahi hai

#image load hone ke baad ye chale





