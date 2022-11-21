package main

import (
	"api-gateway/plugins/login/gateway"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func init() {
	fmt.Println("login plugin loaded!!!")
}

var ClientRegisterer = registerer("login")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), func(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
		cfg := parse(extra)
		if cfg == nil {
			return nil, errors.New("wrong config")
		}
		if cfg.name != string(r) {
			return nil, fmt.Errorf("unknown register %s", cfg.name)
		}
		sessionSecret := DefaultValue("SESSION_SECRET", "secret")
		return gateway.New(ctx, cfg.name, gateway.Config{
			CookieName:     cfg.cookie,
			JwtSecretKey:   os.Getenv("JWT_SECRET"),
			JwtIss:         os.Getenv("JWT_ISS"),
			CookieSecure:   ConvertBool("COOKIE_SECURE"),
			CookieHttpOnly: ConvertBool("COOKIE_HTTP_ONLY"),
			CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
			SessionSecret:  sessionSecret,
		})
	})
}

func parse(extra map[string]interface{}) *opts {
	name, ok := extra["name"].(string)
	if !ok {
		return nil
	}
	cookieName, ok := extra["cookie_name"].(string)
	if !ok {
		return nil
	}
	return &opts{
		name:   name,
		cookie: cookieName,
	}
}

func ConvertBool(env string) bool {
	v, _ := strconv.ParseBool(os.Getenv(env))
	return v
}

func DefaultValue(env, defaultValue string) string {
	val := os.Getenv(env)
	if val == "" {
		return defaultValue
	}

	return val
}

type opts struct {
	name   string
	cookie string
}

func main() {}
