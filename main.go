package main

import (
	"context"
	"fmt"

	"github.com/DIMO-Network/shared"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/warp-contracts/syncer/src/utils/arweave"
	"github.com/warp-contracts/syncer/src/utils/bundlr"
	"github.com/warp-contracts/syncer/src/utils/config"
)

func main() {
	ctx := context.Background()
	settings, err := shared.LoadConfig[struct {
		PrivateKey string `yaml:"PRIVATE_KEY"`
	}]("settings.yaml")
	if err != nil {
		panic(err)
	}

	signer, err := bundlr.NewEthereumSigner("0x" + settings.PrivateKey)
	fmt.Println(crypto.PubkeyToAddress(signer.PrivateKey.PublicKey))

	client := bundlr.NewClient(ctx, &config.Bundlr{
		Urls: []string{
			"https://devnet.bundlr.network",
		},
		Wallet: crypto.PubkeyToAddress(signer.PrivateKey.PublicKey).String(),
	})

	dataItem := bundlr.BundleItem{
		Data: arweave.Base64String([]byte("test")),
	}

	err = dataItem.Sign(signer)

	resp1, resp2, err := client.Upload(ctx, &dataItem)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp1, resp2)
}
