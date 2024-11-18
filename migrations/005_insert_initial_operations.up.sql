INSERT INTO operations (type, cost, status) VALUES
    ('addition', 20.0, 'active'),
    ('subtraction', 20.0, 'active'),
    ('multiplication', 30.0, 'active'),
    ('division', 40.0, 'active'),
    ('square_root', 25.0, 'active'),
    ('random_string', 50.0, 'active')
ON DUPLICATE KEY UPDATE cost = VALUES(cost), status = VALUES(status);
