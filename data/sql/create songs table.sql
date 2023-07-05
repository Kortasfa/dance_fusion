USE dance_fusion;
CREATE TABLE  songs(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY, 
    `song_name` VARCHAR(255) NOT NULL,
    `author_name` VARCHAR(255) NOT NULL,
    `image_src` VARCHAR(255) NOT NULL,
    `audio_src` VARCHAR(255) NOT NULL,
    `video_src` VARCHAR(255) NOT NULL,
    `preview_video_src` VARCHAR(255) NOT NULL,
    `style_id` INTEGER NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

