BEGIN;

-- 建立 intelligence_transfers 表
CREATE TABLE intelligence_transfers
(
    id                        INT AUTO_INCREMENT PRIMARY KEY,
    game_id                   INT          NOT NULL,
    card_id                   INT          NOT NULL,
    sender_player_id          INT          NOT NULL,
    current_target_player_id  INT          NOT NULL,
    original_target_player_id INT          NOT NULL,
    face_up                   BOOLEAN      NOT NULL DEFAULT FALSE,
    status                    VARCHAR(20)  NOT NULL DEFAULT 'IN_TRANSIT',
    created_at                DATETIME     NOT NULL,
    updated_at                DATETIME     NOT NULL,
    deleted_at                DATETIME     NULL,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES cards (id) ON DELETE CASCADE,
    FOREIGN KEY (sender_player_id) REFERENCES players (id) ON DELETE CASCADE,
    FOREIGN KEY (current_target_player_id) REFERENCES players (id) ON DELETE CASCADE,
    FOREIGN KEY (original_target_player_id) REFERENCES players (id) ON DELETE CASCADE,
    INDEX idx_game_status (game_id, status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

COMMIT;
