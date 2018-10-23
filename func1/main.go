package main

type Demo struct {
	Name string  `csv:"name"`
	ID   int     `csv:"ID"`
	Num  float64 `csv:"number"`
}

func Handlee(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
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
		StatusCode: 201,
	}, nil
}
