-- DROP TABLE IF EXISTS COMPTES;
-- DROP TABLE IF EXISTS AGENDA;
-- DROP TABLE IF EXISTS MESSAGES;

CREATE TABLE COMPTES (
  id TEXT NOT NULL,
  mdp  TEXT NOT NULL,
  description  TEXT NOT NULL,
  PRIMARY KEY(id));
CREATE TABLE AGENDA (
  numero INTEGER PRIMARY KEY AUTOINCREMENT,
  id  TEXT NOT NULL,
  date  TEXT NOT NULL,
  contenu TEXT NOT NULL,
  FOREIGN KEY(id) REFERENCES comptes(id));
CREATE TABLE MESSAGES (
  numero INTEGER PRIMARY KEY AUTOINCREMENT,
  id TEXT NOT NULL,
  dest TEXT NOT NULL,
  date TEXT NOT NULL,  
  contenu TEXT NOT NULL,
  FOREIGN KEY(id) REFERENCES comptes(id),
  FOREIGN KEY(dest) REFERENCES comptes(id));

INSERT INTO COMPTES VALUES ('bot', '123', 'le compte du chatbot')
