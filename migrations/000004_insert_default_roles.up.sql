INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator role with full access'),
    ('user', 'Regular user role')
ON CONFLICT (name) DO NOTHING;
