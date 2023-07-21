USE dance_fusion;
CREATE TABLE faces(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
    `recommended_level` INTEGER NOT NULL,
    `face_src` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO faces (id, recommended_level, face_src) VALUES
(1, 0, '../static/img/face/7.png'),
(2, 0, '../static/img/face/6.png'),
(3, 1, '../static/img/face/5.png'),
(4, 2, '../static/img/face/4.png'),
(5, 3, '../static/img/face/3.png'),
(6, 4, '../static/img/face/2.png'),
(7, 4, '../static/img/face/1.png'),
(8, 10, '../static/img/hat/smileFace.png');
