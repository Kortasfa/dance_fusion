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
    ("good_4_u", "static/motion_list/good_4_u.json" ,"static/motion_list/good_4_u.json", "static/motion_list/good_4_u.json", "static/motion_list/good_4_u.json"),
    ("Rasputin", "static/motion_list/rasputin.json" ,"static/motion_list/rasputin.json", "static/motion_list/rasputin.json", "static/motion_list/rasputin.json"),
    ("mi_mi_mi", "static/motion_list/mi_mi_mi.json" ,"static/motion_list/mi_mi_mi.json", "static/motion_list/mi_mi_mi.json", "static/motion_list/mi_mi_mi.json");