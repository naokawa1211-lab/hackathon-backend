-- 文字コードを明示的に指定（Cloud SQL対策・絵文字対応）
SET CHARACTER SET utf8mb4;

-- 1. ユーザーテーブル (userdatabase)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(128) PRIMARY KEY, -- FirebaseのUIDを入れる
    username VARCHAR(50) NOT NULL,
    avatar_url TEXT,
    space_credits INT DEFAULT 10000, -- 通貨単位: cr
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. 出品商品テーブル (productdatabase)
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    category VARCHAR(50),

    -- 商品画像URL
    image_url_1 VARCHAR(512) NOT NULL,      -- メイン画像（必須）
    image_url_2 VARCHAR(512),               -- サブ画像（任意）
       
    seller_id VARCHAR(128),
    status VARCHAR(20) DEFAULT 'available', -- available, sold
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 🕒 【追加】更新日時（商品が 'sold' になった時間を追うためにあると便利）
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
    
    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. メッセージテーブル (messagedatabase)
CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id VARCHAR(128),
    receiver_id VARCHAR(128),
    content TEXT NOT NULL, -- 宇宙人語など
    translated_content TEXT, -- Gemini APIで翻訳したテキストを入れる用
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. 取引テーブル (transactiondatabase)
CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT,
    buyer_id VARCHAR(128),
    seller_id VARCHAR(128),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- 🌌 🚀 【テスト用モックデータの注入】 🚀 🌌
INSERT INTO users (id, username, space_credits) VALUES
('mock_uid_shibata', 'シバタ大佐', 5000000),
('mock_uid_ecoearth', 'エコアース社', 12000),
('mock_uid_dr_mase', 'Dr.マセ', 85000),
('mock_uid_naoya', 'Naoya', 30000);

INSERT INTO messages (sender_id, receiver_id, content, translated_content) VALUES
('mock_uid_shibata', 'mock_uid_naoya', '👽 ギガ・グラビティ！！（探査艇、値引きできる？）', '（翻訳：こんにちは。出品されている探査艇ですが、400万crに値下げしていただくことは可能でしょうか？）');