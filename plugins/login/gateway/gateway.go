package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginHandler struct {
	ctx            context.Context
	Client         *http.Client
	CookieName     string
	CookieExpired  time.Duration
	Name           string
	JwtSecretKey   []byte
	JwtIss         string
	CookieSecure   bool
	CookieHttpOnly bool
	CookieDomain   string
	SessionSecret  string
}

type ResponseKrakend struct {
	Data DataKrakend `json:"data"`
}

type DataKrakend struct {
	Id        string   `json:"id"`
	Email     string   `json:"email"`
	Msisdn    string   `json:"msisdn"`
	Roles     []string `json:"roles"`
	LastLogin string   `json:"lastLogin"`
}

type ResponseKrakendError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type ResponseAuth struct {
	Data PayloadAuth `json:"data"`
}

type PayloadAuth struct {
	Id        string   `json:"id"`
	Email     string   `json:"email"`
	Msisdn    string   `json:"msisdn"`
	Roles     []string `json:"roles"`
	LastLogin string   `json:"lastLogin"`
}

type ResponseError struct {
	Code    PayloadAuth   `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}

func (h *LoginHandler) setCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     h.CookieName,
		Value:    token,
		Path:     "/",
		Secure:   h.CookieSecure,
		HttpOnly: h.CookieHttpOnly,
		SameSite: http.SameSiteNoneMode,
		Domain:   h.CookieDomain,
		Expires:  time.Now().Add(h.CookieExpired),
	}
	http.SetCookie(w, &cookie)
}

func (h *LoginHandler) generateToken(payload *PayloadAuth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": payload.Email,
		"sub":   payload.Id,
		"iss":   os.Getenv("JWT_ISS"),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(h.JwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data ResponseAuth

	req, err := http.NewRequestWithContext(h.ctx, r.Method, fmt.Sprintf("%s://%s%s", r.URL.Scheme, r.URL.Host, r.URL.Path), r.Body)

	log.Println("req after ", req.Body)
	if err != nil {
		log.Println("failed to call backend:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Host = r.Host
	req.URL.Opaque = r.URL.RequestURI()
	for key, value := range r.Header {
		req.Header[key] = value
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		log.Println("failed to write response:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println("failed defer:", err.Error())
					return
				}
			}(resp.Body)

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Println("failed decode:", err.Error())
				return
			}

			for key, value := range resp.Header {
				w.Header()[key] = value
			}

			token, err := h.generateToken(&data.Data)
			if err != nil {
				log.Println(err)
				return
			}
			h.setCookie(w, token)

			w.Header().Del("Content-Length")
			w.WriteHeader(resp.StatusCode)

			responseObj := ResponseKrakend{
				Data: DataKrakend{
					Id:        data.Data.Id,
					Email:     data.Data.Email,
					Msisdn:    data.Data.Msisdn,
					Roles:     data.Data.Roles,
					LastLogin: data.Data.LastLogin,
				},
			}

			err = json.NewEncoder(w).Encode(&responseObj)
			if err != nil {
				log.Println("failed to write response:", err.Error())
			}
		} else {
			for key, value := range resp.Header {
				w.Header()[key] = value
			}
			w.Header().Del("Content-Length")
			w.WriteHeader(resp.StatusCode)
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				log.Println("failed decode:", err.Error())
				return
			}
		}
	}
}

type Config struct {
	CookieName     string
	JwtSecretKey   string
	JwtIss         string
	CookieSecure   bool
	CookieHttpOnly bool
	CookieDomain   string
	SessionSecret  string
}

func New(ctx context.Context, name string, cfg Config) (http.Handler, error) {
	lgn := LoginHandler{
		ctx: ctx,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		CookieName:     cfg.CookieName,
		CookieExpired:  time.Hour * 24,
		Name:           name,
		JwtSecretKey:   []byte(cfg.JwtSecretKey),
		JwtIss:         cfg.JwtIss,
		CookieSecure:   cfg.CookieSecure,
		CookieHttpOnly: cfg.CookieHttpOnly,
		CookieDomain:   cfg.CookieDomain,
		SessionSecret:  cfg.SessionSecret,
	}
	mux := http.NewServeMux()
	mux.Handle("/", &lgn)
	return mux, nil
}
