USE dance_fusion;
CREATE TABLE songs(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `song_name` VARCHAR(255) NOT NULL,
    `author_name` VARCHAR(255) NOT NULL,
    `image_src` VARCHAR(255) NOT NULL,
    `video_src` VARCHAR(255) NOT NULL,
    `preview_video_src` VARCHAR(255) NOT NULL,
    `style_id` INTEGER NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO songs(song_name, author_name, image_src, video_src, preview_video_src, style_id)
VALUES
	("Forget You", "ceeLoGreen", "static/img/forgetYou.jpg", "static/video/forgetYou.mp4", "static/video/forgetYou.mp4", 3),
    ("1999", "charlXCX&troyeSivan", "static/img/1999.jpg", "static/video/1999.mp4", "static/video/1999.mp4", 1),
    ("American Girl", "bonnieMcKeeSong", "static/img/americanGirl.jpg", "static/video/americanGirl.mp4", "static/video/americanGirl.jpg", 2),
    ("Burn", "ellieGoulding", "static/img/burn.jpg", "static/video/burn.mp4", "static/video/burn.mp4", 1),
    ("Firework", "katyPerry", "static/img/firework.jpg", "static/video/firework.mp4", "static/video/firework.mp4", 1),
    ("Kiss Me More", "dojaCatFtSza", "static/img/kissMeMore.jpg", "static/video/kissMeMore.mp4", "static/video/kissMeMore.mp4", 1),
    ("Dynamite", "ceeLoGreen", "static/img/dynamite.jpg", "static/video/dynamite.mp4", "static/video/dynamite.mp4", 1),
    ("Oops I Did It Again", "britneySpears", "static/img/oopsIDidItAgain.jpg", "static/video/ooopsIDidItAgain.mp4", "static/video/ooopsIDidItAgain.mp4", 1),
    ("Montero", "lilNasX", "static/img/montero.jpg", "static/video/montero.mp4", "static/video/montero.mp4", 1);