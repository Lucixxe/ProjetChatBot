CREATE TABLE COMPTES (
  id TEXT NOT NULL,
  mdp  TEXT NOT NULL,
  description  TEXT NOT NULL,
  PRIMARY KEY(id));
CREATE TABLE AGENDA (
  id  TEXT NOT NULL,
  date  TEXT NOT NULL,
  contenu TEXT NOT NULL,
  FOREIGN KEY(id) REFERENCES comptes(id),
  PRIMARY KEY(id, date));
CREATE TABLE MESSAGE (
  id TEXT NOT NULL,
  date TEXT NOT NULL,
  contenu TEXT NOT NULL,
  FOREIGN KEY(id) REFERENCES comptes(id),
  PRIMARY KEY(id, date, contenu));
