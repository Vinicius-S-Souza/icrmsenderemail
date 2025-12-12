-- Alteração da tabela MENSAGEMEMAIL para suportar templates
-- Data: 12/12/2025 18:30
-- Versão: 1.2.0

-- Adicionar coluna para vincular e-mail ao template utilizado
ALTER TABLE MENSAGEMEMAIL ADD TEMPLATE_ID NUMBER(10);

-- Adicionar foreign key para TEMPLATEEMAIL
ALTER TABLE MENSAGEMEMAIL ADD CONSTRAINT FK_MENSAGEM_TEMPLATE
    FOREIGN KEY (TEMPLATE_ID) REFERENCES TEMPLATEEMAIL(ID);

-- Criar índice para otimizar consultas
CREATE INDEX IDX_MENSAGEMEMAIL_TEMPLATE ON MENSAGEMEMAIL(TEMPLATE_ID);

-- Adicionar comentário na coluna
COMMENT ON COLUMN MENSAGEMEMAIL.TEMPLATE_ID IS 'ID do template utilizado para gerar este e-mail';
