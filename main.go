package main

import (
  "github.com/nuclio/nuclio-sdk-go"
  "github.com/artonge/go-csv-tag"
)

type Demo struct {
	Name string  `csv:"name"`
	ID   int     `csv:"ID"`
	Num  float64 `csv:"number"`
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	tab := []Demo{
		Demo{
			Name: "some name",
			ID:   1,
			Num:  42.5,
		},
	}

	err := csvtag.DumpToFile(tab, "csv_file_name.csv")
	if err != nil {
		return nuclio.Response{
			StatusCode: 500,
		}, nil
	}

	return nuclio.Response{
		StatusCode: 200,
	}, nil
}
