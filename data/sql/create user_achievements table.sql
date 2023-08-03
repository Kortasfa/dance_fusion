USE dance_fusion;
drop table user_achievements;
drop table achievements;

CREATE TABLE achievements(
                             `achievement_id` INTEGER AUTO_INCREMENT PRIMARY KEY,
                             `achievement_name` VARCHAR(255) NOT NULL,
                             `max_progress` INTEGER NOT NULL
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
    INSERT INTO user_achievements (user_id, achievement_id, achievement_name,  progress, completed, collected, max_progress)
    SELECT
        u.id AS user_id,
        NEW.achievement_id,
        NEW.achievement_name,
        0 AS progress,
        0 AS completed,
        0 AS collected,
        NEW.max_progress
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
    INSERT INTO user_achievements (user_id, achievement_id, achievement_name,  progress, completed, collected, max_progress)
    SELECT
        NEW.id,
        a.achievement_id AS achievement_id,
        a.achievement_name AS achievement_name,
        0 AS progress,
        0 AS completed,
        0 AS collected,
        a.max_progress AS max_progress
    FROM
        achievements a;
END;
//
DELIMITER ;

INSERT INTO achievements(achievement_name, max_progress)
VALUES
    ("Debil", 5);


SELECT * FROM user_achievements;

