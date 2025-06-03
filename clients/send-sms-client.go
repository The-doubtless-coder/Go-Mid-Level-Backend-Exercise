package clients

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendSMS(to, message string) error {

	era := godotenv.Load(".env")
	if era != nil {
		log.Fatal("Error loading .env file:: Using os files instead")
	}

	apiKey := os.Getenv("AFRICASTALKING_API_KEY")
	username := os.Getenv("AFRICASTALKING_USERNAME")

	data := url.Values{}
	data.Set("username", username)
	data.Set("to", to)
	data.Set("message", message)

	req, err := http.NewRequest("POST", "https://api.sandbox.africastalking.com/version1/messaging", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("apiKey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func SendSMSAsync(to, message string) {
	go func(to, message string) {
		apiKey := os.Getenv("AFRICASTALKING_API_KEY")
		username := os.Getenv("AFRICASTALKING_USERNAME")

		data := url.Values{}
		data.Set("username", username)
		data.Set("to", to)
		data.Set("message", message)

		req, err := http.NewRequest("POST", "https://api.sandbox.africastalking.com/version1/messaging", strings.NewReader(data.Encode()))
		if err != nil {
			log.Printf("SMS error (request build): %v", err)
			return
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("apiKey", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("SMS error (send): %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			log.Printf("SMS failed with status: %s", resp.Status)
		} else {
			log.Println("SMS sent successfully to:", to)
		}
	}(to, message)
}
