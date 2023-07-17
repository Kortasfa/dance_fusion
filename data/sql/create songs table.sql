USE dance_fusion;
CREATE TABLE songs(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `song_name` VARCHAR(255) NOT NULL,
    `author_name` VARCHAR(255) NOT NULL,
    `image_src` VARCHAR(255) NOT NULL,
    `video_src` VARCHAR(255) NOT NULL,
    `preview_video_src` VARCHAR(255) NOT NULL,
    `style_id` INTEGER NOT NULL,
    `motion_list_path` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO songs(song_name, author_name, image_src, video_src, preview_video_src, style_id, motion_list_path)
VALUES
	("Forget You", "ceeLoGreen", "static/img/forgetYou.jpg", "static/video/forgetYou.mp4", "static/video/forgetYou.mp4", 3, "static/motion_list/forgetYou.json"),
    ("1999", "charlXCX&troyeSivan", "static/img/1999.jpg", "static/video/1999.mp4", "static/video/1999.mp4", 1, "static/motion_list/forgetYou.json"),
    ("American Girl", "bonnieMcKeeSong", "static/img/americanGirl.jpg", "static/video/americanGirl.mp4", "static/video/americanGirl.mp4", 2, "static/motion_list/forgetYou.json"),
    ("Burn", "ellieGoulding", "static/img/burn.jpg", "static/video/burn.mp4", "static/video/burn.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Firework", "katyPerry", "static/img/firework.jpg", "static/video/firework.mp4", "static/video/firework.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Kiss Me More", "dojaCatFtSza", "static/img/kissMeMore.jpg", "static/video/kissMeMore.mp4", "static/video/kissMeMore.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Dynamite", "ceeLoGreen", "static/img/dynamite.jpg", "static/video/dynamite.mp4", "static/video/dynamite.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Oops I Did It Again", "britneySpears", "static/img/oopsIDidItAgain.jpg", "static/video/ooopsIDidItAgain.mp4", "static/video/ooopsIDidItAgain.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Montero", "lilNasX", "static/img/montero.jpg", "static/video/montero.mp4", "static/video/montero.mp4", 1, "static/motion_list/forgetYou.json");