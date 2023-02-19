CREATE TABLE IF NOT EXISTS users(
	id VARCHAR (30) PRIMARY KEY,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,

	name VARCHAR (100) NOT NULL,
	username VARCHAR (100) NOT NULL UNIQUE,
	role VARCHAR (50) NOT NULL,
	active BOOLEAN DEFAULT true
);

CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(active);

INSERT INTO users(id, name, username, role, active) VALUES('0', 'Anonymous', 'anonymous', 'anonymous', false);

CREATE TABLE IF NOT EXISTS identities(
	id         VARCHAR (30) PRIMARY KEY,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,

	last_login TIMESTAMP,
	provider VARCHAR (30), -- phone, email, wechat, github...
	uid VARCHAR (30), 		 -- e-mail, google id, facebook id, etc
	password VARCHAR (300),
	verified BOOLEAN DEFAULT false,
	confirmed_at TIMESTAMP,

	user_id VARCHAR (30) NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS links(
	id VARCHAR (30) PRIMARY KEY,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,

	domain VARCHAR (100) NOT NULL,
	keyword VARCHAR (100) NOT NULL,
	url VARCHAR (300) NOT NULL,
	title VARCHAR (100),
	active BOOLEAN DEFAULT true,

	user_id VARCHAR (30) NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX idx_links_domain ON links (domain);
CREATE UNIQUE INDEX idx_links_keyword ON links (keyword);

CREATE TABLE IF NOT EXISTS groups(
	id VARCHAR (30) PRIMARY KEY,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,

	name VARCHAR (100) NOT NULL,
	display_name VARCHAR (100),
	description VARCHAR (300),

	created_by VARCHAR (30) NOT NULL,
	owner_id VARCHAR (30) NOT NULL,

	CONSTRAINT fk_created_by FOREIGN KEY(created_by) REFERENCES users(id),
	CONSTRAINT fk_owner_id FOREIGN KEY(owner_id) REFERENCES users(id)
);

CREATE INDEX idx_groups_created_by ON groups (created_by);
CREATE INDEX idx_groups_owner_id ON groups (owner_id);

CREATE TABLE group_user(
	group_id VARCHAR (30) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	user_id VARCHAR (30) REFERENCES users(id) ON UPDATE CASCADE
);

CREATE INDEX idx_group_user_group_id ON group_user(group_id);
CREATE INDEX idx_group_user_user_id ON group_user(user_id);

CREATE TABLE IF NOT EXISTS groups_invites(
	id VARCHAR (30) PRIMARY KEY,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP,

	group_id VARCHAR (30) NOT NULL,
	user_id VARCHAR (30) NOT NULL,
	invited_by VARCHAR (30) NOT NULL, -- invited by user_id

	accepted BOOLEAN DEFAULT false,

	CONSTRAINT fk_group_id FOREIGN KEY(group_id) REFERENCES groups(id),
	CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id),
	CONSTRAINT fk_invited_by FOREIGN KEY(invited_by) REFERENCES users(id)
);
