package main

import (
	"fmt"
	"os"

	"github.com/etheryen/chessarbiter-scraper/scrape"
	"github.com/etheryen/chessarbiter-scraper/utils"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	tournaments, err := scrape.GetByYearMonth(2024, 6)
	if err != nil {
		return err
	}

	utils.PrintStructSlice(tournaments)

	fmt.Println("\nFound:", len(tournaments))

	return nil
}
