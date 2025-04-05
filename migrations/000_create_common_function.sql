create function set_updated_at() returns trigger AS '
  	BEGIN
    	new.updated_at := NOW();
    	return new;
  	END;
' language 'plpgsql';