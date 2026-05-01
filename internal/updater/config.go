package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

// configFileName é o nome do arquivo dentro de UserConfigDir/cfs_spool/.
const configFileName = "updater.json"

// ConfigStore persiste preferências do verificador de updates (versões
// que o usuário escolheu ignorar) em um arquivo JSON simples.
//
// Path pode ser injetado em testes — em produção, DefaultConfigPath()
// resolve para o diretório de config do SO.
type ConfigStore struct {
	Path string
}

// configFile é o schema serializado em JSON. Mantido pequeno e versionado
// implicitamente: campos novos são lidos como zero-value se ausentes.
type configFile struct {
	IgnoredVersions []string `json:"ignoredVersions"`
}

// DefaultConfigPath devolve <UserConfigDir>/cfs_spool/updater.json.
// Em macOS: ~/Library/Application Support/cfs_spool/updater.json.
// Em Linux: ~/.config/cfs_spool/updater.json. Em Windows: %APPDATA%\cfs_spool\updater.json.
func DefaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("updater: descobrir config dir: %w", err)
	}
	return filepath.Join(dir, "cfs_spool", configFileName), nil
}

// NewConfigStore cria um store apontando para DefaultConfigPath. Se o
// caminho não puder ser resolvido, devolve um store com Path vazio — todas
// as operações degradam silenciosamente (IsIgnored sempre false, IgnoreVersion
// retorna erro). Não fazemos fallback para diretório arbitrário para evitar
// gravar fora do espaço de config esperado.
func NewConfigStore() *ConfigStore {
	path, _ := DefaultConfigPath()
	return &ConfigStore{Path: path}
}

// IsIgnored devolve true se a versão foi marcada como "não notificar".
// Erros de leitura (arquivo ausente, JSON corrompido, permissão) são tratados
// como "nada ignorado" — o usuário recebe a notificação normalmente.
func (c *ConfigStore) IsIgnored(version string) bool {
	if c == nil || c.Path == "" || version == "" {
		return false
	}
	cfg, err := c.load()
	if err != nil {
		return false
	}
	return slices.Contains(cfg.IgnoredVersions, version)
}

// IgnoreVersion adiciona a versão à lista de ignoradas e persiste o arquivo.
// Idempotente: chamar duas vezes para a mesma versão é no-op após a primeira.
func (c *ConfigStore) IgnoreVersion(version string) error {
	if c == nil || c.Path == "" {
		return errors.New("updater: config path não configurado")
	}
	if version == "" {
		return errors.New("updater: versão vazia")
	}

	cfg, err := c.load()
	if err != nil {
		// Arquivo corrompido / ausente: começamos uma config nova.
		cfg = &configFile{}
	}
	if slices.Contains(cfg.IgnoredVersions, version) {
		return nil
	}
	cfg.IgnoredVersions = append(cfg.IgnoredVersions, version)

	return c.save(cfg)
}

// load lê o arquivo JSON. Devolve erro se não existir ou estiver corrompido —
// IsIgnored e IgnoreVersion convertem isso em "config vazia" para resiliência.
func (c *ConfigStore) load() (*configFile, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &configFile{}, nil
		}
		return nil, err
	}
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// save serializa cfg em JSON indentado e escreve atomicamente (write-rename)
// para que uma falha no meio da escrita não corrompa o arquivo existente.
func (c *ConfigStore) save(cfg *configFile) error {
	if err := os.MkdirAll(filepath.Dir(c.Path), 0o755); err != nil {
		return fmt.Errorf("updater: criar config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("updater: serializar config: %w", err)
	}

	tmp := c.Path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("updater: escrever config tmp: %w", err)
	}
	if err := os.Rename(tmp, c.Path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("updater: renomear config: %w", err)
	}
	return nil
}
