-- +goose Up
CREATE TABLE IF NOT EXISTS  events (
	id SERIAL, 
	title varchar(100) NOT NULL DEFAULT '', 
	start_time timestamp NOT NULL, 
	end_time timestamp NOT NULL,
	duration double precision NOT NULL, 
	description varchar(100) DEFAULT '', 
	owner_id varchar(36) NOT NULL DEFAULT '', 
	remind_time timestamp NOT NULL , 
	PRIMARY KEY (id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE events;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
