CREATE TABLE buildings (
    id SERIAL PRIMARY KEY,
    doitt_id TEXT NOT NULL,
    bin TEXT NOT NULL,
    construct_year integer NOT NULL,
    height_roof double precision NOT NULL DEFAULT 0.0,
    area double precision NOT NULL DEFAULT 0.0
  );

-- index
CREATE UNIQUE INDEX building_pkey ON buildings(id int4_ops);
CREATE INDEX building_year ON buildings(construct_year int4_ops);
