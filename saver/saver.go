package main

import (
	"compress/flate"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"testTask/data"
	"testTask/shared"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SqliteDBPath  = "queriesLog.sqlite"
	FilesaverPath = "log"
)

type Env struct {
	senderSaver *FileSaver
	querySaver  *SqliteSaver
}

var (
	env *Env
)

func applyDeflate(reader io.Reader) io.Reader {
	return flate.NewReader(reader)
}

// JsonToQuery reads data from the reader and tries to deserialize the object.
func JsonToQuery(reader io.Reader) (data.Query, error) {
	var obj data.Query
	err := json.NewDecoder(reader).Decode(&obj)
	return obj, err
}

func SenderHandler(w http.ResponseWriter, r *http.Request) {
	reader := applyDeflate(r.Body)

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		shared.ThrowError(shared.ErrInternalError, "trying to read sender bytes", err)
	}

	err = env.senderSaver.SaveData(data)
	if err != nil {
		shared.ThrowError(shared.ErrInternalError, "trying to write sender data", err)
	}
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	query, err := JsonToQuery(applyDeflate(r.Body))
	if err != nil {
		shared.ThrowError(shared.ErrBadRequest, "trying to decode json", err)
	}

	err = env.querySaver.SaveData(query)
	if err != nil {
		shared.ThrowError(shared.ErrInternalError, "trying to save data", err)
	}
}

func main() {
	dSaver, err := NewSqliteSaver(SqliteDBPath)
	if err != nil {
		panic(err)
	}

	err = dSaver.createTable()
	if err != nil {
		panic(err)
	}

	env = &Env{
		querySaver:  dSaver,
		senderSaver: NewFileSaver(FilesaverPath),
	}

	defer env.querySaver.Close()
	defer env.senderSaver.Close()

	senderHandler := http.Handler(http.HandlerFunc(SenderHandler))
	queryHandler := http.Handler(http.HandlerFunc(QueryHandler))
	adapters := []shared.Adapter{
		shared.DefineHTTPMethod("PUT"),
		shared.HandleErrors(),
		shared.LogQueries()}

	userEndpoint := shared.Endpoint{
		Handler:  queryHandler,
		Adapters: adapters,
		Path:     "/query",
	}

	errorEndpoint := shared.Endpoint{
		Handler:  senderHandler,
		Adapters: adapters,
		Path:     "/sender",
	}

	server := shared.NewServer(nil, errorEndpoint, userEndpoint)
	panic(server.Start())
}
