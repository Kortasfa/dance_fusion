package main

import (
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func gyrPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pages/gyr_page.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := struct{}{}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	log.Println("gyr_page.html Request completed successfully")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pages/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	err = tmpl.Execute(w, "wss://"+r.Host+"/echo")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func previewVideoWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message) // Добавить поиск в бд ссылки на превью и отправить его
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type stylesData struct {
	StyleID   int    `db:"id"`
	StyleName string `db:"name"`
}

type songsData struct {
	SongID          int    `db:"id"`
	SongName        string `db:"song_name"`
	SongAuthor      string `db:"author_name"`
	PreviewVideoSrc string `db:"preview_video_src"`
	ImageSrc        string `db:"image_src"`
	StyleID         int    `db:"style_id"`
}

type menuPageData struct {
	Styles []stylesData
	Songs  []songsData
}

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

func menuPage(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("pages/main_menu.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
		styles, err := getStylesData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		songs, err := getSongsData(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}
		data := menuPageData{
			Styles: styles,
			Songs:  songs,
		}
		//err = tmpl.Execute(w, "wss://"+r.Host+"/preview_video")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}
	}
}
