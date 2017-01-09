CREATE TABLE IF NOT EXISTS owl_query (
	qr_uuid	BINARY(16) PRIMARY KEY,
	qr_named_id	VARCHAR(32),
	qr_md5_content BINARY(16),
	qr_content VARBINARY(20480) NOT NULL,
	qr_time_creation DATETIME NOT NULL,
	qr_time_access DATETIME NOT NULL,
	CONSTRAINT unq_owl_query__qr_named_id_qr_md5_content UNIQUE(
		qr_named_id, qr_md5_content
	),
	INDEX ix_owl_query__qr_time_access (qr_time_access ASC)
);
