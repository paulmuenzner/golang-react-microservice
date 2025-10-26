
CREATE SCHEMA IF NOT EXISTS templates;

CREATE TABLE IF NOT EXISTS templates.profile_objects (
    -- Eindeutige ID für den Datensatz
    id BIGSERIAL PRIMARY KEY, 
    
    -- NEUE PFLICHTFELDER FÜR PROFILGRUPPIERUNG
    
    -- ID des Lastprofils (maximal zweistellig, z.B. 1-99). Pflichtfeld.
    -- Damit man zuordnen kann ob LP1, LP2, ... gemeint ist.
    load_profile SMALLINT NOT NULL,
    -- Die Versionsnummer des Profils (ganzzahlig, aufsteigend). Pflichtfeld.
    -- Falls es eine Änderung des spezifischen Objektes gibt, wird die Versionsnummer erhöht. Alte Objekte werden nicht gelöscht, da sie noch von gespeicherten Templates genutzt werden können und dann fehlen würden
    version INT NOT NULL,
    -- Eindeutige Kennung des Objekts (ganzzahlig). Pflichtfeld.
    -- Jeder Objekttyp muss eine eindeutige Kennung haben. Sie ist nur identisch bei Objekten wo sich die Versionsnummer erhöht hat. So kann man via objekt_kennung alle versionen abfragen.
    objekt_kennung INT NOT NULL,
    
    -- ALLGEMEIN GÜLTIGE FELDER (enthalten in DataItem1 und DataItem2)
    
    -- Interner technischer Name
    name VARCHAR(255) NOT NULL, 
    -- Deutscher Name
    name_de VARCHAR(255) NOT NULL, 
    -- Englischer Name
    name_en VARCHAR(255) NOT NULL, 
    
    -- Identifikationsattribute des Objekts
    obis_code VARCHAR(50) NOT NULL,
    class_id VARCHAR(50) NOT NULL,
    attribute_id VARCHAR(50) NOT NULL,
    
    -- Anzeige- und Verarbeitungsparameter
    column_order INT NOT NULL, 
    type VARCHAR(50) NOT NULL, 
    data_type VARCHAR(50) NOT NULL,
    
    -- Optionales Feld (lpName ist bereits optional in TS)
    lp_name VARCHAR(255) NULL, 
    
    -- SPEZIFISCHE FELDER VON DataItem1 (optional/NULL)
    
    -- Maßeinheit, z.B. 'kWh', 'A'
    unit VARCHAR(50) NULL, 
    -- Skalierungsfaktor (Potenz von 10)
    scaler INT NULL, 
    -- IEC-Code oder anderer Branchenstandardcode
    iec_code VARCHAR(50) NULL
);

-- Indices zur Beschleunigung der Suche nach Identifikatoren

-- Index für die Lastprofilsuche
CREATE INDEX idx_profile_objects_load_profile ON templates.profile_objects (load_profile);