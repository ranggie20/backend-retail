package util

import (
	"crypto/rand"
	"database/sql"
	"io"
	"log"
	"strconv"
	"time"
)

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func ConvertStringIDToInt32(ID string) int32 {
	IDInt, err := strconv.Atoi(ID)
	if err != nil {
		// Penanganan kesalahan jika konversi gagal
		log.Fatal(err)
	}
	// Convert IDInt to int32
	var IDInt32 int32 = int32(IDInt)
	return IDInt32
}

func ConvertStringToDate(dateString string) time.Time {

	layout := "2006-01-02" // The layout matches the format of the input string

	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		// Penanganan kesalahan jika konversi gagal
		log.Fatal("Error parsing date:", err)
	}

	return parsedTime
}

func SqlFloat32(f float32) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: float64(f), Valid: true}
}
