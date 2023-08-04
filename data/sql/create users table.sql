USE dance_fusion;

CREATE TABLE users(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `name` VARCHAR(255) NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `img_hat` VARCHAR(255) NOT NULL  DEFAULT 'static/img/hat/withoutHat.png',
    `img_body` VARCHAR(255) NOT NULL  DEFAULT 'static/img/body/fff76e.png',
    `img_face` VARCHAR(255) NOT NULL  DEFAULT 'static/img/face/2.png',
    `total_score` INTEGER NOT NULL DEFAULT 0
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;