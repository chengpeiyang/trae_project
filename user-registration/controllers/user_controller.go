package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"user-registration/config"
	"user-registration/models"
	"user-registration/utils"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Source   string `json:"source"`
}

type RegisterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	} `json:"data"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名已存在",
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: utils.MD5Encode(req.Password),
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	go createRegisterLog(c, &user, &req)

	if config.Redis != nil {
		cacheUser(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func validateRegisterRequest(req *RegisterRequest) error {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return &ValidationError{Message: "用户名不能为空"}
	}

	if len(username) < 3 || len(username) > 20 {
		return &ValidationError{Message: "用户名长度必须在3-20个字符之间"}
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return &ValidationError{Message: "用户名只能包含字母、数字和下划线"}
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		return &ValidationError{Message: "密码不能为空"}
	}

	if len(password) < 6 || len(password) > 32 {
		return &ValidationError{Message: "密码长度必须在6-32个字符之间"}
	}

	if req.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			return &ValidationError{Message: "邮箱格式不正确"}
		}
	}

	if req.Phone != "" {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(req.Phone) {
			return &ValidationError{Message: "手机号格式不正确"}
		}
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func createRegisterLog(c *gin.Context, user *models.User, req *RegisterRequest) {
	userAgent := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")
	source := req.Source
	if source == "" {
		source = "web"
	}

	deviceType, os, browser := parseUserAgent(userAgent)

	log := models.RegisterLog{
		UserID:     user.ID,
		Username:   user.Username,
		IP:         utils.GetClientIP(c),
		UserAgent:  userAgent,
		Source:     source,
		DeviceType: deviceType,
		OS:         os,
		Browser:    browser,
		Referrer:   referrer,
		RequestURI: c.Request.RequestURI,
		CreatedAt:  time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		config.DB.Logger.Error(c, "创建注册日志失败: %v", err)
	}
}

func parseUserAgent(userAgent string) (deviceType, os, browser string) {
	deviceType = "Unknown"
	os = "Unknown"
	browser = "Unknown"

	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		deviceType = "Mobile"
	} else {
		deviceType = "Desktop"
	}

	if strings.Contains(ua, "windows") {
		os = "Windows"
	} else if strings.Contains(ua, "mac os") {
		os = "MacOS"
	} else if strings.Contains(ua, "linux") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "iOS"
	}

	if strings.Contains(ua, "chrome") {
		browser = "Chrome"
	} else if strings.Contains(ua, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(ua, "safari") {
		browser = "Safari"
	} else if strings.Contains(ua, "edge") {
		browser = "Edge"
	} else if strings.Contains(ua, "opera") {
		browser = "Opera"
	} else if strings.Contains(ua, "msie") || strings.Contains(ua, "trident") {
		browser = "IE"
	}

	return deviceType, os, browser
}

func cacheUser(user *models.User) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return
	}

	key := "user:" + user.Username
	config.Redis.Set(config.Ctx, key, userJSON, 24*time.Hour)
}
