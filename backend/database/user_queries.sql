-- 用户注册 (邮箱)
INSERT INTO users (username, password, email, login_type) 
VALUES (?, ?, ?, 1);

-- 用户注册 (手机号)
INSERT INTO users (username, password, phone, login_type) 
VALUES (?, ?, ?, 2);

-- 用户注册后创建默认设置
INSERT INTO user_settings (user_id) 
VALUES (?);

-- 通过邮箱查询用户
SELECT * FROM users 
WHERE email = ? AND status = 1 
LIMIT 1;

-- 通过手机号查询用户
SELECT * FROM users 
WHERE phone = ? AND status = 1 
LIMIT 1;

-- 通过用户ID查询用户
SELECT * FROM users 
WHERE id = ? AND status = 1 
LIMIT 1;

-- 通过用户名查询用户
SELECT * FROM users 
WHERE username = ? AND status = 1 
LIMIT 1;

-- 验证邮箱是否已存在
SELECT COUNT(*) AS count FROM users 
WHERE email = ?;

-- 验证手机号是否已存在
SELECT COUNT(*) AS count FROM users 
WHERE phone = ?;

-- 验证用户名是否已存在
SELECT COUNT(*) AS count FROM users 
WHERE username = ?;

-- 更新用户登录信息
UPDATE users 
SET last_login_time = ?, last_login_ip = ? 
WHERE id = ?;

-- 创建用户会话
INSERT INTO user_sessions (user_id, token, expire_time) 
VALUES (?, ?, ?);

-- 查询用户会话
SELECT * FROM user_sessions 
WHERE token = ? AND expire_time > NOW() 
LIMIT 1;

-- 删除用户会话
DELETE FROM user_sessions 
WHERE token = ?;

-- 删除用户过期会话
DELETE FROM user_sessions 
WHERE expire_time <= NOW();

-- 更新用户密码
UPDATE users 
SET password = ? 
WHERE id = ?;

-- 更新用户信息
UPDATE users 
SET username = ?, avatar = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- 更新用户邮箱
UPDATE users 
SET email = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- 更新用户手机号
UPDATE users 
SET phone = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- 获取用户设置
SELECT * FROM user_settings 
WHERE user_id = ? 
LIMIT 1;

-- 更新用户设置
UPDATE user_settings 
SET theme = ?, language = ?, notification_enabled = ? 
WHERE user_id = ?;

-- 禁用用户
UPDATE users 
SET status = 0, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- 启用用户
UPDATE users 
SET status = 1, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?; 