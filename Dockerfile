# base image
FROM node:12.16.0 as builder
# set working directory
WORKDIR /app
RUN git clone https://github.com/RobusK/AngularCtaUi.git .

ENV PATH /app/node_modules/.bin:$PATH

RUN npm install
RUN npm run build -- --prod --aot


FROM golang
ADD . /go/src/GoCtaApi

RUN go get "github.com/graphql-go/graphql"
RUN go get "github.com/graphql-go/handler"
RUN go get "pault.ag/go/haversine"

WORKDIR /go/src/GoCtaApi
RUN go build
COPY --from=builder /app/dist/CtaApiUi ./static
ENTRYPOINT ./GoCtaApi
EXPOSE 80
