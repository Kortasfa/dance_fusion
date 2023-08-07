USE dance_fusion;
CREATE TABLE faces(
    `id` INTEGER AUTO_INCREMENT PRIMARY KEY,
    `recommended_level` INTEGER NOT NULL,
    `face_src` VARCHAR(255) NOT NULL
) ENGINE = InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_unicode_ci;

INSERT INTO faces (id, recommended_level, face_src) VALUES
(1, 0, 'static/img/face/2.png'),
(2, 0, 'static/img/face/6.png'),
(3, 1, 'static/img/face/5.png'),
(4, 2, 'static/img/face/4.png'),
(5, 3, 'static/img/face/3.png'),
(6, 4, 'static/img/face/7.png'),
(7, 4, 'static/img/face/1.png'),
(8, 5, 'static/img/face/smileFace.png'),
(9, 5, 'static/img/face/simpleFace.png'),
(10, 6, 'static/img/face/frightenedFace.png'),
(11, 6, 'static/img/face/angryFace.png'),
(12, 7, 'static/img/face/puzzledFace.png'),
(13, 8, 'static/img/face/whatFace.png'),
(14, 9, 'static/img/face/deadInsideFace.png'),
(15, 10, 'static/img/face/pickUpFace.png'),
(16, 11, 'static/img/face/scaredFace.png'),
(17, 12, 'static/img/face/angriestFace.png'),
(18, 13, 'static/img/face/stupidFace.png'),
(19, 14, 'static/img/face/unhappyFace.png'),
(20, 14, 'static/img/face/wFace.png'),
(21, 14, 'static/img/face/whatFace.png'),
(22, 15, 'static/img/face/yourFace.png'),
(23, 15, 'static/img/face/loveFace.png'),
(24, 17, 'static/img/face/funFace.png'),
(25, 17, 'static/img/face/badFace.png'),
(26, 17, 'static/img/face/fullFace.png');
