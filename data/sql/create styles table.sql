USE dance_fusion;
CREATE TABLE  styles(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `name` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO styles(name)
VALUES
	("Hip-hop"),
    ("Latin music"),
    ("Russian music");