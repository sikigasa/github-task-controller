ALTER TABLE task DROP COLUMN IF EXISTS github_issue_url;
ALTER TABLE task DROP COLUMN IF EXISTS github_issue_number;
ALTER TABLE task DROP COLUMN IF EXISTS github_item_id;

ALTER TABLE project DROP COLUMN IF EXISTS github_project_number;
ALTER TABLE project DROP COLUMN IF EXISTS github_repo;
ALTER TABLE project DROP COLUMN IF EXISTS github_owner;

ALTER TABLE github_account DROP COLUMN IF EXISTS pat_encrypted;
