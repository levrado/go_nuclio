package csvtag

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

// DumpToFile - writes a slice content into a file specified by filePath
// @param slice: An object typically of the form []struct, where the struct using csv tag
// @param filePath: The file path string of where you want the file to be created
// @return an error if one occures
func DumpToFile(slice interface{}, filePath string) error {
	// Create file object
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// Dump the slice into the file
	return Dump(slice, file)
}

// Dump - writes a slice content into an io.Writer
// @param slice: an object typically of the form []struct, where the struct using csv tags
// @param writer: the location of where you will write the slice content to. Example: File, Stdout, etc
// @return an error if one occures
func Dump(slice interface{}, writer io.Writer) error {
	// Determines the type of the elements of the passed slice
	reflectedValue := reflect.ValueOf(slice)

	// If slice is a pointer, get the value it points to
	// (if it isn't, Indirect() does nothing and returns the value it was called with)
	reflectedValue = reflect.Indirect(reflectedValue)

	// Return when slice is not a slice
	if reflectedValue.Kind() != reflect.Array && reflectedValue.Kind() != reflect.Slice {
		return errors.New("Unsupported data type")
	}
	// Get the headers
	var headers []string
	for i := 0; i < reflectedValue.Type().Elem().NumField(); i++ {
		name := reflectedValue.Type().Elem().Field(i).Tag.Get("csv")
		if name != "" {
			headers = append(headers, name)
		}
	}
	// Create a csv.Writer to write the content in the csv format
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush() // Ensures that all data in the buffer has been written to io.Writer
	err := csvWriter.Write(headers)
	if err != nil {
		return err
	}
	// Extract the content of the slice
	for i := 0; i < reflectedValue.Len(); i++ {
		itemData := []string{}
		fields := reflectedValue.Index(i).NumField()
		for j := 0; j < fields; j++ {
			if reflectedValue.Index(i).Type().Field(j).Tag.Get("csv") != "" {
				switch reflectedValue.Index(i).Type().Field(j).Type.Kind() {
				case reflect.Float64, reflect.Float32:
					itemData = append(itemData, strconv.FormatFloat(reflectedValue.Index(i).Field(j).Float(), 'f', -1, 64))
				default:
					itemData = append(itemData, fmt.Sprint(reflectedValue.Index(i).Field(j)))
				}
			}
		}
		// Write the line into io.Writer
		err = csvWriter.Write(itemData)
		if err != nil {
			return err
		}
	}
	return nil
}
