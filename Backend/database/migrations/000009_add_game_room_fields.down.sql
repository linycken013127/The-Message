BEGIN;

DROP TABLE IF EXISTS game_players;

ALTER TABLE games
    DROP COLUMN host_account_id,
    DROP COLUMN max_players,
    DROP COLUMN current_players;

COMMIT;
