USE dance_fusion;

CREATE TABLE songs(
                      `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                      `song_name` VARCHAR(255) NOT NULL,
                      `author_name` VARCHAR(255) NOT NULL,
                      `image_src` VARCHAR(255) NOT NULL,
                      `video_src` VARCHAR(255) NOT NULL,
                      `preview_video_src` VARCHAR(255) NOT NULL,
                      `style_id` INTEGER NOT NULL,
                      `difficulty` INTEGER NOT NULL,
                      `best_player_id` INTEGER,
                      `best_score` INTEGER
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE utf8mb4_unicode_ci;

INSERT INTO songs(song_name, author_name, image_src, video_src, preview_video_src, style_id, difficulty)
VALUES
    ("Forget You", "Cee Lo Green", "static/img/forgetYou.jpg", "static/video/full/forgetYouFull.mp4", "static/video/previews/forgetYou.mp4", 1, 2),
    ("Rasputin"," 'Boney M'", "static/img/rasputin.jpg", "static/video/full/rasputin.mp4", "static/video/previews/rasputin.mp4", 3, 2),
    ("Mi Mi Mi"," 'Serebro'", "static/img/mimimi.jpg", "static/video/full/mimimi.mp4", "static/video/previews/mimimi.mp4", 2, 2);

