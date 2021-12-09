package middlewares

import (
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	j "github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Access returns a middleware that records an access log message for every HTTP request being processed.
func Access(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := c.Request.Context()
		ctx = log.WithRequest(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)

		// Start logging request access log.
		logger.With(ctx, "http", "request", "start", start).
			Infof("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)

		c.Next()

		// End logging response access log.
		logger.With(ctx, "http", "response", "duration", time.Since(start).Milliseconds(), "status", c.Writer.Status()).
			Infof("%s %s %s %d %d", c.Request.Method, c.Request.URL.Path, c.Request.Proto, c.Writer.Status(), c.Writer.Size())
	}
}

// Auth check if auth ok and set claims in request header.
func Auth(logger log.Logger, secret string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		sessionAuth := sessions.DefaultMany(c, "session_auth")
		if sessionAuth == nil {
			logg.Info("session_auth not found")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrNoTokenFound)
			return
		}

		tokenInterface := sessionAuth.Get("access_token")
		if tokenInterface == nil {
			logg.Info("access token not found in session cookies")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrNoTokenFound)
			return
		}
		accessToken := tokenInterface.(string)

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logg.Warnf("fail to parse access token")
				_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenInvalid)
				return token, e.ErrParseToken
			}
			return []byte(secret), nil
		})

		if err != nil {
			logg.Infof("error to parse access token: %s", err.Error())
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenExpired)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("identity_id", claims["identity_id"].(string))
			c.Set("identity_provider", claims["identity_provider"].(string))
			c.Set("identity_uid", claims["identity_uid"].(string))
			c.Set("user_id", claims["user_id"].(string))
			c.Set("user_role", claims["user_role"].(string))

			// TODO: try to update access and refresh tokens.
			accessTokenExpires := int64(claims["exp"].(float64))
			refreshTokenID := sessionAuth.Get("refresh_token_id").(string)
			refreshTokenExpires := sessionAuth.Get("refresh_token_exp").(int64)

			// Logic for update access-token.
			tm := time.Unix(accessTokenExpires, 0)
			remains := -(time.Since(tm).Minutes())

			// If there is less 6 minutes left, create a new access token and put in Cookie.
			var jwtRefreshToken, jwtAccessToken entity.Token
			var refreshTokenValid, updateCookies bool
			if remains < 6 {
				err = db.Debug().Model(&entity.Token{}).Where("id = ?", refreshTokenID).Take(&jwtRefreshToken).Error
				if err == nil {
					_, refreshTokenValid := j.ValidToken(secret, jwtRefreshToken.RefreshToken)
					if refreshTokenValid {
						jwtAccessToken, err = j.UpdateAccessToken(secret, claims)
						sessionAuth.Set("access_token", jwtAccessToken.AccessToken)
						updateCookies = true
					}
				}
			}

			// Logic for update refresh-token.
			tm = time.Unix(refreshTokenExpires, 0)
			remains = -(time.Since(tm).Hours())
			// If there is less 2 hours left, create a new access token.
			if remains < 2 && refreshTokenValid {
				// Create a new refresh token and put in Cookie.
				jwtRefreshToken, err = j.UpdateRefreshToken(secret, claims)

				// Create first and delete after.
				if err = db.Debug().Model(&entity.Token{}).Create(&jwtRefreshToken).Error; err == nil {
					if err = db.Debug().Model(&entity.Token{}).Where("id = ?", refreshTokenID).Delete(&entity.Token{ID: refreshTokenID}).Error; err == nil {
						sessionAuth.Set("refresh_token_id", jwtRefreshToken.RefreshToken)
						sessionAuth.Set("refresh_token_exp", jwtRefreshToken.RefreshExpires)
						updateCookies = true
					}
				}
			}

			if updateCookies {
				_ = sessionAuth.Save()
			}
			c.Next()
		} else {
			logg.Warnf("invalid token! so sorry")
			_ = c.AbortWithError(http.StatusUnauthorized, e.ErrTokenInvalid)
		}
	}
}

// Checks returns a middleware that verify some points before business logic.
func Checks(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		if c.Request.Method == "POST" || c.Request.Method == "PATCH" {
			if c.Request.Body == http.NoBody {
				logg.Warnf("Empty body in POST or PATCH request")
				_ = c.AbortWithError(http.StatusBadRequest, e.ErrRequestNeedBody)
			}
		}
		c.Next()
	}
}

// Authorizer check if user role has access to resource.
func Authorizer(en *casbin.Enforcer, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logg := logger.With(c.Request.Context())

		var role string
		roleInterface, ok := c.Get("user_role")
		if ok {
			role = roleInterface.(string)
		}

		if role == "" {
			role = "anonymous"
		}

		// casbin rule enforcing
		res, err := en.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			logg.Error("error to enforce casbin authorization: ", err.Error())
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if res {
			c.Next()
		} else {
			_ = c.AbortWithError(http.StatusForbidden, e.ErrUnauthorized)
			return
		}
	}
}
