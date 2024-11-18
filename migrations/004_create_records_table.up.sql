CREATE TABLE records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    operation_id INT,
    user_id INT,
    amount FLOAT NOT NULL,
    user_balance FLOAT NOT NULL,
    operation_response TEXT,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_operation_id FOREIGN KEY (operation_id) REFERENCES operations(id) ON DELETE SET NULL,
    CONSTRAINT fk_user_id_records FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
