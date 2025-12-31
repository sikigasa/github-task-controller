ALTER TABLE task DROP COLUMN github_issue_url;
ALTER TABLE task DROP COLUMN github_issue_number;
ALTER TABLE task DROP COLUMN github_item_id;

ALTER TABLE project DROP COLUMN github_project_number;
ALTER TABLE project DROP COLUMN github_repo;
ALTER TABLE project DROP COLUMN github_owner;

ALTER TABLE github_account DROP COLUMN pat_encrypted;
