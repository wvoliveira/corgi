CREATE TABLE IF NOT EXISTS links(
	id VARCHAR (30) PRIMARY KEY,
  user_id VARCHAR (30) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,

	domain VARCHAR (100) NOT NULL,
	keyword VARCHAR (100) NOT NULL,
	url VARCHAR (300) NOT NULL,
	title VARCHAR (100),
	active BOOLEAN DEFAULT true,

	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX idx_links_domain ON links (domain);
CREATE UNIQUE INDEX idx_links_keyword ON links (keyword);
