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
	("Forget You", "Cee Lo Green", "static/img/forgetYou.jpg", 'static/video/full/forgetYouFull.mp4', "static/video/previews/forgetYou.mp4", 3, "static/motion_list/forgetYou.json"),
    ("1999", "Charli XCX & Troye Sivan", "static/img/1999.jpg", "static/video/1999.mp4", "static/video/previews/1999.mp4", 1, "static/motion_list/raRaRasputin.json"),
    ("American Girl", "Bonnie McKee", "static/img/americanGirl.jpg", "static/video/americanGirl.mp4", "static/video/previews/americanGirl.mp4", 2, "static/motion_list/forgetYou.json"),
    ("Burn", "Ellie Goulding", "static/img/burn.jpg", "static/video/burn.mp4", "static/video/previews/burn.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Firework", "katy Perry", "static/img/firework.jpg", "static/video/firework.mp4", "static/video/previews/firework.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Kiss Me More", "Doja Cat Ft. Sza", "static/img/kissMeMore.jpg", "static/video/kissMeMore.mp4", "static/video/previews/kissMeMore.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Dynamite", "Tajo Cruz", "static/img/dynamite.jpg", "static/video/dynamite.mp4", "static/video/previews/dynamite.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Oops I Did It Again", "Britney Spears", "static/img/oopsIDidItAgain.jpg", "static/video/previews/ooopsIDidItAgain.mp4", "static/video/previews/ooopsIDidItAgain.mp4", 1, "static/motion_list/forgetYou.json"),
    ("Montero", "Lil Nas X", "static/img/montero.jpg", "static/video/montero.mp4", "static/video/previews/montero.mp4", 1, "static/motion_list/forgetYou.json");
    ("Rasputin", "Lil Nas X", "static/img/montero.jpg", "static/video/full/rasputin.mp4", "static/video/previews/montero.mp4", 1, "static/motion_list/raRaRasputin.json");