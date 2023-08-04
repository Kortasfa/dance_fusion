USE dance_fusion;

CREATE TABLE bosses(
                     `boss_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                     `boss_name` VARCHAR(255) NOT NULL,
                     `boss_health_point` INTEGER NOT NULL,
                     `img_hat` VARCHAR(255) NOT NULL,
                     `img_body` VARCHAR(255) NOT NULL,
                     `img_face` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO bosses(boss_name, boss_health_point, img_hat, img_body, img_face)
VALUES
    ("Rhythmo", 3200, "static/img/hat/robinHoodHat.png", "static/img/body/fff76e.png", "static/img/face/7.png"),
    ("Electric Maestro", 3600, "static/img/hat/robinHoodHat.png", "static/img/body/fff76e.png", "static/img/face/7.png"),
    ("Funk Fusionist", 4000, "static/img/hat/robinHoodHat.png", "static/img/body/fff76e.png", "static/img/face/7.png"),
    ("Tempo Titan", 4500, "static/img/hat/robinHoodHat.png", "static/img/body/fff76e.png", "static/img/face/7.png");

