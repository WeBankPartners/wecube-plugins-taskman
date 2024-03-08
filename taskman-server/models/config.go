package models

import (
	"encoding/json"
	"github.com/WeBankPartners/go-common-lib/cipher"
	"github.com/WeBankPartners/go-common-lib/smtp"
	"github.com/WeBankPartners/go-common-lib/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type HttpServerConfig struct {
	Port              string `json:"port"`
	Cross             bool   `json:"cross"`
	ErrorTemplateDir  string `json:"error_template_dir"`
	ErrorDetailReturn bool   `json:"error_detail_return"`
}

type LogConfig struct {
	Level            string `json:"level"`
	LogDir           string `json:"log_dir"`
	AccessLogEnable  bool   `json:"access_log_enable"`
	DbLogEnable      bool   `json:"db_log_enable"`
	ArchiveMaxSize   int    `json:"archive_max_size"`
	ArchiveMaxBackup int    `json:"archive_max_backup"`
	ArchiveMaxDay    int    `json:"archive_max_day"`
	Compress         bool   `json:"compress"`
}

type DatabaseConfig struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DataBase string `json:"database"`
	MaxOpen  int    `json:"maxOpen"`
	MaxIdle  int    `json:"maxIdle"`
	Timeout  int    `json:"timeout"`
}

type WecubeConfig struct {
	BaseUrl       string `json:"base_url"`
	JwtSigningKey string `json:"jwt_signing_key"`
	SubSystemCode string `json:"sub_system_code"`
	SubSystemKey  string `json:"sub_system_key"`
}

type MailConfig struct {
	SenderName   string `json:"sender_name"`
	SenderMail   string `json:"sender_mail"`
	AuthServer   string `json:"auth_server"`
	AuthPassword string `json:"auth_password"`
	Ssl          string `json:"ssl"`
}

type AttachFileConfig struct {
	MinioAddress   string `json:"minio_address"`
	MinioAccessKey string `json:"minio_access_key"`
	MinioSecretKey string `json:"minio_secret_key"`
	Bucket         string `json:"bucket"`
	SSL            bool   `json:"ssl"`
}

type GlobalConfig struct {
	DefaultLanguage string           `json:"default_language"`
	HttpServer      HttpServerConfig `json:"http_server"`
	Log             LogConfig        `json:"log"`
	Database        DatabaseConfig   `json:"database"`
	RsaKeyPath      string           `json:"rsa_key_path"`
	Wecube          WecubeConfig     `json:"wecube"`
	Mail            MailConfig       `json:"mail"`
	AttachFile      AttachFileConfig `json:"attach_file"`
	EncryptSeed     string           `json:"encrypt_seed"`
}

var (
	Config                   *GlobalConfig
	CoreToken                *token.CoreToken
	ProcessFetchTabs         string
	MailEnable               bool
	MailSender               smtp.MailSender
	PriorityLevelMap         = map[int]string{1: "high", 2: "medium", 3: "low"}
	RequestTemplateImportMap = map[string]RequestTemplateExport{}
)

func InitConfig(configFile string) (errMessage string) {
	if configFile == "" {
		errMessage = "config file empty,use -c to specify configuration file"
		return
	}
	_, err := os.Stat(configFile)
	if os.IsExist(err) {
		errMessage = "config file not found," + err.Error()
		return
	}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		errMessage = "read config file fail," + err.Error()
		return
	}
	var c GlobalConfig
	err = json.Unmarshal(b, &c)
	if err != nil {
		errMessage = "parse file to json fail," + err.Error()
		return
	}
	rsaFileContent, _ := ioutil.ReadFile(c.RsaKeyPath)
	c.Database.Password, err = cipher.DecryptRsa(c.Database.Password, string(rsaFileContent))
	if err != nil {
		errMessage = "init database password fail,%s " + err.Error()
		return
	}
	Config = &c
	RequestTemplateImportMap = make(map[string]RequestTemplateExport)
	tmpCoreToken := token.CoreToken{}
	tmpCoreToken.BaseUrl = Config.Wecube.BaseUrl
	tmpCoreToken.JwtSigningKey = Config.Wecube.JwtSigningKey
	tmpCoreToken.SubSystemCode = Config.Wecube.SubSystemCode
	tmpCoreToken.SubSystemKey = Config.Wecube.SubSystemKey
	tmpCoreToken.InitCoreToken()
	CoreToken = &tmpCoreToken
	ProcessFetchTabs = os.Getenv("TASKMAN_PROCESS_TAGS")
	// attach file
	if strings.Contains(Config.AttachFile.MinioAddress, "//") {
		Config.AttachFile.MinioAddress = strings.Split(Config.AttachFile.MinioAddress, "//")[1]
	}
	Config.AttachFile.MinioAccessKey, _ = cipher.DecryptRsa(Config.AttachFile.MinioAccessKey, string(rsaFileContent))
	Config.AttachFile.MinioSecretKey, _ = cipher.DecryptRsa(Config.AttachFile.MinioSecretKey, string(rsaFileContent))
	// init mail
	MailEnable = false
	if c.Mail.AuthServer != "" && c.Mail.SenderMail != "" {
		MailEnable = true
		MailSender = smtp.MailSender{SenderName: c.Mail.SenderName, SenderMail: c.Mail.SenderMail, AuthServer: c.Mail.AuthServer, AuthPassword: c.Mail.AuthPassword, SSL: false}
		if c.Mail.Ssl == "Y" {
			MailSender.SSL = true
		}
		mailInitErr := MailSender.Init()
		if mailInitErr != nil {
			log.Printf("Init mail sender fail,%s \n", mailInitErr.Error())
			MailEnable = false
		} else {
			log.Println("Mail sender init success ")
		}
	} else {
		log.Println("Mail sender disable")
	}
	return
}
