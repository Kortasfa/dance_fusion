package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
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
			style_id
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

func getMotionListPath(db *sqlx.DB, songName string) (string, error) {
	const query = `
		SELECT
			motion_list_path
		FROM
			songs
		WHERE
		   song_name=?
	`

	var motionListPath string
	err := db.QueryRow(query, songName).Scan(&motionListPath)
	if err != nil {
		return "", err
	}

	return motionListPath, nil
}

func insertNewUser(db *sqlx.DB, userName string, password string) (int, error) {
	user := struct {
		UserName     string
		Password     string
		UserImageSrc string
	}{
		UserName:     userName,
		Password:     password,
		UserImageSrc: "static/img/user_1.png",
	}
	query := `
		INSERT INTO users(name, password, img_src)
		VALUES (?, ?, ?)`
	result, err := db.Exec(query, user.UserName, user.Password, user.UserImageSrc)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userID), nil
}

func userExists(db *sqlx.DB, userName string) (bool, error) {
	const query = `
			SELECT COUNT(*)
			FROM users
			WHERE name = ?`
	var count int
	err := db.QueryRow(query, userName).Scan(&count)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func credentialExists(db *sqlx.DB, userName string, password string) (int, bool, error) {
	const query = `
		SELECT id
		FROM users
		WHERE name = ? and password = ?`
	var userIDs []int
	err := db.Select(&userIDs, query, userName, password)
	if len(userIDs) == 0 {
		return 0, false, nil
	} else if err != nil {
		log.Println(err.Error())
		return 0, false, err
	}
	return userIDs[0], true, nil
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
		WHERE id = ?
	`

	_, err := db.Exec(query, field.HatSrc, field.FaceSrc, field.BodySrc, userID)
	return err
}

func getScoreByUserID(db *sqlx.DB, userID int) (int, error) {
	const query = `
		SELECT
			score
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
