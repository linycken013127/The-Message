BEGIN;

-- 為 games 表新增 phase 欄位
ALTER TABLE games
    ADD COLUMN phase VARCHAR(20) NOT NULL DEFAULT 'ACTION' AFTER status;

-- 建立 action_passes 表
CREATE TABLE action_passes
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    game_id    INT      NOT NULL,
    player_id  INT      NOT NULL,
    round      INT      NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE CASCADE,
    UNIQUE KEY uk_game_player_round (game_id, player_id, round)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

COMMIT;
