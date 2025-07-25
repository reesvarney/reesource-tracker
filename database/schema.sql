CREATE TABLE
    IF NOT EXISTS locations (
        id BLOB(16) PRIMARY KEY NOT NULL,
        name VARCHAR(128) NOT NULL,
        description TEXT,
        parent_location_id BLOB(16) REFERENCES locations (id)
    );

CREATE TABLE
    IF NOT EXISTS products (
        id BLOB(16) PRIMARY KEY NOT NULL,
        name VARCHAR(128) NOT NULL,
        parent_product_id BLOB(16) REFERENCES products (id)
    );

CREATE TABLE
    IF NOT EXISTS samples (
        id BLOB(4) PRIMARY KEY NOT NULL,
        location_id BLOB(16) REFERENCES locations (id),
        product_id BLOB(16) REFERENCES products (id),
        time_registered DATETIME,
        last_update DATETIME,
        state TEXT CHECK (
            state IN (
                "in_use",
                "broken",
                "available",
                "archived",
                "unassigned"
            )
        ) DEFAULT 'unassigned' NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS sample_mods (
        id BLOB(16) PRIMARY KEY NOT NULL,
        sample_id BLOB(4) NOT NULL REFERENCES samples (id),
        name VARCHAR(128) NOT NULL,
        time_added DATETIME NOT NULL,
        time_removed DATETIME
    );

CREATE TABLE
    IF NOT EXISTS sample_notes (
        id BLOB(16) PRIMARY KEY NOT NULL,
        sample_id BLOB(4) NOT NULL REFERENCES samples (id),
        contents TEXT NOT NULL,
        time_made DATETIME NOT NULL
    );