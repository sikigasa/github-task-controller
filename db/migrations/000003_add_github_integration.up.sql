-- github_accountにPATカラムを追加
ALTER TABLE github_account ADD COLUMN IF NOT EXISTS pat_encrypted VARCHAR;

-- プロジェクトにGitHub連携情報を追加
ALTER TABLE project ADD COLUMN IF NOT EXISTS github_owner VARCHAR;
ALTER TABLE project ADD COLUMN IF NOT EXISTS github_repo VARCHAR;
ALTER TABLE project ADD COLUMN IF NOT EXISTS github_project_number INT;

-- タスクにGitHub連携情報を追加
ALTER TABLE task ADD COLUMN IF NOT EXISTS github_item_id VARCHAR;
ALTER TABLE task ADD COLUMN IF NOT EXISTS github_issue_number INT;
ALTER TABLE task ADD COLUMN IF NOT EXISTS github_issue_url VARCHAR;
