-- タスクにpriorityカラムを追加 (0: Low, 1: Medium, 2: High)
ALTER TABLE task ADD COLUMN IF NOT EXISTS priority INT NOT NULL DEFAULT 1;
