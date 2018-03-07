package main

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io"
	"net/http"
	"testTask/data"
	"testTask/shared"
)

// decodeIncomingData tries to cast incoming data to a data.Message object.
func decodeIncomingData(r io.Reader) data.Message {
	decoder := json.NewDecoder(r)

	var input data.Message
	err := decoder.Decode(&input)
	if err != nil {
		shared.ThrowError(shared.ErrBadRequest, "parsing json", err)
	}

	return input
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	input := decodeIncomingData(r.Body)

	Send("http://localhost:8080/query", prepareData(input.Query))
	Send("http://localhost:8080/sender", prepareData(input.Sender))
}

func prepareData(d interface{}) *bytes.Buffer {
	dataBuf := &bytes.Buffer{}
	defWriter, err := flate.NewWriter(dataBuf, flate.BestCompression)
	if err != nil {
		shared.ThrowError(shared.ErrInternalError, "trying to init deflate writers", err)
	}

	json.NewEncoder(defWriter).Encode(d)
	err = defWriter.Close()
	if err != nil {
		shared.ThrowError(shared.ErrInternalError, "trying to close deflate writers", err)
	}

	return dataBuf
}

// Send tries to send data to another microservice.
func Send(url string, data io.Reader) {
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		shared.ThrowError(shared.ErrInternalError, "sending data", err)
	}
	defer resp.Body.Close()
}

func main() {
	sParams := &shared.ServerParameters{
		AddressToListen: ":8081",
	}

	dataHandler := http.Handler(http.HandlerFunc(MessageHandler))
	adapters := []shared.Adapter{
		shared.DefineHTTPMethod("PUT"),
		shared.HandleErrors(),
		shared.LogQueries()}

	dataEndpoint := shared.Endpoint{
		Handler:  dataHandler,
		Adapters: adapters,
		Path:     "/message",
	}

	server := shared.NewServer(sParams, dataEndpoint)
	panic(server.Start())
}
