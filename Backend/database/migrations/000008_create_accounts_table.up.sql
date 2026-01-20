BEGIN;

CREATE TABLE accounts
(
  id         INT AUTO_INCREMENT PRIMARY KEY,
  name       VARCHAR(50)  NOT NULL,
  created_at DATETIME     NOT NULL,
  updated_at DATETIME     NOT NULL,
  deleted_at DATETIME,
  UNIQUE KEY uk_accounts_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
