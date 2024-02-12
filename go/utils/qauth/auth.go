package qauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mhaqqiw/sdk/go/qconstant"
	"github.com/mhaqqiw/sdk/go/qentity"
	h "github.com/mhaqqiw/sdk/go/utils/qhttp"
	"github.com/mhaqqiw/sdk/go/utils/qlog"
	"github.com/mhaqqiw/sdk/go/utils/qmodule"
	"github.com/mhaqqiw/sdk/go/utils/qredis"

	"github.com/gin-gonic/gin"
)

// Non Login User, Doesn't Need Captcha
func WS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start", time.Now())
		c.Next()
	}
}

// Non Login User, Doesn't Need Captcha
func Type1(recaptcha qentity.Recaptcha) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start", time.Now())
		err := validateRecaptcha(c, recaptcha)
		if err != nil {
			return
		}
		c.Next()
	}
}

// Need Login User, Need Captcha
func Type2(recaptcha qentity.Recaptcha) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start", time.Now())
		err := validateRecaptcha(c, recaptcha)
		if err != nil {
			return
		}
		session, err := validateSession(c)
		if err != nil {
			return
		}
		user, err := getUserData(session)
		if err != nil {
			h.Return(c, http.StatusInternalServerError, err)
			return
		}
		if user.Name == "" {
			h.Return(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// Need Login User, Can Use Token
func Type3(initToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start", time.Now())
		session, err := validateSession(c)
		if err != nil {
			return
		}
		token := c.Request.FormValue("api_key")
		//TODO: Validate Token
		if token != initToken {
			user, err := getUserData(session)
			if err != nil {
				h.Return(c, http.StatusInternalServerError, err)
				return
			}
			if user.Name == "" {
				h.Return(c, http.StatusUnauthorized, "unauthorized")
				return
			}
			c.Set("user", user)
		}
		c.Next()
	}
}

// Non Login User, Need Captcha
func Type4(recaptcha qentity.Recaptcha) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start", time.Now())
		err := validateRecaptcha(c, recaptcha)
		if err != nil {
			return
		}
		_, err = validateSession(c)
		if err != nil {
			return
		}
		c.Next()
	}
}

func getUserData(session string) (qentity.SessionData, error) {
	var user qentity.SessionData
	userData, _, err := qredis.Get("session", session)
	if err != nil {
		qlog.LogPrint(qconstant.ERROR, "qredis.Get", qlog.Trace(), err.Error())
		return user, err
	}
	if userData == "" {
		return user, nil
	}
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		qlog.LogPrint(qconstant.ERROR, "json.Unmarshal", qlog.Trace(), err.Error())
		return user, err
	}
	user.SessionID = session
	return user, nil
}

func validateRecaptcha(c *gin.Context, recaptcha qentity.Recaptcha) error {
	captchaToken := c.GetHeader(qconstant.CAPTCHA_TOKEN)
	if captchaToken == "" {
		h.Return(c, http.StatusBadRequest, "bad request")
		return errors.New("missing captcha token")
	}
	body, err := qmodule.CheckRecaptcha(recaptcha.Secret, captchaToken, recaptcha.ValidateURL)
	if err != nil {
		h.Return(c, http.StatusInternalServerError, err)
		return err
	}
	if !body.Success {
		h.Return(c, http.StatusBadRequest, "bad request")
		return errors.New("failed to validate recaptcha")
	}
	if body.Score < recaptcha.Threshold {
		h.Return(c, http.StatusForbidden, "you are detected as a bot")
		return errors.New("you are detected as a bot")
	}
	if body.Action != recaptcha.Action {
		h.Return(c, http.StatusBadRequest, "mismatch recaptcha action")
		return errors.New("mismatch recaptcha action")
	}
	return nil
}

func validateSession(c *gin.Context) (string, error) {
	session, err := c.Cookie(qconstant.SESSION)
	if err != nil {
		session, err = qmodule.GenerateUUIDV1()
		if err != nil {
			h.Return(c, http.StatusInternalServerError, err)
			return session, errors.New("invalid generate session")
		}
		c.SetCookie(qconstant.SESSION, session, 0, "/", "*", false, true)
	}
	if session != "" {
		c.Set("session", session)
	}
	return session, nil
}
