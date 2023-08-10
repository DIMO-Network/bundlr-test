package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/DIMO-Network/shared"
	"github.com/warp-contracts/syncer/src/utils/arweave"
	"github.com/warp-contracts/syncer/src/utils/bundlr"
)

func main() {
	settings, err := shared.LoadConfig[struct {
		PrivateKey string `yaml:"PRIVATE_KEY"`
	}]("settings.yaml")
	if err != nil {
		panic(err)
	}

	signer, err := bundlr.NewEthereumSigner("0x" + settings.PrivateKey)
	dataItem := bundlr.BundleItem{
		Data: arweave.Base64String([]byte("abc")),
		Tags: bundlr.Tags{
			bundlr.Tag{Name: "Content-Type", Value: "text"},
		},
	}

	err = dataItem.Sign(signer)
	if err != nil {
		panic(err)
	}

	reader, err := dataItem.Reader()
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	responseBody := bytes.NewBuffer(body)

	resp, err := http.Post("https://devnet.bundlr.network/tx/matic", "application/octet-stream", responseBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	sb := string(body)
	log.Printf(sb)
}
