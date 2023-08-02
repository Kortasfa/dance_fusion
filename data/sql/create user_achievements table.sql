USE dance_fusion;

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

INSERT INTO achievements(achievement_name, level, max_progress)
VALUES
    ("Debil", 2, 5);
    
SELECT * FROM user_achievements;