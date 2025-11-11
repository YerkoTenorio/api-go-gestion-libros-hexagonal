package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TursoURL   string
	TursoToken string
	Port       int
}

// Load lee .env (si existe) y variables del entorno
func Load() (Config, error) {
	_ = godotenv.Load() // no falla si .env no existe

	c := Config{
		TursoURL:   os.Getenv("TURSO_DATABASE_URL"),
		TursoToken: os.Getenv("TURSO_AUTH_TOKEN"),
		Port:       8080,
	}

	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			c.Port = v
		}
	}

	if c.TursoURL == "" {
		return Config{}, fmt.Errorf("missing TURSO_DATABASE_URL")
	}
	if c.TursoToken == "" {
		return Config{}, fmt.Errorf("missing TURSO_AUTH_TOKEN")
	}

	return c, nil
}
