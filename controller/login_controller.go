package controller

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kirtanwyn/allPanelexch/config"
	"github.com/kirtanwyn/allPanelexch/model"
	"github.com/kirtanwyn/allPanelexch/service"
)

type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func LoginCheck(c *gin.Context) {
	db := config.DB

	// Check Maintenance Mode
	var siteStatus int
	err := db.QueryRow("SELECT site_status FROM site_under_maintenance LIMIT 1").Scan(&siteStatus)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database error"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		req.Username = c.PostForm("username")
		req.Password = c.PostForm("password")
	}

	// if req.Username != "democ2" {
	// 	c.JSON(http.StatusOK, gin.H{"status": "maintenance", "msg": "We are planning to have scheduled maintenance."})
	// 	return
	// }

	if siteStatus == 1 {
		c.JSON(http.StatusOK, gin.H{"status": "maintenance", "msg": "We are planning to have scheduled maintenance."})
		return
	}

	loginIPAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Find user
	var user model.UserLoginMaster
	err = db.QueryRow("SELECT UserID, Password, user_password_salt, user_password_salt_key, SecretKey, first_password_changed, parentDL, parentMDL, parentSuperMDL, parentKingAdmin FROM user_login_master WHERE Email_ID=? AND UserType IN (1)", req.Username).Scan(
		&user.UserID, &user.Password, &user.UserPasswordSalt, &user.UserPasswordSaltKey, &user.SecretKey, &user.FirstPasswordChanged, &user.ParentDL, &user.ParentMDL, &user.ParentSuperMDL, &user.ParentKingAdmin,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"status": "0", "message": "Invalid username"})
		return
	} else if err != nil {
		log.Println("here am i ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database error"})
		return
	}
	log.Println("here am i 2")

	// Handle empty salt (Migration path)
	if (user.UserPasswordSalt == nil || *user.UserPasswordSalt == "") && (user.UserPasswordSaltKey == nil || *user.UserPasswordSaltKey == "") {
		pSalt := strconv.Itoa(111111 + time.Now().Nanosecond()%888889) // Random 6 digit
		siteSalt := "huhefcvringybh"
		saltedHash := service.HashPassword(req.Password, siteSalt, pSalt)

		_, err := db.Exec("UPDATE user_login_master SET user_password_salt=?, user_password_salt_key=?, Password=? WHERE UserID=?", pSalt, siteSalt, saltedHash, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "0"})
		return
	}

	var userSaltKey, userSalt string
	if user.UserPasswordSaltKey != nil {
		userSaltKey = *user.UserPasswordSaltKey
	}
	if user.UserPasswordSalt != nil {
		userSalt = *user.UserPasswordSalt
	}

	// Validate Password
	saltedHash := service.HashPassword(req.Password, userSaltKey, userSalt)
	if user.Password != saltedHash {
		c.JSON(http.StatusOK, gin.H{"status": "0"})
		return
	}

	// Fetch User Master Details
	var userMaster model.UserMaster
	err = db.QueryRow("SELECT Name, Status, user_verification_status, user_verification_type FROM user_master WHERE Id=?", user.UserID).Scan(
		&userMaster.Name, &userMaster.Status, &userMaster.UserVerificationStatus, &userMaster.UserVerificationType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database error"})
		return
	}

	// Validate Parent Hierarchy
	checkParentStatus := func(id int) (string, error) {
		var status string
		err := db.QueryRow("SELECT Status FROM user_master WHERE Id=?", id).Scan(&status)
		return status, err
	}

	dlStatus, _ := checkParentStatus(user.ParentDL)
	mdlStatus, _ := checkParentStatus(user.ParentMDL)
	smdlStatus, _ := checkParentStatus(user.ParentSuperMDL)
	kingadminStatus, _ := checkParentStatus(user.ParentKingAdmin)

	if userMaster.Status == "1" && dlStatus == "1" && mdlStatus == "1" && smdlStatus == "1" && kingadminStatus == "1" {
		if user.SecretKey == nil || *user.SecretKey == "" {
			// 2FA Check
			authCheck := true // Replace with environment variable or config if needed
			if userMaster.UserVerificationStatus == "ENABLED" && authCheck {
				if userMaster.UserVerificationType == "Telegram" {
					if err := service.GenerateTelegramOTP(user.UserID); err != nil {
						c.JSON(http.StatusOK, gin.H{"status": "telegram_error", "msg": err.Error()})
						return
					}
				}

				// Set Session state for Auth Code Verification
				session := sessions.Default(c)
				// session.Options(sessions.Options{
				// 	Path:     "/",
				// 	MaxAge:   60 * 60 * 24, // 24 hours
				// 	HttpOnly: true,
				// 	Secure:   false, // true in production with HTTPS
				// })

				session.Set("CLIENT_AUTH_STATUS", true)
				session.Set("CLIENT_AUTH_UID", user.UserID)
				// session.Save()
				if err := session.Save(); err != nil {
					log.Println("failed to save session:", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  "error",
						"message": "Failed to create session",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"status":       "auth",
					"user_auth_id": user.UserID,
				})
				return
			}

			dateTime := time.Now().Format("2006-01-02 15:04:05")

			// Logging
			_, err = db.Exec("INSERT INTO login_ip_address(user_id, ip_address, login_date_time, user_agent) VALUES(?, ?, ?, ?)", user.UserID, loginIPAddress, dateTime, userAgent)
			if err != nil {
				log.Println("Failed to log login ip address:", err)
			}
			_, err = db.Exec("INSERT INTO activity_log(user_id, username, ip_address, user_agent, date_time, log_type) VALUES(?, ?, ?, ?, ?, ?)", user.UserID, req.Username, loginIPAddress, userAgent, dateTime, "login")
			if err != nil {
				log.Println("Failed to log activity:", err)
			}

			loginRandomString := service.GenerateRandomString(18) // Simplified implementation

			// Token logic
			firstTwoChars := ""
			if len(userMaster.Name) >= 2 {
				firstTwoChars = userMaster.Name[:2]
			}
			//apiAuthToken := fmt.Sprintf("%s%s", firstTwoChars, service.GenerateRandomString(18))
			apiAuthToken := strings.ToUpper(
				fmt.Sprintf("%s%s", firstTwoChars, service.GenerateRandomString(18)),
			)

			_, err = db.Exec("UPDATE user_login_master SET loginString=?, api_auth_token=? WHERE Id=?", loginRandomString, apiAuthToken, user.UserID)
			if err != nil {
				log.Println("Failed to update user login master:", err)
			}

			// Save to Redis Session
			session := sessions.Default(c)
			session.Set("CLIENT_LOGIN_STATUS", true)
			session.Set("CLIENT_LOGIN_ID", user.UserID)
			session.Set("CLIENT_LOGIN_NAME", userMaster.Name)
			session.Set("FIRST_PASSWORD_CHANGED", user.FirstPasswordChanged)
			session.Set("LOGIN_ENC_ID", base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(user.UserID))))
			session.Set("LOGIN_STRING", loginRandomString)
			// session.Save()
			if err := session.Save(); err != nil {
				log.Println("failed to save session:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Failed to create session",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":                 user.UserID,
				"login_id":               base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(user.UserID))),
				"login_string":           loginRandomString,
				"first_password_changed": user.FirstPasswordChanged,
			})
		} else {
			session := sessions.Default(c)
			session.Set("temp_id", user.UserID)
			// session.Save()
			if err := session.Save(); err != nil {
				log.Println("failed to save session:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Failed to create session",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "skey"})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "NA"})
	}
}

// Logout clears the user's authentication cookies and logs them out
func Logout(c *gin.Context) {
	// Clear the session from Redis
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})
	// session.Save()
	if err := session.Save(); err != nil {
		log.Println("failed to save session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create session",
		})
		return
	}

	// Since this might be called as an API or a direct link, we can handle both.
	// If it's an AJAX call, return JSON. If it's a direct navigation, we can redirect.
	// if c.GetHeader("Accept") == "application/json" || c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
	// 	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Logged out successfully"})
	// } else {
	// 	// Redirect to the login page (legacy behavior)
	// 	c.Redirect(http.StatusFound, "/login")
	// }
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Logged out successfully"})

}


func GetSession(c *gin.Context) {
	session := sessions.Default(c)

	c.JSON(200, gin.H{
		"CLIENT_LOGIN_STATUS":     session.Get("CLIENT_LOGIN_STATUS"),
		"CLIENT_LOGIN_ID":         session.Get("CLIENT_LOGIN_ID"),
		"CLIENT_LOGIN_NAME":       session.Get("CLIENT_LOGIN_NAME"),
		"FIRST_PASSWORD_CHANGED":  session.Get("FIRST_PASSWORD_CHANGED"),
		"LOGIN_ENC_ID":            session.Get("LOGIN_ENC_ID"),
		"LOGIN_STRING":            session.Get("LOGIN_STRING"),
		"CLIENT_AUTH_STATUS":      session.Get("CLIENT_AUTH_STATUS"),
		"CLIENT_AUTH_UID":         session.Get("CLIENT_AUTH_UID"),
		"TEMP_ID":                 session.Get("temp_id"),
	})
}