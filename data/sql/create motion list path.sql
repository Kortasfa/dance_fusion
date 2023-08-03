USE dance_fusion;
CREATE TABLE motion_list_path(
                                 `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                                 `song_name` VARCHAR(255) NOT NULL,
                                 `first_player` VARCHAR(255) NOT NULL,
                                 `second_player` VARCHAR(255) NOT NULL,
                                 `third_player` VARCHAR(255) NOT NULL,
                                 `fourth_player` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;
INSERT INTO motion_list_path(song_name, first_player, second_player, third_player, fourth_player)
VALUES
    ("Forget You", "static/motion_list/forgetYou.json", "static/motion_list/forgetYou.json", "static/motion_list/forgetYou.json", "static/motion_list/forgetYou.json"),
    ("Good 4 u", "static/motion_list/good_4_u.json" ,"static/motion_list/good_4_u.json", "static/motion_list/good_4_u.json", "static/motion_list/good_4_u.json"),
    ("Rasputin", "static/motion_list/rasputin.json" ,"static/motion_list/rasputin.json", "static/motion_list/rasputin.json", "static/motion_list/rasputin.json"),
    ("Mi Mi Mi", "static/motion_list/mimimi-man.json" ,"static/motion_list/mimimi-woman.json", "static/motion_list/mimimi-man.json", "static/motion_list/mimimi-woman.json");