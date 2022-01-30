package models

import (
	"context"
	"fmt"
	"time"

	"github.com/darcys22/godbledger-web/backend/models/backend"
	"github.com/darcys22/godbledger-web/backend/setting"

	pb "github.com/darcys22/godbledger/proto/transaction"

	"google.golang.org/grpc"
)

type Account struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type PostAccountCommand struct {
	Name string     `json:"name" binding:"required"`
	Tags []string `json:"tags"`
}

type GetAccounts struct {
	Results []Account `json:"results"`
}

func NewAccountsListing() *GetAccounts {
	accounts := []Account{}
	return &GetAccounts{accounts}
}

func (a *GetAccounts) SearchAccounts() error {
  db := backend.GetConnection()

	queryDB := `
	SELECT 
		accounts.name 
	FROM accounts 
	;`

	log.Debug("Querying Database")
	rows, err := db.Query(queryDB)
	if err != nil {
		log.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	a.Results = []Account{}

	index := 0

	for rows.Next() {
		//Scan one account record
		index++
		t := Account{ID: index}
		if err := rows.Scan(&t.Text); err != nil {
			log.Errorf("Could not scan rows of query (%v)", err)
		}
		a.Results = append(a.Results, t)
	}
	if rows.Err() != nil {
		log.Errorf("rows errored with (%v)", rows.Err())
	}
	return nil
}

func (j *PostAccountCommand) Save() error {
	log.Trace("Calling Save Account function")
  cfg := setting.GetConfig()
	address := fmt.Sprintf("%s:%s", cfg.GoDBLedgerHost, cfg.GoDBLedgerPort)
	log.WithField("address", address).Info("GRPC Dialing on port")
	opts := []grpc.DialOption{}

	if cfg.GoDBLedgerCACert != "" && cfg.GoDBLedgerCert != "" && cfg.GoDBLedgerKey != "" {
		tlsCredentials, err := loadTLSCredentials(cfg)
		if err != nil {
			return fmt.Errorf("Could not load TLS credentials (%v)", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return fmt.Errorf("Could not connect to GRPC (%v)", err)
	}
	defer conn.Close()
	client := pb.NewTransactorClient(conn)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.AccountTagRequest{
		Account:  j.Name,
		Tag:      j.Tags,
	}
	r, err := client.AddAccount(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Add Account Method (%v)", err)
	}
	log.Infof("Add Account Response: %s", r.GetMessage())
	return nil
}

func DeleteAccountCommand(id string) error {
	log.Trace("Calling Delete Account function")

  cfg := setting.GetConfig()
	address := fmt.Sprintf("%s:%s", cfg.GoDBLedgerHost, cfg.GoDBLedgerPort)
	log.WithField("address", address).Info("GRPC Dialing on port")
	opts := []grpc.DialOption{}

	if cfg.GoDBLedgerCACert != "" && cfg.GoDBLedgerCert != "" && cfg.GoDBLedgerKey != "" {
		tlsCredentials, err := loadTLSCredentials(cfg)
		if err != nil {
			return fmt.Errorf("Could not load TLS credentials (%v)", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return fmt.Errorf("Could not connect to GRPC (%v)", err)
	}
	defer conn.Close()
	client := pb.NewTransactorClient(conn)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.DeleteAccountTagRequest{
		Account: id,
	}
	r, err := client.DeleteAccount(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Delete Transaction Method (%v)", err)
	}
	log.Infof("Delete Transaction Response: %s", r.GetMessage())

	return nil
}

func GetAccountCommand(id string) (PostAccountCommand, error) {
	log.Trace("Calling Get Journal Command function")
	j := PostAccountCommand{}
	j.Tags = []string{}

  db := backend.GetConnection()

	queryAccountDetail := `
	SELECT 
		accounts.name 
	FROM accounts 
  WHERE name = ?
  LIMIT 1
	;`

	log.Debug("Querying Database")
	rows, err := db.Query(queryDB)
	if err != nil {
		log.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

  for rows.Next() {
    var t LineItem
    var decimals float64
    var narration string
    if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Currency, &decimals, &t.Amount, &t.Account, &narration); err != nil {
      return j, fmt.Errorf("Could not scan rows of query (%v)", err)
    }
    centsAmount, err := strconv.ParseFloat(t.Amount, 64)
    if err != nil {
      return j, fmt.Errorf("Could not process the amount as a float (%v)", err)
    }
    t.Amount = fmt.Sprintf("%.2f", centsAmount/math.Pow(10, decimals))
    j.LineItems = append(j.LineItems, t)
    j.Narration = narration
    j.Date = t.Date
  }
  if rows.Err() != nil {
    return j, fmt.Errorf("rows errored with (%v)", rows.Err())
  }

  j.LineItemCount = len(j.LineItems)
	return j, nil
}
