-- Tabela para controle de envio de Email
-- Criada em: 2025-12-10
-- Última atualização: 11/12/2025 15:00
-- Versão: 1.1.0

CREATE TABLE MENSAGEMEMAIL (
    -- Identificador único
    ID NUMBER(10) NOT NULL PRIMARY KEY,

    -- Código do cliente (FK para CLIENTES.CLICODIGO)
    CLICODIGO NUMBER(10),

    -- Remetente
    REMETENTE VARCHAR2(255) NOT NULL,

    -- Destinatário (apenas um por registro)
    DESTINATARIO VARCHAR2(255) NOT NULL,

    -- Assunto
    ASSUNTO VARCHAR2(500) NOT NULL,

    -- Corpo do e-mail (suporta HTML ou texto simples)
    CORPO CLOB NOT NULL,

    -- Tipo de corpo: 'text/plain' ou 'text/html'
    TIPO_CORPO VARCHAR2(20) DEFAULT 'text/plain' NOT NULL,

    -- Status do envio
    -- 0 = Pendente
    -- 2 = Enviado com sucesso
    -- 3 = Erro no envio (pode retentar)
    -- 4 = Falha permanente
    -- 125 = E-mail inválido
    STATUS_ENVIO NUMBER(3) DEFAULT 0 NOT NULL,

    -- Datas
    DATA_CADASTRO DATE DEFAULT SYSDATE NOT NULL,
    DATA_AGENDAMENTO DATE,
    DATA_ENVIO DATE,

    -- Controle de tentativas
    QTD_TENTATIVAS NUMBER(2) DEFAULT 0 NOT NULL,

    -- Mensagem de erro
    DETALHES_ERRO VARCHAR2(4000),

    -- ID no provider (para rastreamento)
    ID_PROVIDER VARCHAR2(100),

    -- Código do provider usado
    METODO_ENVIO NUMBER(10),

    -- Prioridade (1=Alta, 2=Normal, 3=Baixa)
    PRIORIDADE NUMBER(1) DEFAULT 2 NOT NULL,

    -- Anexo em base64 (CLOB para suportar arquivos grandes)
    ANEXO_REFERENCIA CLOB,
    ANEXO_NOME VARCHAR2(255),
    ANEXO_TIPO VARCHAR2(100),

    -- IP de origem (para disparo manual via HTTP)
    IP_ORIGEM VARCHAR2(50)
);

-- Índices para otimizar consultas
CREATE INDEX IDX_MENSAGEMEMAIL_STATUS ON MENSAGEMEMAIL(STATUS_ENVIO, DATA_AGENDAMENTO);
CREATE INDEX IDX_MENSAGEMEMAIL_CLI ON MENSAGEMEMAIL(CLICODIGO, STATUS_ENVIO);
CREATE INDEX IDX_MENSAGEMEMAIL_DEST ON MENSAGEMEMAIL(DESTINATARIO);
CREATE INDEX IDX_MENSAGEMEMAIL_DATA ON MENSAGEMEMAIL(DATA_CADASTRO);
CREATE INDEX IDX_MENSAGEMEMAIL_PRIO ON MENSAGEMEMAIL(PRIORIDADE, STATUS_ENVIO, DATA_AGENDAMENTO);

-- Sequence para geração de IDs
CREATE SEQUENCE SEQ_MENSAGEMEMAIL
    START WITH 1
    INCREMENT BY 1
    NOCACHE
    NOCYCLE;

-- Comentários nas colunas para documentação
COMMENT ON TABLE MENSAGEMEMAIL IS 'Controle de envio de mensagens de e-mail';
COMMENT ON COLUMN MENSAGEMEMAIL.ID IS 'Identificador único da mensagem';
COMMENT ON COLUMN MENSAGEMEMAIL.CLICODIGO IS 'Código do cliente (FK)';
COMMENT ON COLUMN MENSAGEMEMAIL.REMETENTE IS 'Endereço de e-mail do remetente';
COMMENT ON COLUMN MENSAGEMEMAIL.DESTINATARIO IS 'Endereço de e-mail do destinatário';
COMMENT ON COLUMN MENSAGEMEMAIL.ASSUNTO IS 'Assunto do e-mail';
COMMENT ON COLUMN MENSAGEMEMAIL.CORPO IS 'Corpo do e-mail (texto ou HTML)';
COMMENT ON COLUMN MENSAGEMEMAIL.TIPO_CORPO IS 'Tipo do corpo: text/plain ou text/html';
COMMENT ON COLUMN MENSAGEMEMAIL.STATUS_ENVIO IS '0=Pendente, 2=Enviado, 3=Erro, 4=Falha permanente, 125=Email inválido';
COMMENT ON COLUMN MENSAGEMEMAIL.DATA_CADASTRO IS 'Data/hora de criação do registro';
COMMENT ON COLUMN MENSAGEMEMAIL.DATA_AGENDAMENTO IS 'Data/hora agendada para envio (NULL=imediato)';
COMMENT ON COLUMN MENSAGEMEMAIL.DATA_ENVIO IS 'Data/hora do último envio';
COMMENT ON COLUMN MENSAGEMEMAIL.QTD_TENTATIVAS IS 'Número de tentativas de envio realizadas';
COMMENT ON COLUMN MENSAGEMEMAIL.DETALHES_ERRO IS 'Última mensagem de erro retornada';
COMMENT ON COLUMN MENSAGEMEMAIL.ID_PROVIDER IS 'ID da mensagem no provedor de e-mail';
COMMENT ON COLUMN MENSAGEMEMAIL.METODO_ENVIO IS 'Código do provedor utilizado';
COMMENT ON COLUMN MENSAGEMEMAIL.PRIORIDADE IS '1=Alta, 2=Normal, 3=Baixa';
COMMENT ON COLUMN MENSAGEMEMAIL.ANEXO_REFERENCIA IS 'Anexo em base64 (CLOB para suportar arquivos grandes)';
COMMENT ON COLUMN MENSAGEMEMAIL.ANEXO_NOME IS 'Nome do arquivo anexo';
COMMENT ON COLUMN MENSAGEMEMAIL.ANEXO_TIPO IS 'Tipo MIME do anexo';
COMMENT ON COLUMN MENSAGEMEMAIL.IP_ORIGEM IS 'IP de origem da requisição (disparo manual)';
