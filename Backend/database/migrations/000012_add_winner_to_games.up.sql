BEGIN;

-- 為 games 表新增 winner 欄位
ALTER TABLE games
    ADD COLUMN winner VARCHAR(50) NULL DEFAULT '' AFTER current_players;

COMMIT;
