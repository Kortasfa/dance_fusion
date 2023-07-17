USE dance_fusion;
CREATE TABLE faces(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
    `recommended_level` INTEGER NOT NULL,
    `face_src` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;