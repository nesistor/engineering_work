package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ScopeAuthentication = "authentication"
	ScopeRefresh        = "refresh"
	RoleAdmin           = "admin"
	RoleUser            = "user"
)

// TokenModel holds the Redis client and KeyManager for token handling.
type TokenModel struct {
	RedisClient *redis.Client
	KeyManager  *KeyManager
}

// Models represents all models in the application.
type Models struct {
	Token TokenModel
}

// New creates a new instance of Models with initialized TokenModel.
func New(redisClient *redis.Client, keyManager *KeyManager) Models {
	return Models{
		Token: TokenModel{
			RedisClient: redisClient,
			KeyManager:  keyManager,
		},
	}
}

// JWTClaims stores extra JWT information, including role and scope.
type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`  // Role (admin/user)
	Scope  string `json:"scope"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT for the specified user with role, scope, and TTL.
func (m *TokenModel) GenerateToken(ctx context.Context, userID int, role string, ttl time.Duration, scope, kid string) (string, error) {
	claims := JWTClaims{
		UserID: int64(userID),
		Role:   role,
		Scope:  scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	privateKey := m.KeyManager.GetPrivateKey()
	if privateKey == nil {
		return "", fmt.Errorf("failed to retrieve private key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// InsertDeactivatedToken stores a deactivated token in Redis.
func (m *TokenModel) InsertDeactivatedToken(ctx context.Context, tokenString string, ttl time.Duration) error {
	key := fmt.Sprintf("deactivated_token:%s", tokenString)

	err := m.RedisClient.Set(ctx, key, "deactivated", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to insert deactivated token: %v", err)
	}

	return nil
}

// IsTokenDeactivated checks if a token has been added to the deactivated list in Redis.
func (m *TokenModel) IsTokenDeactivated(ctx context.Context, tokenString string) (bool, error) {
	key := fmt.Sprintf("deactivated_token:%s", tokenString)

	exists, err := m.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if token is deactivated: %v", err)
	}

	return exists > 0, nil
}

// GetUserIDForToken retrieves the user ID and role from a token, ensuring it is valid and has the correct scope.
func (m *TokenModel) GetUserIDForToken(ctx context.Context, tokenString, scope string) (int64, string, error) {

	deactivated, err := m.IsTokenDeactivated(ctx, tokenString)
	if err != nil {
		return 0, "", err
	}
	if deactivated {
		return 0, "", fmt.Errorf("token has been deactivated")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse token: %v", err)
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return 0, "", fmt.Errorf("kid missing from token header")
	}

	publicKey, err := m.KeyManager.GetPublicKey(kid)
	if err != nil {
		return 0, "", fmt.Errorf("public key not found for kid: %s, %v", kid, err)
	}

	token, err = jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("failed to verify token: %v", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid || claims.Scope != scope {
		return 0, "", fmt.Errorf("invalid or unauthorized token")
	}

	return claims.UserID, claims.Role, nil
}

// RefreshAccessToken creates a new access token based on a valid refresh token.
func (m *TokenModel) RefreshAccessToken(ctx context.Context, refreshToken, kid string) (string, error) {

	userID, role, err := m.GetUserIDForToken(ctx, refreshToken, ScopeRefresh)
	if err != nil {
		return "", fmt.Errorf("failed to refresh access token: %v", err)
	}

	accessToken, err := m.GenerateToken(ctx, int(userID), role, 15*time.Minute, ScopeAuthentication, kid)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return accessToken, nil
}
