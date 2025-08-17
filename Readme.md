CREATE TABLE public.posts (
id serial4 NOT NULL,
title varchar(255) NOT NULL,
body text NOT NULL,
user_id int4 NOT NULL,
CONSTRAINT posts_pkey PRIMARY KEY (id)
);

CREATE TABLE public.users (
id serial4 NOT NULL,
username varchar(50) NOT NULL,
email varchar(100) NOT NULL,
"password" text NOT NULL,
is_active bool DEFAULT true NULL,
created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
CONSTRAINT users_email_key UNIQUE (email),
CONSTRAINT users_pkey PRIMARY KEY (id),
CONSTRAINT users_username_key UNIQUE (username)
);
