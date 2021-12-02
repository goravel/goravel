package migrations

type MigrateStubs struct {
}

func (receiver MigrateStubs) CreateUp() string {
	return `CREATE TABLE DummyTable (
  id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_DummyTable_created_at (created_at),
  KEY idx_DummyTable_updated_at (updated_at)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = DummyDatabaseCharset;
`
}

func (receiver MigrateStubs) CreateDown() string {
	return `DROP TABLE IF EXISTS DummyTable;`
}

func (receiver MigrateStubs) UpdateUp() string {
	return `ALTER TABLE DummyTable ADD column varchar(255) COMMENT '';`
}

func (receiver MigrateStubs) UpdateDown() string {
	return `ALTER TABLE DummyTable DROP COLUMN column;`
}
