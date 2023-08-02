USE dance_fusion;
CREATE TABLE bots(
                     `bot_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                     `bot_name` VARCHAR(255) NOT NULL,
                     `bot_scores_path` VARCHAR(255) NOT NULL,
                     `img_hat` VARCHAR(255) NOT NULL,
                     `img_body` VARCHAR(255) NOT NULL,
                     `img_face` VARCHAR(255) NOT NULL,
                     `difficulty` INTEGER NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO bots(bot_id, bot_name, bot_scores_path, img_hat, img_body, img_face, difficulty)
VALUES
(1, 'Maxim', 'static/bots_list/bot_1.json', 'static/img/hat/robinHoodHat.png', 'static/img/body/fff76e.png', 'static/img/face/7.png', 4),
(2, 'Kirill', 'static/bots_list/bot_2.json', 'static/img/hat/bigHat.png', 'static/img/body/greenBody.png', '../static/img/face/pickUpFace.png', 3),
(3, 'Zaxar', 'static/bots_list/bot_2.json', 'static/img/hat/cowboyHat.png', 'static/img/body/89e7dc.png', '../static/img/face/1.png', 2),
(4, 'Valeryan', 'static/bots_list/bot_2.json', 'static/img/hat/stanHat.png', 'static/img/body/50a7ea.png', '../static/img/face/badFace.png', 1),
(5, 'Dasha', 'static/bots_list/bot_2.json', 'static/img/hat/herHat.png', 'static/img/body/fe5955.png', '../static/img/face/angriestFace.png', 1);