
CREATE DATABASE Kliptopia;

-- Drop table

DROP TABLE public.users;

CREATE TABLE public.users (
	user_id serial NOT NULL,
	username varchar(50) NOT NULL,
	email varchar(100) NOT NULL,
	password_hash varchar(100) NOT NULL,
	first_name varchar(50) NULL,
	last_name varchar(50) NULL,
	profile_picture_url varchar(255) NULL,
	registration_date timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	last_login_date timestamp NULL,
	"role" varchar(20) NULL DEFAULT 'user'::character varying,
	account_status varchar(20) NULL DEFAULT 'active'::character varying,
	verification_status varchar(20) NULL DEFAULT 'unverified'::character varying,
	password_reset_token uuid NULL,
	password_reset_expiry timestamp NULL,
	two_factor_enabled bool NULL DEFAULT false,
	CONSTRAINT users_email_key UNIQUE (email)
);
