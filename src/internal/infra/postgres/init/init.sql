CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- User
CREATE TABLE IF NOT EXISTS Users (
    id UUID PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL,
    sex BOOLEAN NOT NULL,
    weight INT NOT NULL,
    age INT NOT NULL,
    is_pregnant BOOLEAN NOT NULL,
    is_driver BOOLEAN NOT NULL
);

-- Medicine
CREATE TABLE IF NOT EXISTS Medicine (
    id UUID PRIMARY KEY,
    form_id UUID NOT NULL,
    unit_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    expire_time int NOT NULL,

    effect_on_driver BOOLEAN NOT NULL,
    effect_on_pregnant BOOLEAN NOT NULL,
    method_of_application VARCHAR(255) NOT NULL,
    is_prescription BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS Form (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Unit (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

ALTER TABLE Medicine
ADD CONSTRAINT fk_unit
FOREIGN KEY (unit_id) REFERENCES Unit(id);

ALTER TABLE Medicine
ADD CONSTRAINT fk_form
FOREIGN KEY (form_id) REFERENCES Form(id);

-- Dosage
CREATE TABLE IF NOT EXISTS Dosage (
    id UUID PRIMARY KEY,
    medicine_id UUID NOT NULL,
    value_from INT NOT NULL,
    value_to INT NOT NULL,
    dosage_type VARCHAR(255) NOT NULL,
    dosage_value FLOAT NOT NULL,
    number_of_doses_per_day INT NOT NULL
);

ALTER TABLE Dosage
ADD CONSTRAINT fk_medicine
FOREIGN KEY (medicine_id) REFERENCES Medicine(id);

-- Substance
CREATE TABLE IF NOT EXISTS Substance (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Illness
CREATE TABLE IF NOT EXISTS Illness (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Medicine_Illness
CREATE TABLE IF NOT EXISTS Recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medicine_id UUID NOT NULL,
    illness_id UUID NOT NULL
);

ALTER TABLE Recommendations
ADD CONSTRAINT fk_medicine_recommendations
FOREIGN KEY (medicine_id) REFERENCES Medicine(id);

ALTER TABLE Recommendations
ADD CONSTRAINT fk_illness_recommendations
FOREIGN KEY (illness_id) REFERENCES Illness(id);

CREATE TABLE IF NOT EXISTS Contraindications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medicine_id UUID NOT NULL,
    illness_id UUID NOT NULL
);

ALTER TABLE Contraindications
ADD CONSTRAINT fk_medicine_contraindications
FOREIGN KEY (medicine_id) REFERENCES Medicine(id);

ALTER TABLE Contraindications
ADD CONSTRAINT fk_illness_contraindications
FOREIGN KEY (illness_id) REFERENCES Illness(id);

-- Medicine_Substance
CREATE TABLE IF NOT EXISTS Medicine_Substance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medicine_id UUID NOT NULL,
    substance_id UUID NOT NULL,
    concentration FLOAT NOT NULL
);

ALTER TABLE Medicine_Substance
ADD CONSTRAINT fk_medicine_substance
FOREIGN KEY (medicine_id) REFERENCES Medicine(id);

ALTER TABLE Medicine_Substance
ADD CONSTRAINT fk_substance_medicine
FOREIGN KEY (substance_id) REFERENCES Substance(id);

-- User_Illness
CREATE TABLE IF NOT EXISTS User_Illness (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    illness_id UUID NOT NULL
);

ALTER TABLE User_Illness
ADD CONSTRAINT fk_user_illness
FOREIGN KEY (user_id) REFERENCES Users(id);

ALTER TABLE User_Illness
ADD CONSTRAINT fk_illness_user_illness
FOREIGN KEY (illness_id) REFERENCES Illness(id);

-- User_Substance
CREATE TABLE IF NOT EXISTS User_Substance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    substance_id UUID NOT NULL
);

ALTER TABLE User_Substance
ADD CONSTRAINT fk_user_substance
FOREIGN KEY (user_id) REFERENCES Users(id);

ALTER TABLE User_Substance
ADD CONSTRAINT fk_substance_user_substance
FOREIGN KEY (substance_id) REFERENCES Substance(id);

-- User_Medicine
CREATE TABLE IF NOT EXISTS User_Medicine (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    medicine_id UUID NOT NULL,
    date_of_manufacture DATE NOT NULL,
    quantity float NOT NULL
);

ALTER TABLE User_Medicine
ADD CONSTRAINT fk_user_medicine
FOREIGN KEY (user_id) REFERENCES Users(id);

ALTER TABLE User_Medicine
ADD CONSTRAINT fk_medicine_user_medicine
FOREIGN KEY (medicine_id) REFERENCES Medicine(id);
