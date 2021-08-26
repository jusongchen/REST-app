
BEGIN;

CREATE TABLE task ( 
     task_id VARCHAR(50) NOT NULL PRIMARY KEY,
     status VARCHAR(20) NOT NULL,			
     submitted_at TIMESTAMPTZ NOT NULL 
     
 );


CREATE TABLE Lock (
	lock_id VARCHAR(100) PRIMARY KEY,
	expires TIMESTAMPTZ NOT NULL
);

ALTER TABLE Lock SET (autovacuum_enabled = true);

CREATE OR REPLACE FUNCTION AcquireLock(VARCHAR(100), INT) RETURNS TIMESTAMP AS $$
	DECLARE
		nowT TIMESTAMP;
		expiresT TIMESTAMP;
	BEGIN
		nowT := CURRENT_TIMESTAMP;
		expiresT := nowT + '1 SECOND'::interval * $2;

		IF EXISTS (SELECT lock_id FROM Lock WHERE lock_id = $1 AND expires > nowT) THEN
			RETURN to_timestamp(0); -- Special value indicating no lock acquired.
		END IF;

		IF EXISTS (SELECT lock_id FROM Lock WHERE lock_id = $1) THEN
			UPDATE Lock SET expires = expiresT WHERE lock_id = $1;
			RETURN expiresT;
		END IF;

		INSERT INTO Lock (lock_id, expires) VALUES ($1, expiresT);
		RETURN expiresT;
	END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ReleaseLock(VARCHAR(100), TIMESTAMP) RETURNS BOOLEAN AS $$
	BEGIN
		IF NOT EXISTS (SELECT lock_id FROM Lock WHERE lock_id = $1 AND expires = $2) THEN
			RETURN FALSE; -- Another process acquired an expired lock
		END IF;

		DELETE FROM Lock WHERE lock_id = $1;
		RETURN TRUE;
	END
$$ LANGUAGE plpgsql;

END;
