package helper

import (
	"database/sql"
	"log"
)

var (
	dbHome  *sql.DB
	err     error
	publish = Publish
)

// User model
type User struct {
	ID          int    `json:"user_id"`
	Name        string `json:"full_name"`
	MSISDN      string `json:"msisdn"`
	Email       string `json:"user_email"`
	BirthDate   string `json:"birth_date"`
	CreatedTime string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Age         int    `json:"age"`
}

func init() {
	dbHome, err = sql.Open("postgres", "postgres://mh180319:RI6zX88gB2SgAp@devel-postgre.tkpd/tokopedia-dev-db?sslmode=disable")
	if err != nil {
		log.Print(err)
	}
}

// GetUsers by name
func GetUsers(key string) []User {
	where := ""

	if key != "" {
		where = " where u.full_name like '%" + key + "%' "
	}
	stmt, err := dbHome.Prepare("select u.user_id, u.full_name,	u.msisdn, u.user_email,	to_char(u.birth_date, 'YYYY/MM/DD') as birth_date, to_char(u.create_time, 'YYYY/MM/DD') as create_time, to_char(u.update_time, 'YYYY/MM/DD') as update_time,  cast(EXTRACT(YEAR FROM age(now(), u.birth_date)) as int) as age  from	public.ws_user u " + where + " limit 10")

	userList := []User{}

	rows, err := stmt.Query()
	if err != nil {
		log.Print(err)
	}

	defer rows.Close()
	for rows.Next() {
		u := &User{}
		rows.Scan(&u.ID, &u.Name, &u.MSISDN, &u.Email, &u.BirthDate, &u.CreatedTime, &u.UpdateTime, &u.Age)

		if u.MSISDN == "" {
			u.MSISDN = "-"
		}

		if u.BirthDate == "" {
			u.BirthDate = "-"
		}

		if u.CreatedTime == "" {
			u.CreatedTime = "-"
		}

		if u.UpdateTime == "" {
			u.UpdateTime = "-"
		}

		userList = append(userList, *u)
	}

	return userList

}

// CloseConnection handler
func CloseConnection() {
	dbHome.Close()
}
