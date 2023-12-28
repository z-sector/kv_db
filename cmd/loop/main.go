package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"

	"kv_db/internal/database"
	"kv_db/internal/database/compute"
	"kv_db/internal/database/compute/analyzer"
	"kv_db/internal/database/compute/parser"
	"kv_db/internal/database/storage"
	"kv_db/internal/database/storage/backend/memory"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	memory.NewHashTable()
	dbStorage := storage.MustStorage(
		memory.NewHashTable(),
		logger.With(slog.String("layer", "storage")),
	)

	queryAnalyzer := analyzer.MustAnalyzer(logger.With(slog.String("layer", "analyzer")))
	queryParser := parser.MustParser(logger.With(slog.String("layer", "parser")))

	queryCompute := compute.MustCompute(
		queryParser,
		queryAnalyzer,
		logger.With(slog.String("layer", "compute")),
	)

	db := database.MustDatabase(
		queryCompute,
		dbStorage,
		logger.With(slog.String("layer", "database")),
	)

	for {
		fmt.Println("input command:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		res := db.HandleQuery(context.Background(), scanner.Text())
		fmt.Printf("result: %s\n", res)
	}
}
