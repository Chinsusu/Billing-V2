ALTER TABLE auth_sessions
    DROP COLUMN IF EXISTS two_factor_satisfied_at;

DROP TABLE IF EXISTS user_two_factor_methods;
