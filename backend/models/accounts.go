package models

import (
	"context"
  "encoding/json"
  "strings"
	"fmt"
  "net/http"
	"time"
  "io/ioutil"

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

type AccountDetail struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tags []string `json:"tags"`
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

func GetAccountCommand(id string) (AccountDetail, error) {
	log.Trace("Calling Get Journal Command function")
	
  var account AccountDetail
	account.Tags = []string{}

  db := backend.GetConnection()

	queryAccountDetail := `
	SELECT 
		accounts.account_id,
		accounts.name 
	FROM accounts 
  WHERE name = ?
  LIMIT 1
	;`

	log.Debug("Querying Database")
	rows, err := db.Query(queryAccountDetail, id)
	if err != nil {
		log.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

  for rows.Next() {
    if err := rows.Scan(&account.ID, &account.Name); err != nil {
      return account, fmt.Errorf("Could not scan rows of query (%v)", err)
    }
  }
  if rows.Err() != nil {
    return account, fmt.Errorf("rows errored with (%v)", rows.Err())
  }

	tagsQuery := `
		SELECT tag_name
		FROM   tags
					 JOIN account_tag
						 ON account_tag.tag_id = tags.tag_id
					 JOIN accounts
						 ON accounts.account_id = account_tag.account_id
		WHERE  accounts.NAME = ?;
		`


  log.Debugf("Querying Database for Tags on Account: %s", account.Name)

  rows, err = db.Query(tagsQuery, account.Name)
  if err != nil {
    return account, fmt.Errorf("tags query errored with (%v)", rows.Err())
  }

  for rows.Next() {
    var tag string
    if err := rows.Scan(&tag); err != nil {
      return account, fmt.Errorf("tags scan errored with (%v)", err)
    }
    log.Debugf("Tag found: %s", tag)
    account.Tags = append(account.Tags, tag)
  }


	return account, nil
}

func DeleteAccountTagCommand(id, tag string) error {
	log.Trace("Calling Delete Account Tag function")

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
		Tag: []string{tag},
	}
	r, err := client.DeleteTag(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Delete Transaction Method (%v)", err)
	}
	log.Infof("Delete Transaction Response: %s", r.GetMessage())

	return nil
}

func ImportAccountsCommand(name string) (error) {
	log.Trace("Calling Import Accounts Command function")
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
  url := fmt.Sprintf("https://raw.githubusercontent.com/darcys22/godbledger-assets/master/Accounts/%s.json",  strings.Title(name))
	resp, err := http.Get(url)
	if err != nil {
    return fmt.Errorf("Error fetching account from url with (%v)", err)
	}
  defer resp.Body.Close()
  accountsJSONstring, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    return fmt.Errorf("Error reading accounts JSON string (%v)", err)
	}
  var accounts = []AccountDetail{}
  err = json.Unmarshal(accountsJSONstring, &accounts)
	if err != nil {
    if terr, ok := err.(*json.UnmarshalTypeError); ok {
        fmt.Printf("Failed to unmarshal field %s \n", terr.Field)
    } else {
        fmt.Println(err)
    }
    return fmt.Errorf("Error parsing accounts JSON string (%v)", err)
	}
  // using for loop
  for _, account := range accounts{
      req := &pb.AccountTagRequest{
        Account:  account.Name,
        Tag:      account.Tags,
      }
      r, err := client.AddAccount(ctxTimeout, req)
      if err != nil {
        return fmt.Errorf("Could not call Add Account Method (%v)", err)
      }
      log.Debugf("Add Account Response: %s", r.GetMessage())
  }
  return nil
}
