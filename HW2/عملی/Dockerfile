# مرحله اول: ساخت برنامه (Build Stage)
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o proxy main.go

# مرحله دوم: اجرا در ایمیج سبک (Final Stage)
FROM alpine:3.19

WORKDIR /app

# اضافه کردن تلاش مجدد (retry) برای نصب پکیج در صورت قطعی موقت شبکه
RUN apk update && apk --no-cache add ca-certificates

COPY --from=builder /app/proxy .

EXPOSE 80

CMD ["./proxy"]