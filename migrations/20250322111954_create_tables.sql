-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		username VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS folders(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		user_id BIGSERIAL NOT NULL,
		name VARCHAR(100) NOT NULL,
		UNIQUE (user_id, name),
		FOREIGN KEY (user_id) REFERENCES users (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TYPE note_type AS ENUM ('message', 'photo', 'audio', 'doc', 'video', 'animation', 'voice');

CREATE TABLE IF NOT EXISTS texts(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		user_id BIGSERIAL NOT NULL,
		folder_id BIGSERIAL NOT NULL,
		type note_type NOT NULL DEFAULT 'message',
		description TEXT,
		media_group_id VARCHAR(100) NOT NULL DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
		ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (folder_id) REFERENCES folders (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS photos(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		texts_id BIGSERIAL NOT NULL,
		file_id VARCHAR(100),
		media_group_id VARCHAR(100) NOT NULL DEFAULT '',
		FOREIGN KEY (texts_id) REFERENCES texts (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS audios(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		texts_id BIGSERIAL NOT NULL,
		file_id VARCHAR(100),
		media_group_id VARCHAR(100) NOT NULL DEFAULT '',
		FOREIGN KEY (texts_id) REFERENCES texts (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS documents(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		texts_id BIGSERIAL NOT NULL,
		file_id VARCHAR(100),
		media_group_id VARCHAR(100) NOT NULL DEFAULT '',
		FOREIGN KEY (texts_id) REFERENCES texts (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS videos(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		texts_id BIGSERIAL NOT NULL,
		file_id VARCHAR(100),
		media_group_id VARCHAR(100) NOT NULL DEFAULT '',
		FOREIGN KEY (texts_id) REFERENCES texts (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS animations(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		texts_id BIGSERIAL NOT NULL,
		file_id VARCHAR(100),
		media_group_id VARCHAR(100) NOT NULL DEFAULT '',
		FOREIGN KEY (texts_id) REFERENCES texts (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS voices(
		id BIGSERIAL NOT NULL PRIMARY KEY,
		texts_id BIGSERIAL NOT NULL,
		file_id VARCHAR(100),
		FOREIGN KEY (texts_id) REFERENCES texts (id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS 
    users,
    texts,
    photos,
    audios,
    documents,
    videos,
    animations,
    voices;
    
DROP TYPE IF EXISTS note_type;
-- +goose StatementEnd
