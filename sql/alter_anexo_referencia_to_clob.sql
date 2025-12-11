-- Script de migração para alterar ANEXO_REFERENCIA para CLOB
-- Data de criação: 11/12/2025 15:00
-- Versão: 1.0.0
--
-- Objetivo: Permitir anexos maiores em base64
-- A coluna VARCHAR2(500) não suporta arquivos grandes

-- Alterar coluna ANEXO_REFERENCIA para CLOB
ALTER TABLE MENSAGEMEMAIL MODIFY (ANEXO_REFERENCIA CLOB);

-- Comentário atualizado
COMMENT ON COLUMN MENSAGEMEMAIL.ANEXO_REFERENCIA IS 'Anexo em base64 (CLOB para suportar arquivos grandes)';
