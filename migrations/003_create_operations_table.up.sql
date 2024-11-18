CREATE TABLE operations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type ENUM('addition', 'subtraction', 'multiplication', 'division', 'square_root', 'random_string') UNIQUE NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    cost FLOAT NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);
