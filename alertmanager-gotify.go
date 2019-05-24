package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gotify/go-api-client/auth"
	"github.com/gotify/go-api-client/client/message"
	"github.com/gotify/go-api-client/gotify"
	"github.com/gotify/go-api-client/models"
)

type POSTData struct {
	Alerts []Alert `json:alerts`
}

type Alert struct {
	Annotations Annotations `json:annotations`
	Labels Labels `json:labels`
	Status string `json:status`
}

type Annotations struct {
	Description string `json:description`
	Title string `json:title`
}

type Labels struct {
	Severity string `json:severity`
}

func main() {
	gotifyURL := os.Getenv("GOTIFY_URL")

	http.HandleFunc("/alert", func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		applicationToken := r.URL.Query().Get("token")
		if applicationToken == "" {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body %v", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		postData := POSTData{}
		err = json.Unmarshal(body, &postData)
		if err != nil {
			log.Printf("Failed to parse json %v", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		parsedURL, _ := url.Parse(gotifyURL)
		client := gotify.NewClient(parsedURL, &http.Client{})

		for _, alert := range postData.Alerts {
			var title strings.Builder

			if alert.Status == "resolved" {
				title.WriteString("RESOLVED: ")
			} else {
				title.WriteString(strings.ToUpper(alert.Labels.Severity))
				title.WriteString(": ")
			}

			title.WriteString(alert.Annotations.Title)

			params := message.NewCreateMessageParams()
			params.Body = &models.MessageExternal{
				Title:    title.String(),
				Message:  alert.Annotations.Description,
			}
			_, err = client.Message.CreateMessage(params, auth.TokenAuth(applicationToken))

			if err != nil {
				log.Printf("Could not send message %v", err)
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	})

	http.ListenAndServe(":8081", nil)
}
