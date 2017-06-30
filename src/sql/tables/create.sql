-- Table: public.acl_group

-- DROP TABLE public.acl_group;

CREATE TABLE public.acl_group
(
  id SERIAL,
  name character varying(20),
  parent_id integer REFERENCES acl_group ON DELETE CASCADE,
  CONSTRAINT acl_group_id_c PRIMARY KEY (id),
  CONSTRAINT acl_group_name_key UNIQUE (name)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.acl_group
  OWNER TO webitel;




-- Table: public.acl_permission

-- DROP TABLE public.acl_permission;

CREATE TABLE public.acl_permission
(
  id SERIAL,
  object_type character varying(20),
  object_id character varying(20),
  group_id integer NOT NULL,
  perm varchar(2)[],
  CONSTRAINT acl_permission_pkey PRIMARY KEY (group_id, id),
  CONSTRAINT acl_permission_group_id_fkey FOREIGN KEY (group_id)
      REFERENCES public.acl_group (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.acl_permission
  OWNER TO webitel;


-- Table: public.acl_tokens

-- DROP TABLE public.acl_tokens;

CREATE TABLE public.acl_tokens
(
  key uuid NOT NULL,
  domain character varying(70),
  username character varying(70),
  expires TIMESTAMP,
  group_id integer,
  enabled BIT,
  createdBy varchar(70),
  createdOn TIMESTAMP  DEFAULT 'NOW()',
  CONSTRAINT acl_tokens_group_id_fkey FOREIGN KEY (group_id)
      REFERENCES public.acl_group (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.acl_tokens
  OWNER TO webitel;



