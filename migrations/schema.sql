--
-- PostgreSQL database dump
--

-- Dumped from database version 10.4
-- Dumped by pg_dump version 10.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: imports; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.imports (
    id uuid NOT NULL,
    author character varying(50) NOT NULL,
    source character varying(100) NOT NULL,
    processed boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.imports OWNER TO postgres;

--
-- Name: pictures; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pictures (
    id uuid NOT NULL,
    upvotes integer NOT NULL,
    downvotes integer NOT NULL,
    sorting double precision NOT NULL,
    confidence_level double precision NOT NULL,
    filename character varying(100) NOT NULL,
    hash character varying(100) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    import_id uuid NOT NULL
);


ALTER TABLE public.pictures OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(255) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.votes (
    id uuid NOT NULL,
    picture_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.votes OWNER TO postgres;

--
-- Name: imports imports_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.imports
    ADD CONSTRAINT imports_pkey PRIMARY KEY (id);


--
-- Name: pictures pictures_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pictures
    ADD CONSTRAINT pictures_pkey PRIMARY KEY (id);


--
-- Name: votes votes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pkey PRIMARY KEY (id);


--
-- Name: pictures_confidence_level_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX pictures_confidence_level_idx ON public.pictures USING btree (confidence_level);


--
-- Name: pictures_sorting_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX pictures_sorting_idx ON public.pictures USING btree (sorting);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: pictures pictures_imports_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pictures
    ADD CONSTRAINT pictures_imports_id_fk FOREIGN KEY (import_id) REFERENCES public.imports(id);


--
-- Name: votes votes_picture_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_picture_id_fkey FOREIGN KEY (picture_id) REFERENCES public.pictures(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

