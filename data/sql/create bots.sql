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

INSERT INTO bots(bot_name, bot_scores_path, img_hat, img_body, img_face, difficulty)
VALUES
    ('Maxim', 'static/bots_list/bot_charged_maksim.json', 'static/img/hat/withoutHat.png', 'static/img/special/ronaldDrink.gif', 'static/img/face/withoutFace.png', 3),
    ('Kirill', 'static/bots_list/bot_fury_kirill.json', 'static/img/hat/cilindrHat.png', 'static/img/body/greenBody.png', '../static/img/face/loveFace.png', 3),
    ('Zaxar', 'static/bots_list/bot_smart_zaxar.json', 'static/img/hat/withoutHat.png', 'static/img/special/papichBarMontage.gif', '../static/img/face/withoutFace.png', 2),
    ('Valeryan', 'static/bots_list/bot_up_down_hand_valeryan.json', 'static/img/hat/withoutHat.png', 'static/img/special/siBody.png', '../static/img/face/withoutFace.png', 1),
    ('Dasha', 'static/bots_list/bot_super_dasha.json', 'static/img/hat/withoutHat.png', 'static/img/special/cat.gif', 'static/img/face/withoutFace.png', 1),
    ('Michael Jackson', 'static/bots_list/bot_michael_jackson.json', 'static/img/hat/withoutHat.png', 'static/img/special/michaelJackson.png', '../static/img/face/withoutFace.png', 4),
    ('Kerim', 'static/bots_list/bot_kerim.json', 'static/img/hat/beretHat.png', 'static/img/body/fe5955.png', '../static/img/face/pickUpFace.png', 4);
select * from bots;





