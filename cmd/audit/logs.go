package audit

import (
	"fmt"
	"nautic/cmd/storage"
	"strings"
	"time"
)

func InsertLog(id_user int, url string, action string, extra_description string, sql string) error {
	db := storage.GetDB()

	year := time.Now().Year()
	logTable := fmt.Sprintf("logs_%d", year)

	query_table := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS public.logs_%d
		(
			id bigint NOT NULL DEFAULT nextval('logs_%d_id_seq'::regclass),
			url character varying(100) COLLATE pg_catalog."default",
			action character varying(100) COLLATE pg_catalog."default",
			extra_description text COLLATE pg_catalog."default",
			sql text COLLATE pg_catalog."default",
			id_user bigint NOT NULL,
			created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT logs_%d_pkey PRIMARY KEY (id),
			CONSTRAINT users_fkey FOREIGN KEY (id_user)
				REFERENCES public.users (id) MATCH FULL
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
				NOT VALID
		)

		TABLESPACE pg_default;

		ALTER TABLE IF EXISTS public.logs_%d
			OWNER to postgres;
	`, year, year, year, year)

	_, err := db.Exec(query_table)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (id_user, url, action, extra_description, sql) VALUES ($1, $2, $3, $4, $5)", logTable)

	_, err = db.Exec(query, id_user, url, action, extra_description, strings.TrimSpace(sql))
	if err != nil {
		return err
	}

	return nil
}
