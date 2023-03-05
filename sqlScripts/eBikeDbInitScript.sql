-- Create User Table

-- Table: public.users

DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    username character varying(32) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;





-- Table: public.reservation

DROP TABLE IF EXISTS public.reservation;

CREATE TABLE IF NOT EXISTS public.reservation
(
    reservationid uuid NOT NULL,
    bikeid integer NOT NULL,
    username character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT reservation_pkey PRIMARY KEY (reservationid),
    CONSTRAINT "username_Unique_contraint" UNIQUE (username),
    CONSTRAINT username_foreign_key FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.reservation
    OWNER to postgres;
-- Index: fki_username_foreign_key

DROP INDEX IF EXISTS public.fki_username_foreign_key;

CREATE INDEX IF NOT EXISTS fki_username_foreign_key
    ON public.reservation USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: reservation_FKEY_user_for_username

DROP INDEX IF EXISTS public."reservation_FKEY_user_for_username";

CREATE INDEX IF NOT EXISTS "reservation_FKEY_user_for_username"
    ON public.reservation USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;




-- Table: public.bike

DROP TABLE IF EXISTS public.bike;

CREATE TABLE IF NOT EXISTS public.bike
(
    bikeid integer NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    reservationid uuid,
    CONSTRAINT "Bikes_pkey" PRIMARY KEY (bikeid),
    CONSTRAINT "bike_reservationId_fkey" FOREIGN KEY (reservationid)
        REFERENCES public.reservation (reservationid) MATCH SIMPLE
        ON UPDATE NO ACTION -- cascade statt no action
        ON DELETE SET NULL -- set null statt NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.bike
    OWNER to postgres;
-- Index: fKey_from_bike_to_reservation_for_reservationId

DROP INDEX IF EXISTS public."fKey_from_bike_to_reservation_for_reservationId";

CREATE INDEX IF NOT EXISTS "fKey_from_bike_to_reservation_for_reservationId"
    ON public.bike USING btree
    (reservationid ASC NULLS LAST)
    TABLESPACE pg_default;


-- Insert Data into bike Table

INSERT INTO public.bike(
	bikeid, name, latitude, longitude, reservationid)
	VALUES (0, 'Henry', 50.119504, 8.638137, NULL);

INSERT INTO public.bike(
	bikeid, name, latitude, longitude, reservationid)
	VALUES (1, 'Hans', 50.119229, 8.640020, NULL);

INSERT INTO public.bike(
	bikeid, name, latitude, longitude, reservationid)
	VALUES (2, 'Thomas', 50.120452, 8.650507, NULL);

INSERT INTO public.bike(
	bikeid, name, latitude, longitude, reservationid)
	VALUES (3, 'Kevin', 50.55, 8.88, NULL);

-- Insert Data into reservation Table

INSERT INTO public.users(
	username)
	VALUES ('userOne');

INSERT INTO public.users(
	username)
	VALUES ('userTwo');