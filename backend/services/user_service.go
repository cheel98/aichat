package services

import (
	"aiChat/backend/database"
	"aiChat/backend/models"
	"errors"
	"time"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound 用户不存在错误
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials 凭证无效错误
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrEmailExists 邮箱已存在错误
	ErrEmailExists = errors.New("email already exists")

	// ErrPhoneExists 手机号已存在错误
	ErrPhoneExists = errors.New("phone number already exists")

	// ErrUsernameExists 用户名已存在错误
	ErrUsernameExists = errors.New("username already exists")

	// ErrInvalidLoginType 无效的登录类型错误
	ErrInvalidLoginType = errors.New("invalid login type")
)

// RegisterUser 注册新用户
func RegisterUser(req models.RegisterRequest) (uint, error) {
	db := database.GetDB()

	// 验证用户名是否已存在
	var count int64
	if err := db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, ErrUsernameExists
	}

	// 验证邮箱或手机号是否已存在（基于登录类型）
	if req.LoginType == 1 {
		// 邮箱登录
		if req.Email == "" {
			return 0, errors.New("email is required for email login type")
		}
		if err := db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
			return 0, err
		}
		if count > 0 {
			return 0, ErrEmailExists
		}
	} else if req.LoginType == 2 {
		// 手机号登录
		if req.Phone == "" {
			return 0, errors.New("phone is required for phone login type")
		}
		if err := db.Model(&models.User{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
			return 0, err
		}
		if count > 0 {
			return 0, ErrPhoneExists
		}
	} else {
		return 0, ErrInvalidLoginType
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// 使用事务
	var userID uint
	err = db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := models.User{
			Username:  req.Username,
			Password:  string(hashedPassword),
			LoginType: req.LoginType,
		}

		if req.LoginType == 1 {
			user.Email = req.Email
		} else {
			user.Phone = req.Phone
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		userID = user.ID

		// 创建用户设置
		userSettings := models.UserSettings{
			UserID: userID,
		}

		if err := tx.Create(&userSettings).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userID, nil
}

// LoginUser 用户登录
func LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	db := database.GetDB()
	var user models.User

	// 基于登录类型选择查询方式
	query := db.Where("status = ?", 1)

	if req.LoginType == 1 {
		// 邮箱登录
		query = query.Where("email = ?", req.Account)
	} else if req.LoginType == 2 {
		// 手机号登录
		query = query.Where("phone = ?", req.Account)
	} else {
		return nil, ErrInvalidLoginType
	}

	// 查询用户
	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成令牌
	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	if err := db.Model(&user).Updates(map[string]interface{}{
		"last_login_time": &now,
	}).Error; err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID 通过ID获取用户
func GetUserByID(userID uint64) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("id = ? AND status = ?", userID, 1).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 清除密码字段
	user.Password = ""

	return &user, nil
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(userID uint64, req models.UpdateProfileRequest) error {
	db := database.GetDB()

	// 如果更新用户名，检查用户名是否已存在
	if req.Username != "" {
		var count int64
		if err := db.Model(&models.User{}).Where("username = ? AND id != ?", req.Username, userID).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return ErrUsernameExists
		}
	}

	// 准备更新数据
	updates := map[string]interface{}{}

	if req.Username != "" {
		updates["username"] = req.Username
	}

	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	return db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userID uint64, req models.UpdatePasswordRequest) error {
	db := database.GetDB()

	// 获取用户当前密码
	var user models.User
	if err := db.Select("password").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 验证当前密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	// 加密新密码
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return db.Model(&models.User{}).Where("id = ?", userID).Update("password", string(newHashedPassword)).Error
}

// GetUserSettings 获取用户设置
func GetUserSettings(userID uint) (*models.UserSettings, error) {
	db := database.GetDB()
	var settings models.UserSettings

	err := db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有设置，则创建默认设置
			settings = models.UserSettings{
				UserID: userID,
			}

			if err := db.Create(&settings).Error; err != nil {
				return nil, err
			}

			return &settings, nil
		}
		return nil, err
	}

	return &settings, nil
}

// UpdateUserSettings 更新用户设置
func UpdateUserSettings(userID uint64, req models.UpdateSettingsRequest) error {
	db := database.GetDB()

	// 准备更新数据
	updates := map[string]interface{}{}

	if req.Theme != "" {
		updates["theme"] = req.Theme
	}

	if req.Language != "" {
		updates["language"] = req.Language
	}

	notificationEnabled := 0
	if req.NotificationEnabled {
		notificationEnabled = 1
	}
	updates["notification_enabled"] = notificationEnabled

	if req.Prompt != "" {
		updates["prompt"] = req.Prompt
	}

	if req.Rules != "" {
		updates["rules"] = req.Rules
	}

	return db.Model(&models.UserSettings{}).Where("user_id = ?", userID).Updates(updates).Error
}
