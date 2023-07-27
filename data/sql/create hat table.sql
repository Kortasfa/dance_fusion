USE dance_fusion;
CREATE TABLE hats(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
    `recommended_level` INTEGER NOT NULL,
    `hat_src` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO hats(id, recommended_level, hat_src)
VALUES
	("1","0","../static/img/hat/withoutHat.png"),
    ("2","0","../static/img/hat/stanHat.png"),
    ("3","0","../static/img/hat/bakerHat.png"),
    ("4","0","../static/img/hat/robinHoodHat.png"),
    ("5","1","../static/img/hat/cowboyHat.png"),
    ("6","1","../static/img/hat/flowerHat.png"),
    ("7","1","../static/img/hat/robinHat.png"),
    ("8","2","../static/img/hat/summerHat.png"),
    ("9","2","../static/img/hat/sherlokHat.png"),
    ("10","2","../static/img/hat/vatsonHat.png"),
    ("11","3","../static/img/hat/herHat.png"),
    ("12","3","../static/img/hat/beretHat.png"),
    ("13","5","../static/img/hat/bigHat.png");
