package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"ktsProject/controllers"
	"ktsProject/routes"
	"log"
	"net/http"
	"os"
	"strings"
)

type DbConfig struct {
	Host    string `toml:"host"`
	Port    string `toml:"port"`
	Sslmode string `toml:"sslmode"`
	Dbname  string `toml:"dbname"`
	User    string `toml:"user"`
	Pass    string `toml:"pass"`
}

func (db DbConfig) String() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		db.Host, db.Port, db.Dbname, db.Sslmode, db.User, db.Pass,
	)
}

func dbSettings() string {
	conf := &DbConfig{}
	_, err := toml.DecodeFile("./config/DBsettings.toml", conf)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s", conf.String())
	return conf.String()
}

func main() {
	handle := controllers.Handle{}
	handle.InitDB("postgres", dbSettings())
	rawFile, err := os.Open("./init.sql")
	file, err := ioutil.ReadAll(rawFile)
	if err != nil {
		log.Fatalln(err)
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err = handle.Db.Exec(request)
	}
	if err != nil {
		log.Fatalln(err)
	}
	r := routes.Router(&handle)
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalln(err)
	}
}
