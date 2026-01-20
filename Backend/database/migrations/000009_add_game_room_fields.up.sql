BEGIN;

-- 為 games 表新增遊戲房相關欄位
ALTER TABLE games
    ADD COLUMN host_account_id INT NULL AFTER current_player_id,
    ADD COLUMN max_players INT NOT NULL DEFAULT 9 AFTER host_account_id,
    ADD COLUMN current_players INT NOT NULL DEFAULT 0 AFTER max_players;

-- 建立 game_players 關聯表
CREATE TABLE game_players
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    game_id    INT      NOT NULL,
    account_id INT      NOT NULL,
    join_order INT      NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    UNIQUE KEY uk_game_account (game_id, account_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

COMMIT;
