package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/oauth2"
    "github.com/google/go-github/v38/github"
)

var (
    githubOauthConfig = oauth2.Config{
        ClientID:     Config("CLIENT_ID"),
        ClientSecret: Config("CLIENT_SECRET"),
        RedirectURL:  "http://localhost:3000/auth/github/callback",
        Scopes:       []string{"user"},
        Endpoint:     oauth2.Endpoint{
            AuthURL:  "https://github.com/login/oauth/authorize",
            TokenURL: "https://github.com/login/oauth/access_token",
        },
    }
)

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("¡Bienvenido a la aplicación!")
    })

    app.Get("/auth/github", func(c *fiber.Ctx) error {
        url := githubOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
        return c.Redirect(url)
    })

    app.Get("/auth/github/callback", func(c *fiber.Ctx) error {
        code := c.Query("code")

        token, err := githubOauthConfig.Exchange(c.Context(), code)
        if err != nil {
            log.Println(err)
            return c.SendStatus(fiber.StatusInternalServerError)
        }

        httpClient := githubOauthConfig.Client(c.Context(), token)
        client := github.NewClient(httpClient)

        user, _, err := client.Users.Get(c.Context(), "")
        if err != nil {
            log.Println(err)
            return c.SendStatus(fiber.StatusInternalServerError)
        }

        // Convertir los datos del usuario a un mapa
        userData := map[string]interface{}{
            "name":      user.Name,
            "avatar_url": user.AvatarURL,
            "id":        user.ID,
            "email":     user.Email,
        }
        // Devolver los datos en formato JSON
        return c.JSON(userData)
    })

    log.Fatal(app.Listen(":3000"))
}
