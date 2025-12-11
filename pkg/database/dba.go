package database

import (
	"database/sql"
	"fmt"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/logger"
	"github.com/go-ini/ini"
	_ "github.com/godror/godror"
	"go.uber.org/zap"
)

type DBConfig struct {
	Username string
	Password string
	TNS      string
}

func LoadDBConfig(path string) (DBConfig, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return DBConfig{}, err
	}
	section := cfg.Section("oracle")
	return DBConfig{
		Username: section.Key("username").String(),
		Password: section.Key("password").String(),
		TNS:      section.Key("tns").String(),
	}, nil
}

func ConnectOracle(conf DBConfig) (*sql.DB, error) {
	logger.Info("Tentando conectar com Oracle",
		zap.String("username", conf.Username),
		zap.String("tns", conf.TNS))

	// Monta a string de conexão usando diferentes formatos para compatibilidade
	// Formato alternativo: usuario/senha@tns
	connStr := fmt.Sprintf("%s/%s@%s", conf.Username, conf.Password, conf.TNS)
	logger.Info("Tentando conexão com formato simples", zap.String("connStr", fmt.Sprintf("%s/***@%s", conf.Username, conf.TNS)))

	// Tenta conectar com o primeiro formato
	db, err := sql.Open("godror", connStr)
	if err != nil {
		logger.Warn("Falha com formato simples, tentando formato completo", zap.Error(err))
		// Se falhar, tenta o formato com aspas
		connStr = fmt.Sprintf(`user="%s" password="%s" connectString="%s"`, conf.Username, conf.Password, conf.TNS)
		logger.Info("Tentando conexão com formato completo")
		db, err = sql.Open("godror", connStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao abrir conexão Oracle: %v", err)
		}
	}

	// Configura parâmetros adicionais
	if err := ConfigureOracle(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao configurar conexão Oracle: %v", err)
	}

	logger.Info("Conexão Oracle estabelecida com sucesso")
	return db, nil
}
