package neon

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
)

type NeonClaims struct {
	AuthID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	Exp    int64  `json:"exp"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type Client struct {
	jwksURL    string
	issuer     string
	jwtSecret  string
	mu         sync.RWMutex
	cachedKeys map[string]*rsa.PublicKey
	lastFetch  time.Time
}

func NewClient(cfg *configs.Config) *Client {
	c := &Client{
		jwksURL:    cfg.NeonAuthJWKSURL,
		issuer:     cfg.NeonAuthIssuer,
		jwtSecret:  cfg.JWTSecret,
		cachedKeys: make(map[string]*rsa.PublicKey),
	}

	if cfg.NeonAuthJWKSURL != "" {
		log.Printf("[Neon Auth] Initialized with JWKS URL: %s\n", cfg.NeonAuthJWKSURL)
		go c.refreshKeys()
	} else {
		log.Println("[Neon Auth] Initialized in local/hybrid JWT mode")
	}

	return c
}

func (c *Client) ValidateToken(tokenStr string) (*NeonClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	// 1. Check internal HMAC JWT first (for dev token & standalone auth)
	if claims, err := utils.ValidateJWT(c.jwtSecret, tokenStr); err == nil {
		authID := claims.SupabaseUserID
		if authID == "" && claims.UserID > 0 {
			authID = "usr_" + claims.Email
		}
		return &NeonClaims{
			AuthID: authID,
			Email:  claims.Email,
			Role:   claims.Role,
			Exp:    claims.Exp,
		}, nil
	}

	// 2. Parse JWT parts for Neon Authorize token
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT segment count")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid token header encoding")
	}

	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errors.New("invalid token header json")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload encoding")
	}

	var claims NeonClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token payload json")
	}

	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, errors.New("token has expired")
	}

	// If JWKS URL is configured, verify RS256 signature
	if c.jwksURL != "" && header.Alg == "RS256" {
		pubKey, err := c.getKey(header.Kid)
		if err != nil {
			return nil, err
		}

		sig, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			return nil, errors.New("invalid signature encoding")
		}

		signedData := []byte(parts[0] + "." + parts[1])
		hasher := crypto.SHA256.New()
		hasher.Write(signedData)
		digest := hasher.Sum(nil)

		if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest, sig); err != nil {
			return nil, errors.New("invalid signature")
		}
	}

	if claims.AuthID == "" {
		return nil, errors.New("token missing subject claim (sub)")
	}

	return &claims, nil
}

func (c *Client) getKey(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	key, ok := c.cachedKeys[kid]
	c.mu.RUnlock()

	if ok {
		return key, nil
	}

	c.refreshKeys()

	c.mu.RLock()
	defer c.mu.RUnlock()
	key, ok = c.cachedKeys[kid]
	if !ok {
		return nil, errors.New("unknown key id (kid) in token")
	}
	return key, nil
}

func (c *Client) refreshKeys() {
	if c.jwksURL == "" {
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(c.jwksURL)
	if err != nil {
		log.Printf("[Neon Auth] Failed to fetch JWKS: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, k := range jwks.Keys {
		if k.Kty != "RSA" {
			continue
		}
		nBytes, err1 := base64.RawURLEncoding.DecodeString(k.N)
		eBytes, err2 := base64.RawURLEncoding.DecodeString(k.E)
		if err1 != nil || err2 != nil {
			continue
		}

		var eInt int
		for _, b := range eBytes {
			eInt = (eInt << 8) | int(b)
		}

		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: eInt,
		}
		newKeys[k.Kid] = pubKey
	}

	c.mu.Lock()
	c.cachedKeys = newKeys
	c.lastFetch = time.Now()
	c.mu.Unlock()
}

// ParseRSAPublicKeyFromPEM is a utility for testing with custom PEM public keys
func ParseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}
