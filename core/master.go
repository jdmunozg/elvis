package core

import (
	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/envar"
	. "github.com/cgalvisleon/elvis/jdb"
	. "github.com/cgalvisleon/elvis/json"
	. "github.com/cgalvisleon/elvis/msg"
	. "github.com/cgalvisleon/elvis/utilities"
	_ "github.com/joho/godotenv/autoload"
)

func GetMarterNodeById(db int, id string) (Item, error) {
	sql := `
	SELECT
	A._DATA||
  jsonb_build_object(
    'mode', A.MODE,
		'index', A.INDEX
  ) AS _DATA
	FROM core.NODES A
	WHERE A._ID=$1
	LIMIT 1;`

	item, err := DBQueryDataOne(db, sql, id)
	if err != nil {
		return Item{}, err
	}

	delete(item.Result, "password")

	return item, nil
}

func UpSetMasterNode(db int, id string, mode int, driver, host string, port int, dbname, user, password string) (Item, error) {
	exist, err := ExistTable(db, "core", "NODES")
	if err != nil {
		return Item{}, err
	}

	if !exist {
		return Item{}, console.AlertF(MARTER_NOT_FOUNT, host)
	}

	if !ValidId(id) {
		return Item{}, console.AlertF(MSG_ATRIB_REQUIRED, "id")
	}

	current, err := GetMarterNodeById(db, id)
	if err != nil {
		return Item{}, err
	}

	now := Now()
	data := Json{
		"driver": driver,
		"host":   host,
		"port":   port,
		"dbname": dbname,
		"user":   user,
	}

	if current.Ok {
		sql := `
		UPDATE core.NODES SET
		DATE_UPDATE=$2,
		MODE=$3,
		PASSWORD=$4,
		_DATA=$5
		WHERE _ID=$1
		RETURNING INDEX;`

		item, err := DBQueryOne(db, sql, id, now, mode, password, data)
		if err != nil {
			return Item{}, err
		}

		return Item{
			Ok: item.Ok,
			Result: Json{
				"message": RECORD_UPDATE,
				"_id":     id,
				"index":   item.Index(),
			},
		}, nil
	}

	index := GetSerie("core.NODES")
	sql := `
		INSERT INTO core.NODES(DATE_MAKE, DATE_UPDATE, _ID, MODE, _DATA, INDEX)
		VALUES($1, $1, $2, $3, $4, $5)
		RETURNING INDEX;`

	item, err := DBQueryOne(db, sql, now, id, mode, data, index)
	if err != nil {
		return Item{}, err
	}

	return Item{
		Ok: item.Ok,
		Result: Json{
			"message": RECORD_CREATE,
			"_id":     id,
			"index":   item.Index(),
		},
	}, nil
}

func JoinToMaster() error {
	if err := DefineCoreSchema(); err != nil {
		return console.PanicE(err)
	}

	sql := `
	SELECT A.*
	FROM core.MODE A
	LIMIT 1;`

	item, err := QueryOne(sql)
	if err != nil {
		return err
	}

	id := item.Id()
	mode := item.ValInt(ModeNone, "_data", "mode")
	driver := item.ValStr(Postgres, "_data", "driver")
	host := item.ValStr("", "_data", "host")
	port := item.ValInt(5432, "_data", "port")
	dbname := item.ValStr("", "_data", "dbname")
	user := item.ValStr("", "_data", "user")
	password := item.ValStr("", "password")

	ModeId = id
	ModeTp = mode

	if ContainsInt([]int{ModeNone, ModeIdle}, ModeTp) {
		return nil
	}

	if driver == "" {
		return console.AlertF(MSG_ATRIB_REQUIRED, "driver")
	}

	if host == "" {
		return console.AlertF(MSG_ATRIB_REQUIRED, "host")
	}

	if dbname == "" {
		return console.AlertF(MSG_ATRIB_REQUIRED, "dbname")
	}

	if user == "" {
		return console.AlertF(MSG_ATRIB_REQUIRED, "user")
	}

	if password == "" {
		return console.AlertF(MSG_ATRIB_REQUIRED, "password")
	}

	idx, err := Connected(driver, host, port, dbname, user, password)
	if err != nil {
		return err
	}

	driver = EnvarStr("", "DB_DRIVE")
	host = EnvarStr("", "DB_HOST")
	port = EnvarInt(5432, "DB_PORT")
	dbname = EnvarStr("", "DB_NAME")
	user = EnvarStr("", "DB_USER")
	password = EnvarStr("", "DB_PASSWORD")
	_, err = UpSetMasterNode(idx, ModeId, ModeTp, driver, host, port, dbname, user, password)
	if err != nil {
		return err
	}

	c := make(chan int, 1)

	go SyncSeries(idx, c)

	select {
	case n := <-c:
		MasterIdx = n
	}

	console.LogKF("MASTER", "Join to master:%s:%d", host, port)

	return nil
}
