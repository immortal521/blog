SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE public.image_folders (
    id uuid NOT NULL,
    parent_id uuid,
    name text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE public.images (
    id uuid NOT NULL,
    folder_id uuid,
    storage_key text NOT NULL,
    origin_name text NOT NULL,
    mime text NOT NULL,
    size bigint NOT NULL,
    width bigint,
    height bigint,
    sha256 text,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE TABLE public.link_categories (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(20) NOT NULL,
    sort_order bigint DEFAULT 0 NOT NULL
);

CREATE SEQUENCE public.link_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.link_categories_id_seq OWNED BY public.link_categories.id;

CREATE TABLE public.links (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    description character varying(255),
    enabled boolean DEFAULT false NOT NULL,
    name character varying(100) NOT NULL,
    sort_order bigint DEFAULT 0 NOT NULL,
    url character varying(255) NOT NULL,
    avatar character varying(255),
    category_id bigint,
    status smallint DEFAULT 1 NOT NULL
);

CREATE SEQUENCE public.links_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.links_id_seq OWNED BY public.links.id;

CREATE TABLE public.post_categories (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL
);

CREATE SEQUENCE public.post_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.post_categories_id_seq OWNED BY public.post_categories.id;

CREATE TABLE public.post_category_relations (
    post_category_id bigint NOT NULL,
    post_id bigint NOT NULL
);

CREATE TABLE public.post_tag_relations (
    post_tag_id bigint NOT NULL,
    post_id bigint NOT NULL
);

CREATE TABLE public.post_tags (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL
);

CREATE SEQUENCE public.post_tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.post_tags_id_seq OWNED BY public.post_tags.id;

CREATE TABLE public.posts (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title character varying(255) NOT NULL,
    summary character varying(500),
    content text NOT NULL,
    cover character varying(255),
    read_time_minutes bigint NOT NULL,
    view_count bigint NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    user_id bigint NOT NULL,
    published_at timestamp with time zone
);

CREATE SEQUENCE public.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    uuid uuid NOT NULL,
    avatar character varying(255),
    email character varying(100) NOT NULL,
    password character varying(255) NOT NULL,
    role smallint DEFAULT 1 NOT NULL,
    username character varying(50) NOT NULL
);

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

ALTER TABLE ONLY public.link_categories ALTER COLUMN id SET DEFAULT nextval('public.link_categories_id_seq'::regclass);

ALTER TABLE ONLY public.links ALTER COLUMN id SET DEFAULT nextval('public.links_id_seq'::regclass);

ALTER TABLE ONLY public.post_categories ALTER COLUMN id SET DEFAULT nextval('public.post_categories_id_seq'::regclass);

ALTER TABLE ONLY public.post_tags ALTER COLUMN id SET DEFAULT nextval('public.post_tags_id_seq'::regclass);

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);

ALTER TABLE ONLY public.image_folders
    ADD CONSTRAINT image_folders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.link_categories
    ADD CONSTRAINT link_categories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.links
    ADD CONSTRAINT links_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.post_categories
    ADD CONSTRAINT post_categories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.post_category_relations
    ADD CONSTRAINT post_category_relations_pkey PRIMARY KEY (post_category_id, post_id);

ALTER TABLE ONLY public.post_tag_relations
    ADD CONSTRAINT post_tag_relations_pkey PRIMARY KEY (post_tag_id, post_id);

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT post_tags_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.link_categories
    ADD CONSTRAINT uni_link_categories_name UNIQUE (name);

ALTER TABLE ONLY public.links
    ADD CONSTRAINT uni_links_url UNIQUE (url);

ALTER TABLE ONLY public.post_categories
    ADD CONSTRAINT uni_post_categories_name UNIQUE (name);

ALTER TABLE ONLY public.post_categories
    ADD CONSTRAINT uni_post_categories_slug UNIQUE (slug);

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT uni_post_tags_name UNIQUE (name);

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT uni_post_tags_slug UNIQUE (slug);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_uuid UNIQUE (uuid);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE INDEX idx_image_folders_deleted_at ON public.image_folders USING btree (deleted_at);

CREATE INDEX idx_image_folders_parent_id ON public.image_folders USING btree (parent_id);

CREATE INDEX idx_images_deleted_at ON public.images USING btree (deleted_at);

CREATE INDEX idx_images_folder_id ON public.images USING btree (folder_id);

CREATE INDEX idx_images_sha256 ON public.images USING btree (sha256);

CREATE INDEX idx_images_storage_key ON public.images USING btree (storage_key);

CREATE INDEX idx_link_categories_deleted_at ON public.link_categories USING btree (deleted_at);

CREATE INDEX idx_links_deleted_at ON public.links USING btree (deleted_at);

CREATE INDEX idx_post_categories_deleted_at ON public.post_categories USING btree (deleted_at);

CREATE INDEX idx_post_tags_deleted_at ON public.post_tags USING btree (deleted_at);

CREATE INDEX idx_posts_deleted_at ON public.posts USING btree (deleted_at);

CREATE INDEX idx_posts_user_id ON public.posts USING btree (user_id);

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);
