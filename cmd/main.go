package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/masnyjimmy/gofig"
)

type Database struct {
	User     string        `conf:"user"`
	Id       int           `conf:"id,25"`
	Bool     bool          `conf:"bool,false"`
	Duration time.Duration `conf:"dur"`
	Date     time.Time     `conf:"dat"`
}

type Config struct {
	Database Database `conf:"database"`
}

func main() {
	godotenv.Load()

	fields := gofig.GenerateFields[Config]()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")

	enc.Encode(fields)

	g := gofig.New(fields)
	g.AddTimeFormats(time.TimeOnly)

	g.Read(&gofig.EnvSource{})
	var config Config

	if err := g.Unmarshall(&config); err != nil {
		panic(err)
	}
	enc.Encode(config)
}
