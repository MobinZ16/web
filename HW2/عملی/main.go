package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// خواندن تنظیمات از متغیرهای محیطی با مقادیر پیش‌فرض
	targetURLStr := os.Getenv("TARGET_URL")
	if targetURLStr == "" {
		targetURLStr = "http://localhost:8081"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_super_secret_key_123456"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// پارس کردن آدرس سرویس مقصد (Legacy Backend)
	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	// ساخت پروکسی معکوس
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// تعریف میدل‌ور و هندلر اصلی پروکسی
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ۱. بررسی وجود هدر Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid token format", http.StatusUnauthorized)
			return
		}

		// استخراج توکن
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// ۲. اعتبارسنجی توکن JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// بررسی الگوریتم امضا (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, fmt.Sprintf("Unauthorized: Invalid token (%v)", err), http.StatusUnauthorized)
			return
		}

		// ۳. هدایت درخواست به سرویس داخلی پس از موفقیت احراز هویت
		log.Printf("[PROXY] Authorized request to: %s", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	// راه‌اندازی سرور روی پورت مشخص‌شده
	log.Printf("Starting JWT Reverse Proxy on port %s, forwarding to %s...", port, targetURLStr)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}