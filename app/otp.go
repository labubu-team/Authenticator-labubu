package app

import (
	"fmt"
	"github.com/pquerna/otp/totp"
	"regexp"
	"strings"
	"time"
)

func isValidBase32(secret string) bool {
	// Regex secret base32
	base32Regex := regexp.MustCompile(`^[A-Z2-7]+$`)
	return base32Regex.MatchString(secret)
}

func extractSecretFromURL(url string) (string, error) {
	otpauthRegex := regexp.MustCompile(`^otpauth:\/\/(totp|hotp)\/[^?]+\?secret=([a-zA-Z0-9%]+)(?:&|$)`)
	matches := otpauthRegex.FindStringSubmatch(url)
	fmt.Println("matches", matches) // In kết quả ra để kiểm tra

	if len(matches) < 3 {
		return "", fmt.Errorf("không tìm thấy secret trong URL")
	}

	secret := matches[2]

	fmt.Println("secret:", secret)

	if isValidBase32(secret) {
		return secret, nil
	}
	return "", fmt.Errorf("secret không hợp lệ")
}

func GenerateOTP(secret string) (string, string, int64, error) {
	fmt.Println("otp", secret)
	// gen OTP for now time
	otp, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", "", 0, err
	}

	// default period of time
	timeStep := 30

	// get now time
	currentTime := time.Now().Unix()

	// calculator time compared to present
	elapsed := currentTime % int64(timeStep)

	// Calculate remaining time in current cycle
	timeLeft := int64(timeStep) - elapsed

	return otp, secret, timeLeft, nil
}

func generateOTPInputHandler(input string) (string, string, int64, error) {
	var secret string

	// check secret
	if isValidBase32(input) {
		secret = input
	} else if strings.HasPrefix(input, "otpauth://") {
		extractedSecret, err := extractSecretFromURL(input)
		fmt.Println("extractedSecret", extractedSecret)
		if err != nil {
			fmt.Println("Lỗi khi trích xuất secret từ URL:", err)
			return "", "", 0, err
		}
		secret = extractedSecret
	} else {
		fmt.Println("Đầu vào không hợp lệ:", input)
		return "", "", 0, fmt.Errorf("đầu vào không hợp lệ: %v", input)
	}
	fmt.Println("otp", secret)
	// Gen OTP
	otp, otpURL, timestamp, err := GenerateOTP(secret)
	if err != nil {
		return "", "", 0, err
	}

	return otp, otpURL, timestamp, nil
}

func startOTPRoutine(secret string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		otp, _, timeLeft, err := GenerateOTP(secret)
		if err != nil {
			fmt.Printf("Lỗi khi tạo mã OTP: %v\n", err)
			return
		}

		fmt.Printf("Mã OTP: %s\n", otp)
		fmt.Printf("Thời gian còn lại trong chu kỳ 30 giây: %d giây\n", timeLeft)

		<-ticker.C
	}
}
