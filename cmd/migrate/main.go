package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/config"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/database"
)

func main() {
	// Carregar configurações
	cfg, err := config.LoadConfig("dbinit.ini")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Conectar ao banco de dados
	fmt.Println("Conectando ao banco de dados Oracle...")
	db, err := database.ConnectOracle(database.DBConfig{
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		TNS:      cfg.Database.TNS,
	})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	fmt.Println("Conexão estabelecida com sucesso!")

	// Migração em etapas (Oracle não permite alterar VARCHAR2 para CLOB diretamente)
	fmt.Println("Executando migração para mudar ANEXO_REFERENCIA para CLOB...")

	ctx := context.Background()

	// Passo 1: Adicionar coluna temporária do tipo CLOB
	fmt.Println("1. Criando coluna temporária ANEXO_REFERENCIA_NEW do tipo CLOB...")
	_, err = db.ExecContext(ctx, "ALTER TABLE MENSAGEMEMAIL ADD (ANEXO_REFERENCIA_NEW CLOB)")
	if err != nil {
		log.Fatalf("Erro ao criar coluna temporária: %v", err)
	}
	fmt.Println("✅ Coluna temporária criada!")

	// Passo 2: Copiar dados da coluna antiga para a nova
	fmt.Println("2. Copiando dados para a nova coluna...")
	_, err = db.ExecContext(ctx, "UPDATE MENSAGEMEMAIL SET ANEXO_REFERENCIA_NEW = ANEXO_REFERENCIA")
	if err != nil {
		log.Fatalf("Erro ao copiar dados: %v", err)
	}
	fmt.Println("✅ Dados copiados!")

	// Passo 3: Remover coluna antiga
	fmt.Println("3. Removendo coluna antiga ANEXO_REFERENCIA...")
	_, err = db.ExecContext(ctx, "ALTER TABLE MENSAGEMEMAIL DROP COLUMN ANEXO_REFERENCIA")
	if err != nil {
		log.Fatalf("Erro ao remover coluna antiga: %v", err)
	}
	fmt.Println("✅ Coluna antiga removida!")

	// Passo 4: Renomear nova coluna para o nome original
	fmt.Println("4. Renomeando coluna ANEXO_REFERENCIA_NEW para ANEXO_REFERENCIA...")
	_, err = db.ExecContext(ctx, "ALTER TABLE MENSAGEMEMAIL RENAME COLUMN ANEXO_REFERENCIA_NEW TO ANEXO_REFERENCIA")
	if err != nil {
		log.Fatalf("Erro ao renomear coluna: %v", err)
	}
	fmt.Println("✅ Coluna renomeada!")

	// Passo 5: Atualizar comentário
	fmt.Println("5. Atualizando comentário da coluna...")
	_, err = db.ExecContext(ctx, "COMMENT ON COLUMN MENSAGEMEMAIL.ANEXO_REFERENCIA IS 'Anexo em base64 (CLOB para suportar arquivos grandes)'")
	if err != nil {
		log.Printf("⚠️  Aviso: Não foi possível atualizar o comentário: %v", err)
	} else {
		fmt.Println("✅ Comentário atualizado!")
	}

	fmt.Println("\n🎉 Migração concluída com sucesso!")
	os.Exit(0)
}
