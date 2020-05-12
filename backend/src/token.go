package main

import (
	"encoding/base64"
	"fmt"
)

func main() {


	//token:="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJyb290IiwiZXhwIjo3MjAwLCJpc3MiOiJCYWNrZW5kIn0.sJqcHanOdpqvKHGm1c-YJzChf6odbV9Tnjl7qouJGzo"
	header := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	payload := "eyJhdWQiOiJyb290IiwiZXhwIjo3MjAwLCJpc3MiOiJCYWNrZW5kIn0"
	signature := "sJqcHanOdpqvKHGm1c-YJzChf6odbV9Tnjl7qouJGzo"

	decodedHeader, err := base64.URLEncoding.DecodeString(header)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println("header:", string(decodedHeader))

	decodedPayload, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println("payload:", string(decodedPayload))

	decodedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println("signature:", string(decodedSignature))

}
