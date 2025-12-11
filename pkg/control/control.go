package control

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"go.uber.org/zap"
)

const StopFileName = "icrmsenderemail.stop"

// WatchStopFile monitora um arquivo de controle para parar a aplicação.
// Agora suporta múltiplos diretórios candidatos para robustez entre Linux/Windows:
//  1. Caminho absoluto definido em ICRMSENDEREMAIL_STOP_FILE (se incluir separador de diretório)
//  2. Diretório definido em ICRMSENDEREMAIL_STOP_DIR (se existir)
//  3. Diretório do executável
//  4. Diretório de trabalho atual (working dir)
//
// Procura o primeiro arquivo encontrado com nome StopFileName (ou nome explícito em ICRMSENDEREMAIL_STOP_FILE).
func WatchStopFile(ctx context.Context, log *zap.Logger, cancel context.CancelFunc) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Resolver candidatos
	fileName, candidateDirs := resolveStopFileCandidates(log)

	// Montar lista final de caminhos completos a verificar
	var candidateFiles []string
	explicit := os.Getenv("ICRMSENDEREMAIL_STOP_FILE")
	if explicit != "" && isAbsOrHasDir(explicit) {
		candidateFiles = append(candidateFiles, explicit)
	}
	for _, d := range candidateDirs {
		candidateFiles = append(candidateFiles, filepath.Join(d, fileName))
	}

	log.Info("Monitoramento de arquivo de controle iniciado",
		zap.Strings("caminhos", candidateFiles))
	log.Info("Para parar via arquivo, crie um dos caminhos acima. Exemplo rápido:")
	log.Info("  echo stop > " + candidateFiles[0])

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, p := range candidateFiles {
				if _, err := os.Stat(p); err == nil { // encontrado
					log.Info("Arquivo de parada detectado, iniciando encerramento...", zap.String("arquivo", p))
					if err := os.Remove(p); err != nil {
						log.Warn("Não foi possível remover arquivo de controle (prosseguindo mesmo assim)", zap.Error(err))
					}
					cancel()
					return
				}
			}
		}
	}
}

// CreateStopFile cria o arquivo de controle para parar a aplicação
// CreateStopFile cria o arquivo de controle tentando respeitar variáveis de ambiente.
// Ordem de tentativa:
//  1. ICRMSENDEREMAIL_STOP_FILE (se absoluto ou contiver separador de diretório)
//  2. ICRMSENDEREMAIL_STOP_DIR (se existir)
//  3. Diretório do executável
//  4. Diretório de trabalho atual
func CreateStopFile() error {
	fileName, candidateDirs := resolveStopFileCandidates(nil)

	// Se usuário especificou caminho completo
	if explicit := os.Getenv("ICRMSENDEREMAIL_STOP_FILE"); explicit != "" && isAbsOrHasDir(explicit) {
		return writeStopFile(explicit)
	}

	// Tentar em ordem
	for _, d := range candidateDirs {
		path := filepath.Join(d, fileName)
		if err := writeStopFile(path); err == nil {
			return nil
		}
	}
	return errors.New("falha ao criar arquivo de stop em todos os diretórios candidatos")
}

// resolveStopFileCandidates retorna o nome do arquivo (pode ser sobrescrito por ICRMSENDEREMAIL_STOP_FILE sem path)
// e os diretórios candidatos em ordem de prioridade, sem duplicados.
func resolveStopFileCandidates(log *zap.Logger) (string, []string) {
	fileName := StopFileName
	if explicit := os.Getenv("ICRMSENDEREMAIL_STOP_FILE"); explicit != "" && !isAbsOrHasDir(explicit) {
		fileName = explicit
	}

	// Coletar diretórios
	dirsMap := make(map[string]struct{})
	add := func(p string) {
		if p == "" {
			return
		}
		if _, ok := dirsMap[p]; !ok {
			dirsMap[p] = struct{}{}
		}
	}

	if d := os.Getenv("ICRMSENDEREMAIL_STOP_DIR"); d != "" {
		if st, err := os.Stat(d); err == nil && st.IsDir() {
			add(d)
		} else if log != nil && err == nil && !st.IsDir() {
			log.Warn("ICRMSENDEREMAIL_STOP_DIR não é diretório, ignorando", zap.String("valor", d))
		}
	}
	if execPath, err := os.Executable(); err == nil {
		add(filepath.Dir(execPath))
	}
	if wd, err := os.Getwd(); err == nil {
		add(wd)
	}

	// Ordenar para estabilidade
	var dirs []string
	for d := range dirsMap {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)
	return fileName, dirs
}

func isAbsOrHasDir(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}
	// Contém separador? então usuário forneceu diretório relativo + nome
	return filepath.Dir(path) != "."
}

func writeStopFile(p string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString("stop"); err != nil {
		return err
	}
	return nil
}
