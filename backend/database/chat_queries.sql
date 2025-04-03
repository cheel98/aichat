-- 创建新的聊天会话
INSERT INTO chat_sessions (session_id, user_id, title) 
VALUES (?, ?, ?);

-- 获取用户的聊天会话列表
SELECT * FROM chat_sessions 
WHERE user_id = ? 
ORDER BY is_pinned DESC, updated_at DESC;

-- 获取聊天会话详情
SELECT * FROM chat_sessions 
WHERE session_id = ? AND user_id = ? 
LIMIT 1;

-- 更新聊天会话标题
UPDATE chat_sessions 
SET title = ?, updated_at = CURRENT_TIMESTAMP 
WHERE session_id = ? AND user_id = ?;

-- 设置/取消置顶聊天会话
UPDATE chat_sessions 
SET is_pinned = ?, updated_at = CURRENT_TIMESTAMP 
WHERE session_id = ? AND user_id = ?;

-- 删除聊天会话
DELETE FROM chat_sessions 
WHERE session_id = ? AND user_id = ?;

-- 保存聊天消息
INSERT INTO chat_messages (user_id, session_id, role, content) 
VALUES (?, ?, ?, ?);

-- 获取会话的聊天历史记录
SELECT * FROM chat_messages 
WHERE session_id = ? AND user_id = ? 
ORDER BY created_at ASC;

-- 删除会话的聊天历史记录
DELETE FROM chat_messages 
WHERE session_id = ? AND user_id = ?;

-- 获取最近的聊天会话
SELECT cs.* FROM chat_sessions cs 
LEFT JOIN (
  SELECT session_id, MAX(created_at) as latest_msg_time 
  FROM chat_messages 
  GROUP BY session_id
) cm ON cs.session_id = cm.session_id 
WHERE cs.user_id = ? 
ORDER BY cs.is_pinned DESC, cm.latest_msg_time DESC, cs.updated_at DESC 
LIMIT ?;

-- 搜索聊天历史记录
SELECT cm.*, cs.title FROM chat_messages cm 
INNER JOIN chat_sessions cs ON cm.session_id = cs.session_id 
WHERE cm.user_id = ? AND cm.content LIKE ? 
ORDER BY cm.created_at DESC 
LIMIT ?, ?;

-- 获取每个会话的消息数量
SELECT session_id, COUNT(*) as message_count 
FROM chat_messages 
WHERE user_id = ? 
GROUP BY session_id;

-- 获取用户的总消息数量
SELECT COUNT(*) as total_messages 
FROM chat_messages 
WHERE user_id = ?;

-- 获取用户最新的20条消息
SELECT * FROM chat_messages 
WHERE user_id = ? 
ORDER BY created_at DESC 
LIMIT 20; 