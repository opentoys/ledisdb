package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	ErrNoConfigFile = errors.New("Running without a config file")
)

const (
	DefaultDBName  string = "goleveldb"
	DefaultDataDir string = "./data"
	KB             int    = 1024
	MB             int    = KB * 1024
	GB             int    = MB * 1024
)

type LevelDBConfig struct {
	Compression     bool `json:"compression"`
	BlockSize       int  `json:"block_size"`
	WriteBufferSize int  `json:"write_buffer_size"`
	CacheSize       int  `json:"cache_size"`
	MaxOpenFiles    int  `json:"max_open_files"`
	MaxFileSize     int  `json:"max_file_size"`
}

type LMDBConfig struct {
	MapSize int  `json:"map_size"`
	NoSync  bool `json:"nosync"`
}

type ReplicationConfig struct {
	Path             string `json:"path"`
	Sync             bool   `json:"sync"`
	WaitSyncTime     int    `json:"wait_sync_time"`
	WaitMaxSlaveAcks int    `json:"wait_max_slave_acks"`
	ExpiredLogDays   int    `json:"expired_log_days"`
	StoreName        string `json:"store_name"`
	MaxLogFileSize   int64  `json:"max_log_file_size"`
	MaxLogFileNum    int    `json:"max_log_file_num"`
	SyncLog          int    `json:"sync_log"`
	Compression      bool   `json:"compression"`
	UseMmap          bool   `json:"use_mmap"`
	MasterPassword   string `json:"master_password"`
}

type SnapshotConfig struct {
	Path   string `json:"path"`
	MaxNum int    `json:"max_num"`
}

type Config struct {
	m                *sync.RWMutex     `json:"-"`
	FileName         string            `json:"-"`
	Readonly         bool              `json:"readonly"`
	DataDir          string            `json:"data_dir"`
	Databases        int               `json:"databases"`
	DBName           string            `json:"db_name"`
	DBPath           string            `json:"db_path"`
	DBSyncCommit     int               `json:"db_sync_commit"`
	LevelDB          LevelDBConfig     `json:"leveldb"`
	LMDB             LMDBConfig        `json:"lmdb"`
	AccessLog        string            `json:"access_log"`
	UseReplication   bool              `json:"use_replication"`
	Replication      ReplicationConfig `json:"replication"`
	Snapshot         SnapshotConfig    `json:"snapshot"`
	TTLCheckInterval int               `json:"ttl_check_interval"`
}

func NewConfigWithFile(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfigWithData(data)
	if err != nil {
		return nil, err
	}

	cfg.FileName = fileName
	return cfg, nil
}

func NewConfigWithData(data []byte) (*Config, error) {
	cfg := NewConfigDefault()

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("newConfigwithData: unmarashal: %s", err)
	}

	cfg.adjust()

	return cfg, nil
}

func NewConfigDefault() *Config {
	cfg := new(Config)
	cfg.m = new(sync.RWMutex)

	cfg.DataDir = DefaultDataDir

	cfg.DBName = DefaultDBName
	cfg.Readonly = false
	// default databases number
	cfg.Databases = 16

	// disable access log
	cfg.AccessLog = ""

	cfg.LMDB.MapSize = 20 * MB
	cfg.LMDB.NoSync = true

	cfg.UseReplication = false
	cfg.Replication.WaitSyncTime = 500
	cfg.Replication.Compression = true
	cfg.Replication.WaitMaxSlaveAcks = 2
	cfg.Replication.SyncLog = 0
	cfg.Replication.UseMmap = true
	cfg.Snapshot.MaxNum = 1

	cfg.adjust()

	return cfg
}

func getDefault(d int, s int) int {
	if s <= 0 {
		return d
	}

	return s
}

func (cfg *Config) adjust() {
	cfg.LevelDB.adjust()

	cfg.Replication.ExpiredLogDays = getDefault(7, cfg.Replication.ExpiredLogDays)
	cfg.Replication.MaxLogFileNum = getDefault(50, cfg.Replication.MaxLogFileNum)
	cfg.Databases = getDefault(16, cfg.Databases)
}

func (cfg *LevelDBConfig) adjust() {
	cfg.CacheSize = getDefault(4*MB, cfg.CacheSize)
	cfg.BlockSize = getDefault(4*KB, cfg.BlockSize)
	cfg.WriteBufferSize = getDefault(4*MB, cfg.WriteBufferSize)
	cfg.MaxOpenFiles = getDefault(1024, cfg.MaxOpenFiles)
	cfg.MaxFileSize = getDefault(32*MB, cfg.MaxFileSize)
}

func (cfg *Config) GetReadonly() bool {
	cfg.m.RLock()
	b := cfg.Readonly
	cfg.m.RUnlock()
	return b
}

func (cfg *Config) SetReadonly(b bool) {
	cfg.m.Lock()
	cfg.Readonly = b
	cfg.m.Unlock()
}
