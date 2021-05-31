FROM node:alpine AS front

WORKDIR /app

COPY build_front.sh .
COPY front ./front/
RUN sh build_front.sh

FROM golang:alpine as back

WORKDIR /app
RUN mkdir -p ./front/dist
COPY --from=front /app/front/dist/* ./front/dist/
COPY go.sum go.mod main.go ./
COPY ./internal/ ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o simtris .

FROM scratch  
COPY --from=back /app/simtris /bin/

ENTRYPOINT [ "/bin/simtris" ]