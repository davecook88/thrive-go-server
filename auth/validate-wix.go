package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var USER_INSTANCE_CONTEXT_KEY = "wix_user_instance"

type Instance struct {
	InstanceId     string  `json:"instanceId"`
	AppDefId       string  `json:"appDefId"`
	SignDate       string  `json:"signDate"`
	Uid            string  `json:"uid"`
	Permissions    string  `json:"permissions"`
	DemoMode       bool    `json:"demoMode"`
	SiteOwnerId    string  `json:"siteOwnerId"`
	SiteMemberId   string  `json:"siteMemberId"`
	ExpirationDate string  `json:"expirationDate"`
	LoginAccountId string  `json:"loginAccountId"`
	Pai            *string `json:"pai"`
	Lpai           *string `json:"lpai"`
}

func decodeHash(hash string) string {
	hash = strings.ReplaceAll(hash, "-", "+")
	hash = strings.ReplaceAll(hash, "_", "/")
	return hash
}

func GetUserInstance(c *gin.Context) *Instance {
	instance, exists := c.Get(USER_INSTANCE_CONTEXT_KEY)
	if !exists {
		return nil
	}
	return instance.(*Instance)
}

func validateInstance(hash string, payload string, secret string) bool {
	if hash == "" {
		return false
	}
	hash = decodeHash(hash)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedHash := base64.RawStdEncoding.EncodeToString(mac.Sum(nil))
	return hash == expectedHash
}

func parseInstance(instance string, appSecret string) (*Instance, error) {
	parts := strings.Split(instance, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid instance format")
	}
	hash := parts[0]
	payload := parts[1]

	if payload == "" {
		return nil, errors.New("no payload")
	}
	if !validateInstance(hash, payload, appSecret) {
		return nil, errors.New("invalid instance")
	}

	decodedPayload, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		println("Decoding failed:", err)
		return nil, err
	}
	var userInstance Instance
	err = json.Unmarshal(decodedPayload, &userInstance)
	if err != nil {
		println("JSON decoding failed:", err)
		return nil, err
	}
	return &userInstance, nil
}

func ValidateWixHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		println("Validating Wix header")
		instance := c.GetHeader("Authorization")
		println("Instance: ", instance)
		if instance == "" {
			c.JSON(400, gin.H{"error": "No instance header"})
			c.Abort()
			return
		}
		appSecret := os.Getenv("WIX_SECRET_KEY")
		if appSecret == "" {
			c.JSON(500, gin.H{"error": "No app secret"})
			c.Abort()
			return
		}

		user_instance, err := parseInstance(instance, appSecret)
		c.Set(USER_INSTANCE_CONTEXT_KEY, user_instance)

		if err != nil {
			println("Error: ", err)
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
