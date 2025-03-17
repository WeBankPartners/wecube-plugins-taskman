package dao

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"xorm.io/core"
	"xorm.io/xorm"
	xorm_log "xorm.io/xorm/log"
)

var (
	_ xorm_log.Logger = &dbLogger{}
	X *xorm.Engine
)

func InitDatabase() (engine *xorm.Engine, err error) {
	connStr := fmt.Sprintf("%s:%s@%s(%s)/%s?collation=utf8mb4_unicode_ci&allowNativePasswords=true",
		models.Config.Database.User, models.Config.Database.Password, "tcp", fmt.Sprintf("%s:%s", models.Config.Database.Server, models.Config.Database.Port), models.Config.Database.DataBase)
	engine, err = xorm.NewEngine("mysql", connStr)
	if err != nil {
		log.Error(nil, log.LOGGER_APP, "Init database connect fail", zap.Error(err))
		return nil, err
	}
	engine.SetMaxIdleConns(models.Config.Database.MaxIdle)
	engine.SetMaxOpenConns(models.Config.Database.MaxOpen)
	engine.SetConnMaxLifetime(time.Duration(models.Config.Database.Timeout) * time.Second)
	if models.Config.Log.DbLogEnable {
		engine.SetLogger(&dbLogger{LogLevel: 1, ShowSql: true, Logger: log.DatabaseLogger})
	}
	// 使用驼峰式映射
	engine.SetMapper(core.SnakeMapper{})
	log.Info(nil, log.LOGGER_APP, "Success init database connect !!")
	X = engine
	return
}

// logExecuteSql 打印执行sql
func logExecuteSql(session *xorm.Session, module, method string, objectParam interface{}, affected int64, err error) {
	sql, param := session.LastSQL()
	log.Debug(nil, log.LOGGER_DB, fmt.Sprintf("%s exec %s sql", module, method), zap.String("sql", sql), log.JsonObj("param", param), log.JsonObj("objectParam", objectParam), zap.Int64("affected", affected), zap.Error(err))
}

type dbLogger struct {
	LogLevel xorm_log.LogLevel
	ShowSql  bool
	Logger   *zap.SugaredLogger
}

func (d *dbLogger) Debug(v ...interface{}) {
	d.Logger.Debugw(fmt.Sprint(v...))
}

func (d *dbLogger) Debugf(format string, v ...interface{}) {
	d.Logger.Debugw(fmt.Sprintf(format, v...))
}

func (d *dbLogger) Error(v ...interface{}) {
	d.Logger.Errorw(fmt.Sprint(v...))
}

func (d *dbLogger) Errorf(format string, v ...interface{}) {
	d.Logger.Errorw(fmt.Sprintf(format, v...))
}

func (d *dbLogger) Info(v ...interface{}) {
	d.Logger.Infow(fmt.Sprint(v...))
}

func (d *dbLogger) Infof(format string, v ...interface{}) {
	if len(v) < 4 {
		d.Logger.Infow(fmt.Sprintf(format, v...))
		return
	}
	var costMs float64 = 0
	costTime := fmt.Sprintf("%s", v[3])
	if strings.Contains(costTime, "µs") {
		costMs, _ = strconv.ParseFloat(strings.ReplaceAll(costTime, "µs", ""), 64)
		costMs = costMs / 1000
	} else if strings.Contains(costTime, "ms") {
		costMs, _ = strconv.ParseFloat(costTime[:len(costTime)-2], 64)
	} else if strings.Contains(costTime, "s") && !strings.Contains(costTime, "m") {
		costMs, _ = strconv.ParseFloat(costTime[:len(costTime)-1], 64)
		costMs = costMs * 1000
	} else {
		costTime = costTime[:len(costTime)-1]
		mIndex := strings.Index(costTime, "m")
		minTime, _ := strconv.ParseFloat(costTime[:mIndex], 64)
		secTime, _ := strconv.ParseFloat(costTime[mIndex+1:], 64)
		costMs = (minTime*60 + secTime) * 1000
	}
	d.Logger.Infow("db_log", zap.String("sql", fmt.Sprintf("%s", v[1])), zap.String("param", fmt.Sprintf("%v", v[2])), zap.Float64("cost_ms", costMs))
}

func (d *dbLogger) Warn(v ...interface{}) {
	d.Logger.Warnw(fmt.Sprint(v...))
}

func (d *dbLogger) Warnf(format string, v ...interface{}) {
	d.Logger.Warnw(fmt.Sprintf(format, v...))
}

func (d *dbLogger) Level() xorm_log.LogLevel {
	return d.LogLevel
}

func (d *dbLogger) SetLevel(l xorm_log.LogLevel) {
	d.LogLevel = l
}

func (d *dbLogger) ShowSQL(b ...bool) {
	d.ShowSql = b[0]
}

func (d *dbLogger) IsShowSQL() bool {
	return d.ShowSql
}

func QueryCount(sql string, params ...interface{}) int {
	sql = "SELECT COUNT(1) FROM ( " + sql + " ) sub_query"
	params = append([]interface{}{sql}, params...)
	queryRows, err := X.QueryString(params...)
	if err != nil || len(queryRows) == 0 {
		log.Error(nil, log.LOGGER_APP, "Query sql count message fail", zap.Error(err))
		return 0
	}
	if _, b := queryRows[0]["COUNT(1)"]; b {
		countNum, _ := strconv.Atoi(queryRows[0]["COUNT(1)"])
		return countNum
	}
	return 0
}

func GetJsonToXormMap(input interface{}) (resultMap map[string]string, idKeyName string) {
	var label string
	resultMap = make(map[string]string)
	t := reflect.TypeOf(input)
	for i := 0; i < t.NumField(); i++ {
		label = t.Field(i).Tag.Get("xorm")
		if strings.Contains(label, "pk") && t.Field(i).Tag.Get("primary-key") != "" {
			label = t.Field(i).Tag.Get("primary-key")
		}
		resultMap[t.Field(i).Tag.Get("json")] = label
		if i == 0 {
			idKeyName = label
		}
	}
	return resultMap, idKeyName
}

func TransFiltersToSQL(queryParam *models.QueryRequestParam, transParam *models.TransFiltersParam) (filterSql, queryColumn string, param []interface{}) {
	if transParam.Prefix != "" && !strings.HasSuffix(transParam.Prefix, ".") {
		transParam.Prefix = transParam.Prefix + "."
	}
	if transParam.IsStruct {
		transParam.KeyMap, transParam.PrimaryKey = GetJsonToXormMap(transParam.StructObj)
	}
	for _, filter := range queryParam.Filters {
		if transParam.KeyMap[filter.Name] == "" || transParam.KeyMap[filter.Name] == "-" {
			continue
		}
		if filter.Operator == "eq" {
			filterSql += fmt.Sprintf(" AND %s%s=? ", transParam.Prefix, transParam.KeyMap[filter.Name])
			param = append(param, filter.Value)
		} else if filter.Operator == "contains" || filter.Operator == "like" {
			filterSql += fmt.Sprintf(" AND %s%s LIKE ? ", transParam.Prefix, transParam.KeyMap[filter.Name])
			param = append(param, fmt.Sprintf("%%%s%%", filter.Value))
		} else if filter.Operator == "in" {
			inValueList := filter.Value.([]interface{})
			var inValueStringList []string
			for _, inValueInterfaceObj := range inValueList {
				if inValueInterfaceObj == nil {
					inValueStringList = append(inValueStringList, "")
				} else {
					inValueStringList = append(inValueStringList, inValueInterfaceObj.(string))
				}
			}
			tmpSpecSql, tmpListParams := CreateListParams(inValueStringList, "")
			filterSql += fmt.Sprintf(" AND %s%s in (%s) ", transParam.Prefix, transParam.KeyMap[filter.Name], tmpSpecSql)
			param = append(param, tmpListParams...)
		} else if filter.Operator == "lt" {
			filterSql += fmt.Sprintf(" AND %s%s<=? ", transParam.Prefix, transParam.KeyMap[filter.Name])
			param = append(param, filter.Value)
		} else if filter.Operator == "gt" {
			filterSql += fmt.Sprintf(" AND %s%s>=? ", transParam.Prefix, transParam.KeyMap[filter.Name])
			param = append(param, filter.Value)
		} else if filter.Operator == "ne" || filter.Operator == "neq" {
			filterSql += fmt.Sprintf(" AND %s%s!=? ", transParam.Prefix, transParam.KeyMap[filter.Name])
			param = append(param, filter.Value)
		} else if filter.Operator == "notNull" || filter.Operator == "isnot" {
			filterSql += fmt.Sprintf(" AND %s%s is not null ", transParam.Prefix, transParam.KeyMap[filter.Name])
		} else if filter.Operator == "null" || filter.Operator == "is" {
			filterSql += fmt.Sprintf(" AND %s%s is null ", transParam.Prefix, transParam.KeyMap[filter.Name])
		}
	}
	if queryParam.Sorting != nil {
		if transParam.KeyMap[queryParam.Sorting.Field] == "" || transParam.KeyMap[queryParam.Sorting.Field] == "-" {
			queryParam.Sorting.Field = transParam.PrimaryKey
		} else {
			queryParam.Sorting.Field = transParam.KeyMap[queryParam.Sorting.Field]
		}
		if queryParam.Sorting.Asc {
			filterSql += fmt.Sprintf(" ORDER BY %s%s ASC ", transParam.Prefix, queryParam.Sorting.Field)
		} else {
			filterSql += fmt.Sprintf(" ORDER BY %s%s DESC ", transParam.Prefix, queryParam.Sorting.Field)
		}
	}
	if len(queryParam.ResultColumns) > 0 {
		for _, resultColumn := range queryParam.ResultColumns {
			if transParam.KeyMap[resultColumn] == "" || transParam.KeyMap[resultColumn] == "-" {
				continue
			}
			queryColumn += fmt.Sprintf("%s%s,", transParam.Prefix, transParam.KeyMap[resultColumn])
		}
	}
	if queryColumn == "" {
		queryColumn = " * "
	} else {
		queryColumn = queryColumn[:len(queryColumn)-1]
	}
	return
}

func TransPageInfoToSQL(pageInfo models.PageInfo) (pageSql string, param []interface{}) {
	pageSql = " LIMIT ?,? "
	param = append(param, pageInfo.StartIndex)
	param = append(param, pageInfo.PageSize)
	return
}

type ExecAction struct {
	Sql   string
	Param []interface{}
}

func Transaction(actions []*ExecAction) error {
	if len(actions) == 0 {
		log.Warn(nil, log.LOGGER_APP, "Transaction is empty,nothing to do")
		return fmt.Errorf("SQL exec transaction is empty,nothing to do,please check server log ")
	}
	for i, action := range actions {
		if action == nil {
			return fmt.Errorf("SQL exec transaction index%d is nill error,please check server log", i)
		}
	}
	session := X.NewSession()
	err := session.Begin()
	for _, action := range actions {
		params := make([]interface{}, 0)
		params = append(params, action.Sql)
		params = append(params, action.Param...)
		_, err = session.Exec(params...)
		if err != nil {
			session.Rollback()
			break
		}
	}
	if err == nil {
		err = session.Commit()
	}
	session.Close()
	return err
}

func TransactionWithoutForeignCheck(actions []*ExecAction) error {
	if len(actions) == 0 {
		log.Warn(nil, log.LOGGER_APP, "Transaction is empty,nothing to do")
		return fmt.Errorf("SQL exec transaction is empty,nothing to do,please check server log ")
	}
	for i, action := range actions {
		if action == nil {
			return fmt.Errorf("SQL exec transaction index%d is nill error,please check server log", i)
		}
	}
	session := X.NewSession()
	err := session.Begin()
	if err != nil {
		return err
	}
	session.Exec("SET FOREIGN_KEY_CHECKS=0")
	for _, action := range actions {
		params := make([]interface{}, 0)
		params = append(params, action.Sql)
		params = append(params, action.Param...)
		_, err = session.Exec(params...)
		if err != nil {
			session.Rollback()
			break
		}
	}
	if err == nil {
		err = session.Commit()
	}
	session.Exec("SET FOREIGN_KEY_CHECKS=1")
	session.Close()
	return err
}

func CreateListParams(inputList []string, prefix string) (specSql string, paramList []interface{}) {
	if len(inputList) > 0 {
		var specList []string
		for _, v := range inputList {
			specList = append(specList, "?")
			paramList = append(paramList, prefix+v)
		}
		specSql = strings.Join(specList, ",")
	}
	return
}

func CombineDBSql(input ...interface{}) string {
	var buf strings.Builder
	fmt.Fprint(&buf, input...)
	return buf.String()
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func GetInsertTableExecAction(tableName string, data interface{}, transNullStr map[string]string) (action *ExecAction, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	execParams := []interface{}{}
	columnStr := ""
	valueStr := ""
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	tagName := "xorm"
	for i := 0; i < t.NumField(); i++ {
		fType := t.Field(i)
		fTag := fType.Tag.Get(tagName)
		if fTag == "-" {
			continue
		}

		if i > 0 {
			columnStr += ","
			valueStr += ","
		}
		// columnStr += t.Field(i).Tag.Get("xorm")
		columnStr += "`" + t.Field(i).Tag.Get("xorm") + "`"
		valueStr += "?"

		if len(transNullStr) > 0 {
			if _, ok := transNullStr[t.Field(i).Tag.Get("xorm")]; ok {
				execParams = append(execParams, NewNullString(v.FieldByName(t.Field(i).Name).String()))
			} else {
				execParams = append(execParams, v.FieldByName(t.Field(i).Name).Interface())
			}
		} else {
			execParams = append(execParams, v.FieldByName(t.Field(i).Name).Interface())
		}
	}
	execSqlCmd := CombineDBSql("INSERT INTO ", tableName, "(", columnStr, ") VALUE (", valueStr, ")")
	action = &ExecAction{Sql: execSqlCmd, Param: execParams}
	return
}
