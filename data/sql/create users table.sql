USE dance_fusion;
CREATE TABLE users(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `name` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `img_src` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO users(name, password, img_src)
VALUES
	("name", "password", "static/img/user_1.png"),
    ("name2", "password2", "static/img/user_2.png");