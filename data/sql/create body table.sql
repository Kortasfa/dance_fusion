USE dance_fusion;
CREATE TABLE bodies(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
    `recommended_level` INTEGER NOT NULL,
    `body_src` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO bodies (id, recommended_level, body_src) VALUES
(1, 0, '../static/img/body/greenBody.png'),
(2, 0, '../static/img/body/7cf79c.png'),
(3, 0, '../static/img/body/50a7ea.png'),
(4, 1, '../static/img/body/89e7dc.png'),
(5, 1, '../static/img/body/ce63ff.png'),
(6, 1, '../static/img/body/fe73ac.png'),
(7, 2, '../static/img/body/fe5955.png'),
(8, 2, '../static/img/body/ff9f55.png'),
(9, 3, '../static/img/body/ff92ff.png'),
(10, 3, '../static/img/body/fff76e.png');