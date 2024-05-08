package dumper

import (
	"database/sql"
	"errors"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

type table struct {
	Name   string
	SQL    string
	Values string
}

type dump struct {
	DumpVersion   string
	ServerVersion string
	Tables        []*table
	CompleteTime  string
}

const version = "0.2.2"

const tmpl = `-- Go SQL Dump {{ .DumpVersion }}
--
-- ------------------------------------------------------
-- Server version	{{ .ServerVersion }}

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


{{range .Tables}}
--
-- Table structure for table {{ .Name }}
--

DROP TABLE IF EXISTS {{ .Name }};
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
{{ .SQL }};
/*!40101 SET character_set_client = @saved_cs_client */;
--
-- Dumping data for table {{ .Name }}
--

LOCK TABLES {{ .Name }} WRITE;
/*!40000 ALTER TABLE {{ .Name }} DISABLE KEYS */;
{{ if .Values }}
INSERT INTO {{ .Name }} VALUES {{ .Values }};
{{ end }}
/*!40000 ALTER TABLE {{ .Name }} ENABLE KEYS */;
UNLOCK TABLES;
{{ end }}
-- Dump completed on {{ .CompleteTime }}
`

func (d *Dumper) Dump() (string, error) {
	name := d.name
	p := path.Join(d.dir, name+".sql")

	if e, _ := exists(p); e {
		return p, errors.New("Dump '" + name + "' already exists.")
	}

	f, err := os.Create(p)

	if err != nil {
		return p, err
	}

	defer f.Close()

	data := dump{
		DumpVersion: version,
		Tables:      make([]*table, 0),
	}

	if data.ServerVersion, err = getServerVersion(d.db); err != nil {
		return p, err
	}

	tables, err := getTables(d.db)
	if err != nil {
		return p, err
	}

	for _, name := range tables {
		if t, err := createTable(d.db, name); err == nil {
			data.Tables = append(data.Tables, t)
		} else {
			return p, err
		}
	}

	data.CompleteTime = time.Now().String()

	t, err := template.New("sqldump").Parse(tmpl)
	if err != nil {
		return p, err
	}
	if err = t.Execute(f, data); err != nil {
		return p, err
	}

	return p, nil
}

func getTables(db *sql.DB) ([]string, error) {
	tables := make([]string, 0)

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return tables, err
	}
	defer rows.Close()

	for rows.Next() {
		var table sql.NullString
		if err := rows.Scan(&table); err != nil {
			return tables, err
		}
		tables = append(tables, table.String)
	}
	return tables, rows.Err()
}

func getServerVersion(db *sql.DB) (string, error) {
	var server_version sql.NullString
	if err := db.QueryRow("SELECT version()").Scan(&server_version); err != nil {
		return "", err
	}
	return server_version.String, nil
}

func createTable(db *sql.DB, name string) (*table, error) {
	var err error
	t := &table{Name: name}

	if t.SQL, err = createTableSQL(db, name); err != nil {
		return nil, err
	}

	if t.Values, err = createTableValues(db, name); err != nil {
		return nil, err
	}

	return t, nil
}

func createTableSQL(db *sql.DB, name string) (string, error) {
	var table_return sql.NullString
	var table_sql sql.NullString
	err := db.QueryRow("SHOW CREATE TABLE "+name).Scan(&table_return, &table_sql)

	if err != nil {
		return "", err
	}
	if table_return.String != name {
		return "", errors.New("returned table is not the same as requested table")
	}

	return table_sql.String, nil
}

func createTableValues(db *sql.DB, name string) (string, error) {
	rows, err := db.Query("SELECT * FROM " + name)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if len(columns) == 0 {
		return "", errors.New("No columns in table " + name + ".")
	}

	data_text := make([]string, 0)
	for rows.Next() {

		data := make([]*sql.NullString, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range data {
			ptrs[i] = &data[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return "", err
		}

		dataStrings := make([]string, len(columns))

		for key, value := range data {
			if value != nil && value.Valid {
				dataStrings[key] = value.String
			}
		}

		data_text = append(data_text, "('"+strings.Join(dataStrings, "','")+"')")
	}

	return strings.Join(data_text, ","), rows.Err()
}
