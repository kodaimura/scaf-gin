CREATE TABLE IF NOT EXISTS account (
	account_id SERIAL PRIMARY KEY,
	account_name TEXT NOT NULL UNIQUE,
	account_password TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create trigger trg_account_upd BEFORE UPDATE ON account FOR EACH ROW
  execute procedure set_updated_at();