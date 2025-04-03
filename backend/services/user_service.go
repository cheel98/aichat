package services

import (
	"aiChat/backend/database"
	"aiChat/backend/models"
	"database/sql"
	"errors"

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
func RegisterUser(req models.RegisterRequest) (uint64, error) {
	// 验证用户名是否已存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
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
		err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
		if err != nil {
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
		err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phone = ?", req.Phone).Scan(&count)
		if err != nil {
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

	// 开始事务
	tx, err := database.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var result sql.Result
	var userID int64

	// 插入用户记录
	if req.LoginType == 1 {
		// 邮箱注册
		result, err = tx.Exec(
			"INSERT INTO users (username, password, email, login_type) VALUES (?, ?, ?, ?)",
			req.Username, hashedPassword, req.Email, req.LoginType,
		)
	} else {
		// 手机号注册
		result, err = tx.Exec(
			"INSERT INTO users (username, password, phone, login_type) VALUES (?, ?, ?, ?)",
			req.Username, hashedPassword, req.Phone, req.LoginType,
		)
	}

	if err != nil {
		return 0, err
	}

	// 获取用户ID
	userID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 创建用户设置
	_, err = tx.Exec("INSERT INTO user_settings (user_id) VALUES (?)", userID)
	if err != nil {
		return 0, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(userID), nil
}

// LoginUser 用户登录
func LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User
	var row *sql.Row

	// 基于登录类型选择查询方式
	if req.LoginType == 1 {
		// 邮箱登录
		row = database.DB.QueryRow("SELECT * FROM users WHERE email = ? AND status = 1 LIMIT 1", req.Account)
	} else if req.LoginType == 2 {
		// 手机号登录
		row = database.DB.QueryRow("SELECT * FROM users WHERE phone = ? AND status = 1 LIMIT 1", req.Account)
	} else {
		return nil, ErrInvalidLoginType
	}

	// 扫描用户数据
	err := row.Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Avatar, &user.Status, &user.LoginType, &user.LastLoginTime,
		&user.LastLoginIP, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成令牌
	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	_, err = database.DB.Exec(
		"UPDATE users SET last_login_time = NOW() WHERE id = ?",
		user.ID,
	)
	if err != nil {
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
	var user models.User

	err := database.DB.QueryRow(
		"SELECT * FROM users WHERE id = ? AND status = 1 LIMIT 1",
		userID,
	).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Avatar, &user.Status, &user.LoginType, &user.LastLoginTime,
		&user.LastLoginIP, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
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
	// 如果更新用户名，检查用户名是否已存在
	if req.Username != "" {
		var count int
		err := database.DB.QueryRow(
			"SELECT COUNT(*) FROM users WHERE username = ? AND id != ?",
			req.Username, userID,
		).Scan(&count)

		if err != nil {
			return err
		}

		if count > 0 {
			return ErrUsernameExists
		}
	}

	// 准备更新查询
	query := "UPDATE users SET updated_at = NOW()"
	params := []interface{}{}

	if req.Username != "" {
		query += ", username = ?"
		params = append(params, req.Username)
	}

	if req.Avatar != "" {
		query += ", avatar = ?"
		params = append(params, req.Avatar)
	}

	query += " WHERE id = ?"
	params = append(params, userID)

	// 执行更新
	_, err := database.DB.Exec(query, params...)
	return err
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userID uint64, req models.UpdatePasswordRequest) error {
	// 获取当前密码
	var hashedPassword string
	err := database.DB.QueryRow(
		"SELECT password FROM users WHERE id = ?",
		userID,
	).Scan(&hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.OldPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	// 加密新密码
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	_, err = database.DB.Exec(
		"UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?",
		newHashedPassword, userID,
	)
	return err
}

// GetUserSettings 获取用户设置
func GetUserSettings(userID uint64) (*models.UserSettings, error) {
	var settings models.UserSettings

	err := database.DB.QueryRow(
		"SELECT * FROM user_settings WHERE user_id = ? LIMIT 1",
		userID,
	).Scan(
		&settings.ID, &settings.UserID, &settings.Theme, &settings.Language,
		&settings.NotificationEnabled, &settings.CreatedAt, &settings.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有设置，则创建默认设置
			_, err = database.DB.Exec(
				"INSERT INTO user_settings (user_id) VALUES (?)",
				userID,
			)
			if err != nil {
				return nil, err
			}

			// 重新获取设置
			return GetUserSettings(userID)
		}
		return nil, err
	}

	return &settings, nil
}

// UpdateUserSettings 更新用户设置
func UpdateUserSettings(userID uint64, req models.UpdateSettingsRequest) error {
	// 准备更新查询
	query := "UPDATE user_settings SET updated_at = NOW()"
	params := []interface{}{}

	if req.Theme != "" {
		query += ", theme = ?"
		params = append(params, req.Theme)
	}

	if req.Language != "" {
		query += ", language = ?"
		params = append(params, req.Language)
	}

	notificationEnabled := 0
	if req.NotificationEnabled {
		notificationEnabled = 1
	}
	query += ", notification_enabled = ?"
	params = append(params, notificationEnabled)

	query += " WHERE user_id = ?"
	params = append(params, userID)

	// 执行更新
	_, err := database.DB.Exec(query, params...)
	return err
}
