-- API密钥表结构定义
CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    provider_type VARCHAR(50) NOT NULL,
    encrypted_key TEXT NOT NULL,
    key_hash VARCHAR(64) NOT NULL, -- 用于验证密钥是否已存在，不存储明文
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, provider_type) -- 每个用户每个提供商只能有一个密钥
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_provider_type ON api_keys(provider_type);
CREATE INDEX IF NOT EXISTS idx_api_keys_user_provider ON api_keys(user_id, provider_type);
CREATE INDEX IF NOT EXISTS idx_api_keys_is_active ON api_keys(is_active);
CREATE INDEX IF NOT EXISTS idx_api_keys_created_at ON api_keys(created_at);

-- 创建触发器自动更新 updated_at 字段
CREATE TRIGGER IF NOT EXISTS update_api_keys_updated_at 
    AFTER UPDATE ON api_keys
    FOR EACH ROW
BEGIN
    UPDATE api_keys SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;