package main

import ( 
    _ "database/sql"
    "github.com/jmoiron/sqlx"
    _ "log"
    _ "github.com/lib/pq"
)

type Count struct {
    CountVal int `db:"count_val"`
}

func getUpdateCount() (int, error) {
    db, err := sqlx.Connect("postgres", "user=irc password=ircsecurepassword dbname=base sslmode=disable")
    if err != nil {
        return 0, err
    }

    count := Count{}

    db.Get(&count, "SELECT * FROM irc.count ORDER BY count_val")
    count.CountVal += 1

    tx := db.MustBegin()
    tx.MustExec("UPDATE irc.count set count_val = $1;", count.CountVal)
    tx.Commit()

    return count.CountVal, nil
}
