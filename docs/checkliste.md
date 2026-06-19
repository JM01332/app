# Projektcheckliste

Diese Checkliste begleitet die Umsetzung der Go-REST-API. Zuständigkeiten werden mit Namen oder Kürzeln ergänzt. Ein Abschnitt gilt erst als abgeschlossen, wenn seine Prüfungen erfolgreich waren.

## 1. Aufgabe gemeinsam klären

- [ ] Vollständige Fachaufgabe lesen
- [ ] Ressource und Datenfelder festhalten
- [ ] Bestehende PostgreSQL-Tabelle und Spalten prüfen
- [ ] Geforderte GET- und POST-Endpunkte festlegen
- [ ] Request- und Response-JSON vereinbaren
- [ ] Pflichtfelder und Validierungsregeln festlegen
- [ ] Benötigte Umgebungsvariablen bestimmen
- [ ] Verantwortlichkeiten eintragen

**Person A – API und Anwendung:**

**Person B – Datenbank und Fachlogik:**

## 2. Technische Grundlage

### Person A – API und Anwendung

- [ ] Konfiguration aus Umgebungsvariablen laden
- [ ] `PORT` mit sinnvollem Standardwert unterstützen
- [ ] Zap-Logger einrichten
- [ ] Gin-Router erstellen
- [ ] `GET /health` implementieren
- [ ] Health-Check mit `httptest` testen
- [ ] Server über `cmd/api/main.go` starten

### Person B – Datenbank und Modell

- [ ] Zugriff auf den bestehenden PostgreSQL-Server prüfen
- [ ] `DATABASE_URL` oder vereinbarte DB-Variablen dokumentieren
- [ ] GORM-Verbindung implementieren
- [ ] Verbindungsfehler kontrolliert behandeln und loggen
- [ ] Bestehendes Datenbankschema prüfen
- [ ] GORM-Modell anhand des Schemas anlegen
- [ ] Keine Migration ohne vorherige Abstimmung ausführen

### Gemeinsame Prüfung

- [ ] `go test ./...` ist erfolgreich
- [ ] Anwendung startet lokal
- [ ] `GET /health` liefert `200 OK`
- [ ] Health-Check liefert das vereinbarte JSON
- [ ] PostgreSQL-Verbindung funktioniert
- [ ] Keine Secrets befinden sich in Git
- [ ] Zwischenstand committen

## 3. Fachliche REST-Schnittstelle

### Person A – HTTP-Schicht

- [ ] Request-DTO für POST definieren
- [ ] Response-DTO definieren
- [ ] GET-Handler implementieren
- [ ] POST-Handler implementieren
- [ ] JSON-Binding mit Gin einrichten
- [ ] POST-Daten mit Validator prüfen
- [ ] Feldbezogene Validierungsfehler zurückgeben
- [ ] Erfolgreiches GET mit `200 OK` beantworten
- [ ] Erfolgreiches POST mit `201 Created` beantworten
- [ ] Wenn sinnvoll, `Location`-Header setzen
- [ ] Ungültiges JSON mit `400 Bad Request` beantworten
- [ ] Validierungsfehler mit `422 Unprocessable Content` beantworten

### Person B – Fachlogik und Datenzugriff

- [ ] Service-Funktionen mit Person A abstimmen
- [ ] Lesen der Datensätze implementieren
- [ ] Neuanlegen eines Datensatzes implementieren
- [ ] GORM-Fehler in verständliche Anwendungsfehler übersetzen
- [ ] Falls gefordert, Duplikate als Konflikt behandeln
- [ ] Datenbankmodelle nicht ungeprüft als API-Antwort verwenden
- [ ] Service- oder Repository-Tests ergänzen

### Gemeinsame Integration

- [ ] Handler und Service verbinden
- [ ] Routen im Gin-Router registrieren
- [ ] Context durch Handler, Service und Datenbankzugriff reichen
- [ ] Unerwartete Fehler serverseitig loggen
- [ ] Keine internen Fehlerdetails an Clients senden
- [ ] Zwischenstand committen

## 4. Tests und manuelle Requests

- [ ] Health-Check-Test ist erfolgreich
- [ ] GET-Endpunkt ist getestet
- [ ] Gültiger POST ist getestet
- [ ] POST mit ungültigem JSON ist getestet
- [ ] POST ohne Pflichtfeld ist getestet
- [ ] Erwartete Statuscodes sind geprüft
- [ ] Beispielrequests unter `requests/` anlegen
- [ ] `gofmt` auf geänderte Go-Dateien anwenden
- [ ] `go test ./...` ist erfolgreich
- [ ] Manuellen Test gegen die laufende Anwendung durchführen
- [ ] Testergebnis dokumentieren und committen

## 5. Dokumentation und Abgabe

- [ ] Vorgegebene `README.md` ins Repository übernehmen
- [ ] Namen und Repository-Link eintragen
- [ ] Verwendete KI-Werkzeuge und Chat-URLs eintragen
- [ ] REST-Schnittstelle beschreiben
- [ ] POST-Validierung beschreiben
- [ ] GORM- und PostgreSQL-Nutzung beschreiben
- [ ] Integrationstest beschreiben
- [ ] OIDC-Status ehrlich dokumentieren
- [ ] Wichtige Prompts und eigene Prüfungen dokumentieren
- [ ] Lokales `docs/agent-log.md` aktualisieren
- [ ] Keine persönlichen Daten oder Secrets dokumentieren

## 6. Abschlussprüfung

- [ ] `gofmt` meldet keine offenen Formatierungsänderungen
- [ ] `go test ./...` ist erfolgreich
- [ ] Anwendung startet mit dokumentierter Konfiguration
- [ ] PostgreSQL ist erreichbar
- [ ] GET liefert plausible Datensätze
- [ ] POST legt einen gültigen Datensatz an
- [ ] Ungültiger POST liefert einen kontrollierten Fehler
- [ ] README ist vollständig
- [ ] Git-Historie besteht aus kleinen, verständlichen Commits
- [ ] `git status --short` ist sauber
- [ ] Keine `.env`, Zugangsdaten oder Tokens sind versioniert
- [ ] Optionales OIDC wurde erst nach den Pflichtteilen begonnen
