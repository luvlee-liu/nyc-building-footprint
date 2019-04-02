COPY buildings(doitt_id, bin, construct_year, height_roof, area)
FROM :building_csv DELIMITER ',' CSV HEADER;
