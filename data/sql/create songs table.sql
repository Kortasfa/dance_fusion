USE dance_fusion;

CREATE TABLE songs(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `song_name` VARCHAR(255) NOT NULL,
    `author_name` VARCHAR(255) NOT NULL,
    `image_src` VARCHAR(255) NOT NULL,
    `video_src` VARCHAR(255) NOT NULL,
    `preview_video_src` VARCHAR(255) NOT NULL,
    `style_id` INTEGER NOT NULL,
    `motion_list_path` VARCHAR(255) NOT NULL,
    `difficulty` INTEGER NOT NULL,
    `best_player_id` INTEGER,
    `best_score` INTEGER
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO songs(song_name, author_name, image_src, video_src, preview_video_src, style_id, motion_list_path, difficulty)
VALUES
	("Forget You", "Cee Lo Green", "static/img/forgetYou.jpg", 'static/video/full/forgetYouFull.mp4', "static/video/previews/forgetYou.mp4", 3, "static/motion_list/forgetYou.json", 2),
    ("Firework", "katy Perry", "static/img/firework.jpg", "static/video/firework.mp4", "static/video/previews/firework.mp4", 1, "static/motion_list/forgetYou.json", 4),
    ("Kiss Me More", "Doja Cat Ft. Sza", "static/img/kissMeMore.jpg", "static/video/kissMeMore.mp4", "static/video/previews/kissMeMore.mp4", 1, "static/motion_list/forgetYou.json", 1),
    ("Rasputin"," 'Boney M'", "static/img/rasputin.jpg", "static/video/full/rasputin.mp4", "static/video/previews/rasputin.mp4", 3, "static/motion_list/rasputin.json", 2);

INSERT INTO songs(song_name, author_name, image_src, video_src, preview_video_src, style_id, motion_list_path, difficulty, best_player_id, best_score)
VALUES
    ("American Girl", "Bonnie McKee", "static/img/americanGirl.jpg", "static/video/americanGirl.mp4", "static/video/previews/americanGirl.mp4", 2, "static/motion_list/forgetYou.json", 3, 3, 100);