USE dance_fusion;

<<<<<<< HEAD
CREATE TABLE achievements(
                             `achievement_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                             `achievement_name` VARCHAR(255) NOT NULL,
                             `level` INTEGER NOT NULL,
                             `max_progress` INTEGER NOT NULL
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE utf8mb4_unicode_ci;

=======
>>>>>>> origin/main_test_0208
CREATE TABLE user_achievements (
                                   `user_achievement_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                                   `user_id` INT NOT NULL,
                                   `achievement_id` INT NOT NULL,
                                   `achievement_name` VARCHAR(255) NOT NULL,
                                   `progress` INTEGER NOT NULL DEFAULT 0,
                                   `completed` TINYINT NOT NULL DEFAULT 0,
                                   `level` INTEGER NOT NULL,
                                   `max_progress` INTEGER NOT NULL,
                                   FOREIGN KEY (`user_id`) REFERENCES dance_fusion.users(`id`),
                                   FOREIGN KEY (`achievement_id`) REFERENCES dance_fusion.achievements(`achievement_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

<<<<<<< HEAD
DELIMITER //
CREATE TRIGGER add_new_achievement_to_user_achievements
    AFTER INSERT ON achievements
    FOR EACH ROW
=======
CREATE TABLE achievements(
                      `achievement_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                      `achievement_name` VARCHAR(255) NOT NULL,
                      `level` INTEGER NOT NULL,
                      `max_progress` INTEGER NOT NULL
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE utf8mb4_unicode_ci;

DELIMITER //
CREATE TRIGGER add_new_achievement_to_user_achievements
AFTER INSERT ON achievements
FOR EACH ROW
>>>>>>> origin/main_test_0208
BEGIN
    INSERT INTO user_achievements (user_id, achievement_id, achievement_name,  progress, completed, level, max_progress)
    SELECT
        u.id AS user_id,
        NEW.achievement_id,
        NEW.achievement_name,
        0 AS progress,
        0 AS completed,
        NEW.level,
        NEW.max_progress
    FROM
        users u;
END;
//
DELIMITER ;

<<<<<<< HEAD
DELIMITER //
CREATE TRIGGER add_new_user_to_user_achievements
    AFTER INSERT ON users
    FOR EACH ROW
BEGIN
    INSERT INTO user_achievements (user_id, achievement_id, achievement_name,  progress, completed, level, max_progress)
    SELECT
        NEW.id,
        a.achievement_id AS achievement_id,
        a.achievement_name AS achievement_name,
        0 AS progress,
        0 AS completed,
        a.level AS level,
        a.max_progress AS max_progress
    FROM
        achievements a;
END;
//
DELIMITER ;

INSERT INTO achievements(achievement_name, level, max_progress)
VALUES
    ("Debil", 2, 5);

SELECT * FROM user_achievements;

=======
INSERT INTO achievements(achievement_name, level, max_progress)
VALUES
    ("Debil", 2, 5);
    
SELECT * FROM user_achievements;
>>>>>>> origin/main_test_0208
