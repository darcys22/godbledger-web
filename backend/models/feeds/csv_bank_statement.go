package feeds

import (
  "encoding/csv"
  "encoding/base64"
  "strings"
  "io"
  "time"
  "strconv"
  "errors"
  "fmt"
  "crypto/sha1"

	"github.com/sirupsen/logrus"
  "github.com/araddon/dateparse"
)

type CSVBankStatementRow struct {
  Date time.Time
  Hash string
  Amount float64
  Description string
}

var log = logrus.WithField("prefix", "godbledger-feeds")

func ReadCSVBankStatement(base_64_encoded_file string, columns []string, start_row, end_row int) (rows []CSVBankStatementRow, err error) {
  log.Debug("Reading CSV Bank Statement")
  dateCol := -1
  amountCol := -1
  debitCol := -1
  creditCol := -1
  var descriptionCols []int
  for i, col := range columns {
    switch col {
      case "date":
      dateCol = i
      case "amount":
      amountCol = i
      case "debit":
      debitCol = i
      case "credit":
      creditCol = i
      case "description":
      descriptionCols = append(descriptionCols, i)
      default:
    }
  }

  if dateCol == -1 {
    err = errors.New("no date column provided")
    return
  }
  if amountCol == -1 && debitCol == -1 && creditCol == -1 {
    err = errors.New("no amount column provided")
		return 
  }
  if amountCol != -1 && (debitCol != -1 || creditCol != -1){
    err = errors.New("amount column provided in addition to a debit/credit column")
    return
  }
  if amountCol == -1 && (debitCol == -1 && creditCol != -1 ){
		err = errors.New("credt column provided in debit column not provided")
    return
  }
  if amountCol == -1 && (debitCol != -1 && creditCol == -1 ){
		err = errors.New("credt column provided in debit column not provided")
    return
  }
  // Split the string to remove the prefix data:text/csv;base64 
  v := strings.Split(base_64_encoded_file, ",")
  data, err := base64.StdEncoding.DecodeString(v[1])
	if err != nil {
		return 
	}
  reader := csv.NewReader(strings.NewReader(string(data)))
  // Dont check fields, some banks drop commas at end of rows 
  reader.FieldsPerRecord = -1
  rowNumber := 0
  // Skip the first rows provided by the user
  for range [5]int{} {
    reader.Read()
    rowNumber += 1
  }

  // Iterate through the records and put into a CSVBankStatementRow
  for {
    // Read each record from csv
    record, err2 := reader.Read()
    if err2 == io.EOF {
      break
    }
    if err2 != nil {
      err = errors.New("Could not read row")
      return
    }
    var row CSVBankStatementRow 
    h := sha1.New()
    h.Write([]byte(strings.Join(record, "")))
    row.Hash = base64.URLEncoding.EncodeToString(h.Sum(nil))
    row.Amount = 0
    row.Date, err = dateparse.ParseLocal(record[dateCol])
		if err != nil {
      err = errors.New(fmt.Sprintf("Could not parse date \"%s\" on row %d", record[dateCol], rowNumber))
      return
		}
    if amountCol != -1 {
      i, err2 := strconv.ParseFloat(record[amountCol], 64)
      if err2 != nil {
          err = errors.New(fmt.Sprintf("Could not convert amount \"%s\" to a number on row %d", record[amountCol], rowNumber))
          return
      }
      row.Amount = i
    }
    if debitCol != -1 {
      i, err2 := strconv.ParseFloat(record[debitCol], 64)
      if err2 != nil {
          err = errors.New(fmt.Sprintf("Could not convert debit amount \"%s\" to a number on row %d", record[debitCol], rowNumber))
          return
      }
      row.Amount += i
    }
    if creditCol != -1 {
      i, err2 := strconv.ParseFloat(record[creditCol], 64)
      if err2 != nil {
          err = errors.New(fmt.Sprintf("Could not convert credit amount \"%s\" to a number on row %d", record[creditCol], rowNumber))
          return
      }
      row.Amount -= i
    }

    for i, col := range descriptionCols {
      if i != 0 {
        row.Description += " "
      }
      row.Description += record[col]
    }
    rowNumber += 1
    rows = append(rows, row)
  }

	return
}
