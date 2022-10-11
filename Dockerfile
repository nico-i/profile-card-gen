# syntax=docker/dockerfile:1

# builder image
FROM surnet/alpine-wkhtmltopdf:3.8-0.12.5-full as builder

FROM golang:1.16-alpine

# Install dependencies for wkhtmltopdf
RUN  echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/main" > /etc/apk/repositories \
     && echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/community" >> /etc/apk/repositories \
     && apk update && apk add --no-cache \
      libstdc++ \
      libx11 \
      libxrender \
      libxext \
      libssl1.0 \
      ca-certificates \
      fontconfig \
      freetype \
      ttf-dejavu \
      ttf-droid \
      ttf-freefont \
      ttf-liberation \
      ttf-ubuntu-font-family \
    && apk add --no-cache --virtual .build-deps \
    \
    # Install microsoft fonts
    && fc-cache -f \
    \
    # Clean up when done
    && rm -rf /var/cache/apk/* \
    && rm -rf /tmp/* \
    && apk del .build-deps

COPY --from=builder /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=builder /bin/wkhtmltoimage /bin/wkhtmltoimage

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY public/ ./public/

RUN go build -o /docker-profile-card-gen

CMD [ "/docker-profile-card-gen" ]
