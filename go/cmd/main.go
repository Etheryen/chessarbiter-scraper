package main

import (
	"fmt"
	"os"

	"github.com/etheryen/chessarbiter-scraper/scrape"
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

	fmt.Println(tournaments)

	return nil
}
