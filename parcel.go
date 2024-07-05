package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализовано добавление строки в таблицу parcel,
	// при этом используются данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	// возвращается идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализовано чтение строки по заданному number
	// здесь из таблицы возвращается только одна строка
	// объект Parcel заполнен данными из таблицы
	p := Parcel{}
	row := s.db.QueryRow("SELECT *	FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализовано чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	// срез Parcel заполнен данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, err
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализовано обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализовано обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE status = :status AND number = :number",
		sql.Named("address", address),
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("number", number))
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализовано удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE status = :status AND number = :number",
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("number", number))
	return err
}
