package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func getStylesData(db *sqlx.DB) ([]stylesData, error) {
	const query = `
		SELECT
			id,
			name
		FROM
			styles
	`
	var data []stylesData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getSongsData(db *sqlx.DB) ([]songsData, error) {
	const query = `
		SELECT
			id,
			song_name,
			author_name,
			video_src,
			preview_video_src,
			image_src,
			style_id,
			difficulty
		FROM
			songs
	`
	var data []songsData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getMotionListPath(db *sqlx.DB, songName string) ([]string, error) {
	const query = `
		SELECT
			first_player,
			second_player,
			third_player,
			fourth_player
		FROM
			motion_list_path
		WHERE
		   song_name = ?
	`

	var playerOnePath, playerTwoPath, playerThreePath, playerFourPath string
	err := db.QueryRow(query, songName).Scan(&playerOnePath, &playerTwoPath, &playerThreePath, &playerFourPath)
	if err != nil {
		return nil, err
	}

	// Создаем срез и добавляем пути к файлам для каждого игрока
	motionListPath := []string{playerOnePath, playerTwoPath, playerThreePath, playerFourPath}

	return motionListPath, nil
}

func insertNewUser(db *sqlx.DB, userName string, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	query := `
		INSERT INTO
		    users(name, password_hash)
		VALUES
		    (?, ?)`
	result, err := db.Exec(query, userName, string(hash))
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userID), nil
}

func userExists(db *sqlx.DB, userName string) (int, bool, error) {
	const query = `
		SELECT
		    id
		FROM
		    users
		WHERE
		    name = ?`
	var userIDs []int
	err := db.Select(&userIDs, query, userName)
	if len(userIDs) == 0 {
		return 0, false, nil
	} else if err != nil {
		log.Println(err.Error())
		return 0, false, err
	}
	return userIDs[0], true, nil
}

func credentialExists(db *sqlx.DB, userName string, password string) (int, bool, error) {
	const query = `
		SELECT
			id, password_hash
		FROM
			users
		WHERE
			name = ?`

	rows, err := db.Query(query, userName)
	if err != nil {
		log.Println("Failed to execute the query")
		return 0, false, err
	}

	var (
		userID       int
		passwordHash string
	)

	if rows.Next() {
		err = rows.Scan(&userID, &passwordHash)
		if err != nil {
			log.Println("Failed to retrieve data from the row")
			return 0, false, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		if err == nil {
			return userID, true, nil
		}
	}
	return 0, false, nil
}

func getUserInfo(db *sqlx.DB, userID int) (userInfo, error) {
	const query = `
		SELECT
			id,
			name,
			img_hat,
			img_face,
			img_body
		FROM
			users
		WHERE
			id = ?
	`

	var user struct {
		UserID   int    `db:"id"`
		UserName string `db:"name"`
		HatSrc   string `db:"img_hat"`
		FaceSrc  string `db:"img_face"`
		BodySrc  string `db:"img_body"`
	}

	err := db.Get(&user, query, userID)
	if err != nil {
		return userInfo{}, err
	}

	return userInfo{
		UserID:   user.UserID,
		UserName: user.UserName,
		HatSrc:   user.HatSrc,
		FaceSrc:  user.FaceSrc,
		BodySrc:  user.BodySrc,
	}, nil
}

func getConnectedUsers(roomID string, db *sqlx.DB) ([]userInfo, error) {
	userIDs := roomIDDict[roomID]
	var users []userInfo

	for _, userID := range userIDs {
		const query = `
			SELECT
				id,
				name,
				img_hat,
				img_face,
				img_body
			FROM
				users
			WHERE
				id = ?`

		row := db.QueryRow(query, userID)
		var user userInfo

		err := row.Scan(&user.UserID, &user.UserName, &user.HatSrc, &user.FaceSrc, &user.BodySrc)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func getHatData(db *sqlx.DB) ([]hatData, error) {
	const query = `
		SELECT
			id,
			recommended_level,
			hat_src
		FROM
			hats
	`
	var data []hatData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getFaceData(db *sqlx.DB) ([]facesData, error) {
	const query = `
		SELECT
			id,
			recommended_level,
			face_src
		FROM
			faces
	`
	var data []facesData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getBodyData(db *sqlx.DB) ([]bodyData, error) {
	const query = `
		SELECT
			id,
			recommended_level,
			body_src
		FROM
			bodies
	`
	var data []bodyData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func changeUserAvatar(db *sqlx.DB, field userAvatarData, userID int) error {
	const query = `
		UPDATE
			users
		SET
			img_hat = ?,
			img_face = ?,
			img_body = ?
		WHERE
		    id = ?
	`

	_, err := db.Exec(query, field.HatSrc, field.FaceSrc, field.BodySrc, userID)
	return err
}

func getScoreByUserID(db *sqlx.DB, userID int) (int, error) {
	const query = `
		SELECT
			total_score
		FROM
			users
		WHERE
			id = ?
	`

	var score int
	err := db.QueryRow(query, userID).Scan(&score)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("пользователь с ID %d не найден", userID)
		}
		return 0, err
	}

	return score, nil
}

func getBestPlayerInfo(db *sqlx.DB, songID int) (bestPlayerInfo, error) {
	const query = `
		SELECT
		    best_player_id,
		    best_score
		FROM
		    songs
		WHERE
		    id = ?`
	var playerInfo struct {
		UserID sql.NullInt32
		Score  sql.NullInt32
	}
	err := db.QueryRow(query, songID).Scan(&playerInfo.UserID, &playerInfo.Score)
	if err != nil {
		return bestPlayerInfo{}, err
	}

	return bestPlayerInfo{
		UserID: int(playerInfo.UserID.Int32),
		Score:  int(playerInfo.Score.Int32),
	}, nil
}

func getBotInfo(db *sqlx.DB, botName string) (botInfo, error) {
	const query = `
		SELECT
		    bot_id,
		    bot_scores_path,
		    img_hat,
		    img_body,
		    img_face,
		    difficulty
		FROM
		    bots
		WHERE
		    bot_name = ?`
	var botMainInfo botInfo
	err := db.QueryRow(query, botName).Scan(&botMainInfo.BotId, &botMainInfo.BotScoresPath, &botMainInfo.BotImgHat, &botMainInfo.BotImgBody, &botMainInfo.BotImgFace, &botMainInfo.Difficulty)
	if err != nil {
		return botInfo{}, err
	}
	return botMainInfo, nil
}

func updateBestPlayerSQL(db *sqlx.DB, songID int, userID int, score int) error {
	const query = `
			UPDATE
				songs
			SET
				best_player_id = ?, best_score = ?
			WHERE
			    id = ?
		`

	_, err := db.Exec(query, userID, score, songID)
	return err
}

func updateUserName(db *sqlx.DB, userID int, userName string) error {
	const query = `
			UPDATE
				users
			SET
				name = ?
			WHERE
			    id = ?
		`

	_, err := db.Exec(query, userName, userID)
	return err
}

func updateUserPassword(db *sqlx.DB, userID int, userPassword string) error {
	const query = `
			UPDATE
				users
			SET
				password_hash = ?
			WHERE
			    id = ?
		`
	hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), hashCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(query, string(hash), userID)
	return err
}

func getBotNames(db *sqlx.DB) ([]botNameData, error) {
	const query = `
		SELECT
			bot_name,
			difficulty
		FROM
			bots
	`
	var data []botNameData

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getBossInfo(db *sqlx.DB) ([]bossInfo, error) {
	const query = `
		SELECT
		    boss_id,
		    boss_name,
		    boss_health_point,
		    img_hat,
		    img_body,
		    img_face
		FROM
		    bosses
		`
	var data []bossInfo

	err := db.Select(&data, query)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func addUserScoreSQL(db *sqlx.DB, userID int, score int) error {
	totalScore, err := getScoreByUserID(db, userID)
	if err != nil {
		return err
	}
	const query = `
			UPDATE
				users
			SET
				total_score = ?
			WHERE
			    id = ?
		`

	_, err = db.Exec(query, totalScore+score, userID)
	return err
}

func getUserAchievements(db *sqlx.DB, userID int) ([]userAchievement, error) {
	return []userAchievement{}, nil
}
