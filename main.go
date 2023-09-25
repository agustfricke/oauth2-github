package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/agustfricke/go-oauth-example/config"
	"github.com/gofiber/fiber/v2"
)


func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/oauth/redirect", func(c *fiber.Ctx) error {

		code := c.Query("code")
    clientID := config.Config("CLIENT_ID")
    clientSecret := config.Config("CLIENT_SECRET")

		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)

		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "no se pudo crear la solicitud HTTP: %v", err)
			c.Status(http.StatusBadRequest)
			return nil
		}

		req.Header.Set("accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "no se pudo enviar la solicitud HTTP: %v", err)
			c.Status(http.StatusInternalServerError)
			return nil
		}
		defer res.Body.Close()

		var t OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "no se pudo analizar la respuesta JSON: %v", err)
			c.Status(http.StatusBadRequest)
			return nil
		}

		return c.Redirect("/welcome.html?access_token=" + t.AccessToken, http.StatusFound)
	})
	app.Listen(":8080")
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
