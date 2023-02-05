CREATE TABLE IF NOT EXISTS identities(
	id VARCHAR (30) PRIMARY KEY,
	user_id VARCHAR (30) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,
	last_login TIMESTAMP,
	provider VARCHAR (30), -- phone, email, wechat, github...
	uid      VARCHAR (30), -- e-mail, google id, facebook id, etc
	password VARCHAR (300),
	verified BOOLEAN DEFAULT false,
	confirmed_at TIMESTAMP,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);
