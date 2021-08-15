## We'll choose the incredibly lightweight
## Go alpine image to work with
FROM golang:1.15-buster AS builder

## We create an /app directory in which
## we'll put all of our project code


## build ARGS

 


WORKDIR  /go/src/gopayments
COPY go.* ./
RUN go mod download


COPY . ./


## We want to build our application's binary executable

RUN  CGO_ENABLED=0 GOOS=linux go build -o main .

## the lightweight scratch image we'll
## run our application within
FROM alpine:latest AS production
## We have to copy the output from our
## builder stage to our production stage


ARG C2B_ONLINE_PARTY_B
ARG C2B_ONLINE_CHECKOUT_CALLBACK_URL
ARG MPESA_SHORTCODE
ARG MPESA_URL
ARG LNM_PASSKEY
ARG CONSUMER_KEY
ARG CONSUMER_SECRET
ARG MPESA_AUTH_URL
ARG DB_USER
ARG DB_PASS
ARG DB_NAME
ARG DB_HOST

ENV C2B_ONLINE_PARTY_B ${C2B_ONLINE_PARTY_B} 
ENV C2B_ONLINE_CHECKOUT_CALLBACK_URL ${C2B_ONLINE_CHECKOUT_CALLBACK_URL}
ENV MPESA_SHORTCODE ${MPESA_SHORTCODE}
ENV MPESA_URL ${MPESA_URL}
ENV LNM_PASSKEY ${LNM_PASSKEY} 
ENV CONSUMER_KEY ${CONSUMER_KEY}
ENV CONSUMER_SECRET ${CONSUMER_SECRET} 
ENV MPESA_AUTH_URL=${MPESA_AUTH_URL}
ENV DB_USER ${DB_USER} 
ENV DB_PASS ${DB_PASS}
ENV DB_NAME ${DB_NAME} 
ENV DB_HOST ${DB_HOST}

COPY --from=builder /go/src/gopayments .
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["./main"]