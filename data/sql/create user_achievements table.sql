USE dance_fusion;
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS achievements;
DROP TRIGGER IF EXISTS add_new_achievement_to_user_achievements;
DROP TRIGGER IF EXISTS add_new_user_to_user_achievements;

CREATE TABLE achievements(
                             `achievement_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                             `achievement_name` VARCHAR(255) NOT NULL,
                             `max_progress` INTEGER NOT NULL,
                             `song_id` INTEGER NOT NULL,
                             `bot_difficulty` INTEGER NOT NULL,
                             `boss_id` INTEGER NOT NULL DEFAULT 0,
                             `score` INTEGER NOT NULL
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE TABLE user_achievements (
                                   `user_achievement_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                                   `user_id` INT NOT NULL,
                                   `achievement_id` INT NOT NULL,
                                   `achievement_name` VARCHAR(255) NOT NULL,
                                   `progress` INTEGER NOT NULL DEFAULT 0,
                                   `completed` TINYINT NOT NULL DEFAULT 0,
                                   `collected` TINYINT NOT NULL DEFAULT 0,
                                   `max_progress` INTEGER NOT NULL,
                                   `song_id` INTEGER NOT NULL,
                                   `bot_difficulty` INTEGER NOT NULL,
                                   `boss_id` INTEGER NOT NULL DEFAULT 0,
                                   `score` INTEGER NOT NULL,
                                   FOREIGN KEY (`user_id`) REFERENCES dance_fusion.users(`id`),
                                   FOREIGN KEY (`achievement_id`) REFERENCES dance_fusion.achievements(`achievement_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

DELIMITER //
CREATE TRIGGER add_new_achievement_to_user_achievements
AFTER INSERT ON achievements
FOR EACH ROW
BEGIN
    INSERT INTO user_achievements (user_id, achievement_id, achievement_name,  progress, completed, collected, max_progress, song_id, bot_difficulty, boss_id, score)
    SELECT
        u.id AS user_id,
        NEW.achievement_id,
        NEW.achievement_name,
        0 AS progress,
        0 AS completed,
        0 AS collected,
        NEW.max_progress,
        NEW.song_id,
        NEW.bot_difficulty,
        NEW.boss_id,
        NEW.score
    FROM
        users u;
END;
//
DELIMITER ;

DELIMITER //
CREATE TRIGGER add_new_user_to_user_achievements
    AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO user_achievements (user_id, achievement_id, achievement_name,  progress, completed, collected, max_progress, song_id, bot_difficulty, boss_id, score)
    SELECT
        NEW.id,
        a.achievement_id AS achievement_id,
        a.achievement_name AS achievement_name,
        0 AS progress,
        0 AS completed,
        0 AS collected,
        a.max_progress AS max_progress,
        a.song_id AS song_id,
        a.bot_difficulty AS bot_difficulty,
        a.boss_id AS boss_id,
        a.score AS score
    FROM
        achievements a;
END;
//
DELIMITER ;

INSERT INTO achievements(achievement_name, max_progress, song_id, bot_difficulty, boss_id, score)
VALUES
	('Dance "Forget You" 1 time', 1, 1, 0, 0, 2700),
	('Dance "Forget You" 3 times', 3, 1, 0, 0, 4000),
	('Dance "Forget You" 5 times', 5, 1, 0, 0, 5400),
	('Win the easy bot', 1, 0, 1, 0, 2000),
	('Win the medium bot', 1, 2, 0, 0, 3000),
    ('Win the hard bot', 1, 0, 3, 0, 4000),
    ('Win the extreme bot', 1, 0, 4, 0, 5000),
    ('Finish off Rhythmo', 1, 0, 0, 4, 2000),
    ('Finish off Electric Maestro', 1, 0, 0, 3, 3000),
    ('Finish off Funk Fusionist', 1, 0, 0, 2, 4000),
    ('Finish off Tempo Titan', 1, 0, 0, 1, 5000);


SELECT * FROM user_achievements;

