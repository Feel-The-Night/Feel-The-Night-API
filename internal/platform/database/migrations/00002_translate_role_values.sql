-- +goose Up
ALTER TABLE users DROP CONSTRAINT chk_users_role;

UPDATE users SET role = CASE role
    WHEN 'iniciante' THEN 'beginner'
    WHEN 'veterano'  THEN 'veteran'
    ELSE role
END;

ALTER TABLE users ALTER COLUMN role SET DEFAULT 'beginner';
ALTER TABLE users ADD CONSTRAINT chk_users_role CHECK (role IN ('beginner', 'veteran', 'admin'));

-- +goose Down
ALTER TABLE users DROP CONSTRAINT chk_users_role;

UPDATE users SET role = CASE role
    WHEN 'beginner' THEN 'iniciante'
    WHEN 'veteran'  THEN 'veterano'
    ELSE role
END;

ALTER TABLE users ALTER COLUMN role SET DEFAULT 'iniciante';
ALTER TABLE users ADD CONSTRAINT chk_users_role CHECK (role IN ('iniciante', 'veterano', 'admin'));
