package main

import (
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

func getUserInfo(db *sqlx.DB, userID string) (string, string, error) {
	const query = `
		SELECT
			name,
			img_src
		FROM
			users
		WHERE
		   id=?
	`
	row := db.QueryRow(query, userID)
	//data := new(userData)
	data := new(struct {
		UserID   int    `db:"id"`
		UserName string `db:"name"`
		ImgSrc   string `db:"img_src"`
	})
	err := row.Scan(&data.UserName, &data.ImgSrc)
	if err != nil {
		return "", "", err
	}
	return data.UserName, data.ImgSrc, nil
}
