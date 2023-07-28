USE dance_fusion;
CREATE TABLE bots(
                     `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                     `bot_name` VARCHAR(255) NOT NULL,
                     `bot_scores_path` VARCHAR(255) NOT NULL,
                     `img_hat` VARCHAR(255) NOT NULL,
                     `img_body` VARCHAR(255) NOT NULL,
                     `img_face` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO bots(bot_name, bot_scores_path, img_hat, img_body, img_face)
VALUES
    ("bot_1", "static/bots_list/bot_1.json", "static/img/hat/robinHoodHat.png", "static/img/body/fff76e.png", "static/img/face/7.png"),
    ("bot_2", "static/bots_list/bot_2.json", "static/img/hat/robinHoodHat.png", "static/img/body/fff76e.png", "static/img/face/7.png");
SELECT * FROM bots;
