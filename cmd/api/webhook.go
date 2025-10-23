package main

import (
	"io"
	"net/http"
	"strings"
)

func (app *application) triggerDeployWebhook() {
	url := app.config.deployWebhook
	if url == "" {
		return
	}

	go func(target string) {
		req, err := http.NewRequest(http.MethodPost, target, strings.NewReader("{}"))
		if err != nil {
			app.logger.Printf("Error creating deploy webhook request: %s", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := app.httpClient
		if client == nil {
			client = http.DefaultClient
		}

		resp, err := client.Do(req)
		if err != nil {
			app.logger.Printf("Error invoking deploy webhook: %s", err)
			return
		}
		defer resp.Body.Close()

		io.Copy(io.Discard, resp.Body)

		if resp.StatusCode >= 300 {
			app.logger.Printf("Deploy webhook returned status %d", resp.StatusCode)
		}
	}(url)
}
