-- UUID v7を生成する関数を定義
CREATE OR REPLACE FUNCTION uuid_generate_v7() RETURNS uuid AS $$
DECLARE unix_time_ms bytea;
uuid_bytes bytea;
BEGIN -- 現在時刻をミリ秒で取得
unix_time_ms := decode(
  lpad(
    to_hex(
      floor(
        extract(
          epoch
          from clock_timestamp()
        ) * 1000
      )::bigint
    ),
    12,
    '0'
  ),
  'hex'
);
-- ランダムなバイトを生成
uuid_bytes := unix_time_ms || gen_random_bytes(10);
-- UUID v7のビット設定 (version 7, variant 1)
uuid_bytes := set_byte(
  uuid_bytes,
  6,
  (get_byte(uuid_bytes, 6) & 15) | 112
);
uuid_bytes := set_byte(
  uuid_bytes,
  8,
  (get_byte(uuid_bytes, 8) & 63) | 128
);
RETURN encode(uuid_bytes, 'hex')::uuid;
END $$ LANGUAGE plpgsql VOLATILE;
CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v7(),
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255),
  image_url TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE github_account (
  user_id uuid NOT NULL,
  provider VARCHAR NOT NULL,
  provider_account_id VARCHAR NOT NULL,
  access_token VARCHAR,
  refresh_token VARCHAR,
  expires_at BIGINT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT github_account_pk PRIMARY KEY (provider, provider_account_id),
  CONSTRAINT github_account_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE google_account (
  user_id uuid NOT NULL,
  provider VARCHAR NOT NULL,
  provider_account_id VARCHAR NOT NULL,
  access_token VARCHAR,
  refresh_token VARCHAR,
  expires_at BIGINT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT google_account_pk PRIMARY KEY (provider, provider_account_id),
  CONSTRAINT google_account_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE project (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v7(),
  user_id uuid NOT NULL,
  title VARCHAR NOT NULL,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT project_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE task (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v7(),
  project_id uuid NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status INT NOT NULL,
  end_date TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT task_project_fk FOREIGN KEY (project_id) REFERENCES project(id)
);