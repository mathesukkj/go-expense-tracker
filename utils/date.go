package datefmt

import "time"

func FormatDateToPostgres(date time.Time) string {
	postgresTimestampFormat := "2006-01-02 15:04:05-07:00"
	return date.Format(postgresTimestampFormat)
}
