package main

import (
	"database/sql"
	"log"
)

type SubjectTable struct {
	db *sql.DB
}

func (st *SubjectTable) Connect(url string) (err error) {

	st.db, err = sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = st.db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (st *SubjectTable) CreateTable() error {

	create := `
	CREATE TABLE IF NOT EXISTS subjectTable (
		uid INT(11) NOT NULL ,
		product_id TEXT,
		transaction_id TEXT,
		original_transaction_id TEXT,
		purchase_date DATETIME,
		original_purchase_date DATETIME,
		expires_date DATETIME,
		in_app_ownership_type TEXT ,
		veri_data TEXT,
		auto_renewal_info TEXT
	);`

	if _, err := st.db.Exec(create); err != nil {
		return err
	}
	return nil
}
func (st *SubjectTable) TransactionIdisExist(transaction_id string) (bool, error) {

	query := `
	select count(transaction_id) from subjectTable where transaction_id = ?
	`
	row := st.db.QueryRow(query, transaction_id)

	var count int
	err := row.Scan(&count)

	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (st *SubjectTable) InsertLatestReceiptInfo(
	uid int,
	product_id string,
	transaction_id string,
	original_transaction_id string,
	purchase_date string,
	original_purchase_date string,
	expires_date string,
	in_app_ownership_type string,
	veri_data string,
	auto_renewal_info string,

) error {

	insertStr := `
		insert into subjectTable(
			uid,
			product_id,
			transaction_id,
			original_transaction_id,
			purchase_date,
			original_purchase_date,
			expires_date,
			in_app_ownership_type,
			veri_data,
			auto_renewal_info
		) values(?,?,?,?,?,?,?,?,?,?)`

	_, err := st.db.Exec(insertStr,
		uid,
		product_id,
		transaction_id,
		original_transaction_id,
		purchase_date,
		original_purchase_date,
		expires_date,
		in_app_ownership_type,
		veri_data,
		auto_renewal_info,
	)
	if err != nil {
		return err
	}
	return nil
}
