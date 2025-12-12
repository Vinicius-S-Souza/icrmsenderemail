-- Tabela para gerenciamento de templates de e-mail
-- Criada em: 12/12/2025 18:30
-- Versão: 1.2.0

CREATE TABLE TEMPLATEEMAIL (
    -- Identificador único
    ID NUMBER(10) NOT NULL PRIMARY KEY,

    -- Nome único do template (usado para identificação)
    NOME VARCHAR2(100) NOT NULL,

    -- Descrição do template (opcional)
    DESCRICAO VARCHAR2(500),

    -- Seção de cabeçalho (HTML)
    HEADER_HTML CLOB,

    -- Corpo principal do template (HTML) - obrigatório
    BODY_HTML CLOB NOT NULL,

    -- Seção de rodapé (HTML)
    FOOTER_HTML CLOB,

    -- Assunto padrão para e-mails usando este template
    ASSUNTO_PADRAO VARCHAR2(500),

    -- Status do template (1=Ativo, 0=Inativo)
    ATIVO NUMBER(1) DEFAULT 1 NOT NULL,

    -- Datas de auditoria
    DATA_CRIACAO DATE DEFAULT SYSDATE NOT NULL,
    DATA_ATUALIZACAO DATE DEFAULT SYSDATE NOT NULL,

    -- Usuário que criou o template
    CRIADO_POR VARCHAR2(100),

    -- Constraint de nome único
    CONSTRAINT UK_TEMPLATEEMAIL_NOME UNIQUE (NOME),

    -- Constraint para validar campo ATIVO
    CONSTRAINT CHK_TEMPLATEEMAIL_ATIVO CHECK (ATIVO IN (0, 1))
);

-- Sequence para geração de IDs
CREATE SEQUENCE SEQ_TEMPLATEEMAIL
    START WITH 1
    INCREMENT BY 1
    NOCACHE
    NOCYCLE;

-- Índices para otimizar consultas
CREATE INDEX IDX_TEMPLATEEMAIL_ATIVO ON TEMPLATEEMAIL(ATIVO);
CREATE INDEX IDX_TEMPLATEEMAIL_NOME ON TEMPLATEEMAIL(NOME);
CREATE INDEX IDX_TEMPLATEEMAIL_DATA ON TEMPLATEEMAIL(DATA_CRIACAO);

-- Comentários nas colunas para documentação
COMMENT ON TABLE TEMPLATEEMAIL IS 'Templates de e-mail HTML com suporte a seções e macros';
COMMENT ON COLUMN TEMPLATEEMAIL.ID IS 'Identificador único do template';
COMMENT ON COLUMN TEMPLATEEMAIL.NOME IS 'Nome único do template (usado para identificação)';
COMMENT ON COLUMN TEMPLATEEMAIL.DESCRICAO IS 'Descrição do propósito do template';
COMMENT ON COLUMN TEMPLATEEMAIL.HEADER_HTML IS 'HTML do cabeçalho (seção superior do e-mail)';
COMMENT ON COLUMN TEMPLATEEMAIL.BODY_HTML IS 'HTML do corpo principal (conteúdo do e-mail)';
COMMENT ON COLUMN TEMPLATEEMAIL.FOOTER_HTML IS 'HTML do rodapé (seção inferior do e-mail)';
COMMENT ON COLUMN TEMPLATEEMAIL.ASSUNTO_PADRAO IS 'Assunto padrão quando usar este template';
COMMENT ON COLUMN TEMPLATEEMAIL.ATIVO IS 'Status: 1=Ativo, 0=Inativo (soft delete)';
COMMENT ON COLUMN TEMPLATEEMAIL.DATA_CRIACAO IS 'Data/hora de criação do template';
COMMENT ON COLUMN TEMPLATEEMAIL.DATA_ATUALIZACAO IS 'Data/hora da última atualização';
COMMENT ON COLUMN TEMPLATEEMAIL.CRIADO_POR IS 'Usuário que criou o template';
