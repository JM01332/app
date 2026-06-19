# Programmierworkshop am 19.6.2026

## Namen

- Mike Steidle
- Jerome Meier

## Link zum Git-Repository

[github.com/JM01332/app](https://github.com/JM01332/app)

## KI-Werkzeuge

### Agenten

- OpenAI Codex im IDE-Agentenmodus

### Chat-URLs

- „Codex wurde direkt in der IDE verwendet.“

## Frameworks und Bibliotheken

- Gin – HTTP-Router und REST-Schnittstelle
- GORM v2 – ORM und PostgreSQL-Zugriff
- PostgreSQL-Treiber für GORM
- go-playground/validator – Request-Validierung
- coreos/go-oidc – OIDC-Authentifizierung mit Keycloak
- Zap – strukturiertes Logging
- godotenv – Laden der lokalen Konfiguration
- testing und net/http/httptest – automatisierte Tests

### REST-Schnittstelle (Lesen und Neuanlegen)

- `GET /api/carriers`
- `GET /api/carriers/:id`
- `POST /api/carriers`
- `GET /health`

### Validierung (nur Neuanlegen)

POST-Requests werden mit `go-playground/validator` auf Pflichtfelder,
Textlängen, Carrier-Typ und Security-Level geprüft.

### OR-Mapping (für PostgreSQL)

GORM v2 mit `gorm.io/driver/postgres`. Die bestehenden Tabellen
`carrier`, `command_center` und `aircraft` werden ohne Auto-Migration verwendet.

### Optional: OIDC mit Keycloak

Bearer Tokens werden über OIDC geprüft. GET ist für USER und ADMIN erlaubt,
POST nur für ADMIN. Der Health-Endpunkt bleibt öffentlich.

### Einfacher Integrationstest

`test/integration_test.go` prüft Router, Authentifizierung,
Autorisierung und Carrier-Handler gemeinsam.

Ausführen mit:

```powershell
go test ./...
```

## Prompts/Requests an KI-Agent/en

### Agent-Log Mike Steidle

Dieses Log enthält nur Prompts und Entscheidungen, die zu einem überprüfbaren Projektfortschritt geführt haben. Reine Verständnisfragen, kurze Rückfragen und Nachrichten wie „weiter“ werden nicht einzeln protokolliert.

## 2026-06-19 – Projekt analysiert und Agent-Regeln vorbereitet

**Prompt an Agent:**
Das vorbereitete Go-Projekt prüfen, fehlende Schritte bis zum ersten REST-Health-Check nennen und eine projektweite `AGENTS.md` als Grundlage für weitere Regeln anlegen.

**Ergebnis:**
Die vorhandene Struktur mit `cmd/api`, `internal/app` und `internal/config` wurde geprüft. Leere Go-Dateien, fehlende Serverlogik, Konfiguration, Health-Route und Tests wurden als nächste Schritte identifiziert. Eine projektweite Agent-Datei wurde angelegt und später um den festen Go-Stack, Go-Lernhinweise, Änderungsfreigaben und die Commit-Verantwortung des Nutzers ergänzt.

**Prüfung:**
Repository-Struktur, Aufgabenbeschreibung, README-Vorlage und Git-Status wurden gelesen. Die Regeln schreiben kleine Schritte, Erklärungen für Go-Einsteiger, vorherige Änderungsfreigaben und nutzerseitige Commits vor.

**Commits:**

- `d28a46e docs: add project agent instructions`
- `e060f82 docs: add Go learning guidelines`
- `8e79dec docs: define change approval workflow`

## 2026-06-19 – Go-Entwicklungsumgebung und Abhängigkeiten

**Prompt an Agent:**
Die projektweite `AGENTS.md` lesen und alles installieren, was für den Projektstart benötigt wird.

**Ergebnis:**
Die vorhandene Go-Installation 1.26.4 wurde geprüft. Gin, GORM mit PostgreSQL-Treiber, Validator, godotenv, Zap und Testify wurden dem Go-Modul hinzugefügt und in den lokalen Modulcache geladen. Das optionale OIDC-Paket wurde entsprechend den Projektregeln nicht installiert.

**Prüfung:**
`go mod verify` bestätigt alle Module. `go test ./...` erreicht die Projektpakete, scheitert derzeit jedoch erwartbar, weil die vorbereiteten Go-Dateien noch leer sind und keine `package`-Deklarationen enthalten. Go ist im Windows-Systempfad eingetragen; ein neu gestartetes Terminal übernimmt diesen Pfad.

**Commit:**

- `d28a46e docs: add project agent instructions`
- `2a1193c chore: add Go project dependencies`
- `5434926 chore: add local configuration templates`
- `3f105a3 chore: scaffold compilable Go packages`

## 2026-06-19 – PostgreSQL-Grundlage eingerichtet

**Prompt an Agent:**
Die Datenbankverbindung als gemeinsame Grundlage einrichten, bevor die Aufgaben zwischen zwei Personen aufgeteilt werden.

**Ergebnis:**
Die Anwendung lädt `PORT` und die verpflichtende `DATABASE_URL` aus der Umgebung beziehungsweise einer optionalen lokalen `.env`. Eine zentrale GORM-Verbindung für PostgreSQL prüft die Erreichbarkeit über den Standardbibliotheks-Pool und kann kontrolliert geschlossen werden. Es wurden keine Zugangsdaten und keine Migration hinterlegt.

**Prüfung:**
Konfiguration und DB-Eingabevalidierung sind mit Go-Tests abgedeckt. `go test ./...` ist erfolgreich. Eine echte Verbindung wurde noch nicht geprüft, da lokal keine `DATABASE_URL` konfiguriert ist.

**Commits:**

- `bfd9c2a feat: load and validate application config`
- `95a0e80 feat: add PostgreSQL connection`
- `c28f7ea docs: update database setup checklist`

## 2026-06-19 – Bestehende FastAPI-Datenbank übernommen

**Prompt an Agent:**
Die PostgreSQL-Datenbank aus dem vorherigen FastAPI-Projekt unter `carrier-api` weiterverwenden. Die DB soll dort mit Docker laufen; es sollen keine Daten, Compose-Dateien oder Secrets in das Go-Projekt kopiert werden.

**Ergebnis:**
Das vorhandene SQL-Schema und die SQLAlchemy-Modelle wurden lesend ausgewertet. Datenbank, Rolle und Schema heißen `carrier`; verwendet werden die Tabellen `carrier`, `command_center` und `aircraft`. Die Go-Anwendung verbindet sich als Client über `localhost:5432`. `AutoMigrate` bleibt deaktiviert.

**Prüfung:**
Spalten, PostgreSQL-Enum, Identity-Schlüssel, Beziehungen, Constraints und TLS-Konfiguration wurden mit den Dateien des FastAPI-Projekts abgeglichen. Zugangsdaten bleiben ausschließlich in der lokalen `.env`.

**Commits:**

- `f585c8d docs: document existing PostgreSQL setup`
- `4dd213b chore: connect postgres on startup`

## 2026-06-19 – Erster REST-Health-Check

**Prompt an Agent:**
Als ersten REST-Schritt einen Gin-Router und `GET /health` implementieren und mit `httptest` prüfen.

**Ergebnis:**
Die Router-Erzeugung wurde vom späteren Serverstart getrennt. Der Health-Endpunkt liefert `200 OK` und die JSON-Antwort `{"status":"ok"}`. Gin-Recovery fängt unerwartete Panics während der Request-Verarbeitung ab.

**Prüfung:**
Ein automatisierter Router-Test prüft Statuscode, Content-Type und JSON-Inhalt ohne echten Netzwerk-Port.

**Commit:**
`6d267d7 feat: add health endpoint`

## 2026-06-19 – Startbarer HTTP-Server

**Prompt an Agent:**
Den Serverstart so implementieren, dass der Health-Check ohne laufende Datenbank manuell aufgerufen werden kann.

**Ergebnis:**
`cmd/api/main.go` startet einen `http.Server` mit dem Gin-Router, konfigurierbarem Port, Header-Timeout und strukturiertem Zap-Logging. Die Serverkonfiguration kann unabhängig von der verpflichtenden Datenbankkonfiguration geladen werden. Ein Konsolenbanner zeigt beim Start die Basis- und Health-URL an.

**Prüfung:**
Die Port-Konfiguration ist automatisiert getestet und `go test ./...` ist erfolgreich. Der Server wurde auf einem temporären lokalen Port gestartet; `GET /health` lieferte `200 OK`, den Content-Type `application/json` und `{"status":"ok"}`. Die Datenbank war dafür nicht erforderlich.

**Commits:**

- `7416eb1 feat: start HTTP server`
- `a0b5103 feat: add API startup banner`

## 2026-06-19 – Parallele Umsetzung vorbereitet

**Prompt an Agent:**
Eine verbindliche Dokumentation für die parallele Aufteilung zwischen GORM-Modellen und REST-Schicht erstellen und alle vorab nötigen API- und DB-Entscheidungen aufnehmen.

**Ergebnis:**
`docs/arbeitsaufteilung.md` beschreibt bekannte DB-Fakten, offene API-Entscheidungen, Dateiverantwortung, zwei Parallelphasen, Übergabekriterien und Integrationsprüfungen. Kritische Punkte wie Zeitstempel ohne DB-Default, verschachteltes Anlegen, Transaktionen und Eindeutigkeitskonflikte sind ausdrücklich vor der Implementierung zu klären.

**Prüfung:**
Die Datei trennt bekannte Fakten von offenen Entscheidungen und vermeidet gleichzeitige Änderungen an gemeinsamen Integrationsdateien.

**Commit:**
`800a67e docs: define parallel implementation workflow`

## 2026-06-19 – API- und DB-Verträge festgelegt

**Prompt an Agent:**
Die offenen Entscheidungen in der Arbeitsaufteilung mit sinnvollen, zum bestehenden Carrier-Projekt passenden Werten ausfüllen.

**Ergebnis:**
Als Pflichtumfang wurden `GET /api/carriers` und `POST /api/carriers` festgelegt. Request, Response, Validierungsgrenzen, Fehlerformat, verschachtelte Beziehungen, Transaktion, Preloads, Zeitstempelverantwortung und Service-Schnittstelle sind verbindlich dokumentiert. Der alte FastAPI-Vertrag wurde berücksichtigt, aber die Sicherheitsstufe passend zum DB-Constraint auf 1 bis 5 begrenzt.

**Prüfung:**
Die Entscheidungen wurden mit dem vorhandenen SQL-Schema und den FastAPI-Modellen abgeglichen. Vor Beginn bleibt nur zu prüfen, dass beide Personen vom gleichen Repository-Stand arbeiten.

**Commit:**
`7469718 docs: define carrier API contracts`

## 2026-06-19 – GORM-Modelle fertiggestellt

**Prompt an Agent:**
Die GORM-Modelle langsam und nachvollziehbar Datei für Datei erstellen, nach jeder Datei erklären, testen, reviewen und separat committen.

**Ergebnis:**
Für das bestehende PostgreSQL-Schema wurden `CarrierType`, `CommandCenter`, `Aircraft` und `Carrier` modelliert. Tabellen- und Spaltennamen, Identity-Schlüssel, Fremdschlüssel, Enum, Zeitstempel sowie 1:1- und 1:n-Beziehungen sind mit expliziten GORM-Tags abgebildet. API-DTOs und Migrationen wurden bewusst nicht mit den DB-Modellen vermischt.

**Prüfung:**
Alle Modelle wurden gemeinsam gegen das vorhandene SQL-Schema geprüft. `gofmt` und `go test ./...` sind erfolgreich, der Git-Arbeitsbaum war beim Abschlussreview sauber.

**Commits:**

- `1f8e42d feat: add carrier type model`
- `6b878da feat: add command center model`
- `a5fee17 feat: add aircraft model`
- `7ca9bca feat: add carrier model relationships`

## 2026-06-19 – Carrier-Service und Repository implementiert

**Prompt an Agent:**
Die Service-Schicht wieder schrittweise und nachvollziehbar aufbauen: Create-Input, fachliche Fehler, GORM-Repository, Service und abschließend Tests mit einem Fake-Repository.

**Ergebnis:**
Der Service besitzt HTTP-unabhängige Eingabetypen, den fachlichen Fehler `ErrCarrierNameExists`, ein GORM-Repository für sortiertes Lesen mit Beziehungen und transaktionales Anlegen sowie eine testbare Service-Schicht. PostgreSQL-Namenskonflikte werden anhand des Constraints in den fachlichen Fehler übersetzt. Der Service mappt validierte Eingaben in die verschachtelten GORM-Modelle.

**Prüfung:**
Service-Tests mit einem Fake-Store prüfen Ergebnisse, Fehler, Context-Weitergabe und das vollständige Create-Mapping ohne echte Datenbank. `go test ./...` und `go vet ./...` sind erfolgreich.

**Commits:**

- `96c2db3 feat: add carrier service input types`
- `54a16ee feat: add carrier service errors`
- `e0985d6 feat: add carrier repository`
- `565a8ca feat: add carrier service`
- `bc9c2f6 test: add carrier service tests`

> Die Commit-IDs dieses Abschnitts wurden nach dem Rebase auf die parallel entstandenen Router-Commits aktualisiert.

## 2026-06-19 – Router und Service auf einen Vertrag gebracht

**Prompt an Agent:**
Die Ergebnisse des Teamkollegen nach dem Pull prüfen, verbleibende Fehler korrigieren und erst danach mit der nächsten Arbeitsphase beginnen.

**Ergebnis:**
Die getrennt entwickelten Schichten verwendeten zunächst unterschiedliche Create-Signaturen: Der Handler erwartete `CreateCarrierRequest`, der Service dagegen `CreateCarrierInput`. Der Router mappt den HTTP-Request nun ausdrücklich in den Service-Input. Der doppelte Router-Fehler wurde entfernt; Handler und Tests verwenden gemeinsam `service.ErrCarrierNameExists`. Response-IDs wurden auf `int64` an die GORM-Modelle angeglichen.

**Prüfung:**
Fake-Service und echter Service besitzen denselben Vertrag. Der Handler-Test prüft zusätzlich das Mapping von Name, Carrier-Typ, Command Center und Aircrafts. `go test ./...` und `go vet ./...` sind erfolgreich.

**Relevante Commits:**

- `92b0614 feat: add carrier API DTOs`
- `5a2b677 fix: align response ID types with models`
- `93eae7f feat: map carrier models to API responses`
- `cf61c2a feat: add carrier REST handlers`
- `7e37689 fix: align handlers with carrier service`

## 2026-06-19 – Router-Tests und manuelle Requests vervollständigt

**Prompt an Agent:**
Die Router-Tests als eigener paralleler Arbeitsanteil vervollständigen und anschließend alle noch fehlenden manuellen Request-Dateien anlegen.

**Ergebnis:**
Mapper-Tests sichern beide Übersetzungsrichtungen zwischen HTTP, Service und GORM-Modell ab. Validierungstests prüfen minimale und maximale Feldlängen, fehlendes Command Center sowie die Sicherheitsstufen 1 bis 5. Zusätzliche HTTP-Requests decken unbekannte JSON-Felder und doppelte Carrier-Namen ab.

**Prüfung:**
Leere Aircraft-Listen bleiben als leere Slices erhalten und werden später als `[]` statt `null` ausgegeben. Alle automatisierten Tests sowie `go vet ./...` sind erfolgreich. Die Request-Sammlung enthält GET, gültigen POST, Validierungsfehler, unbekanntes Feld und Namenskonflikt.

**Commits:**

- `745163b test: add carrier mapper tests`
- `55901ab test: extend carrier validation coverage`
- `9de838f test: add carrier API error requests`

## 2026-06-19 – Anwendung vollständig mit PostgreSQL verdrahtet

**Prompt an Agent:**
Den gepullten Integrationsstand prüfen und feststellen, ob Server und Carrier-Requests nun gegen die bestehende PostgreSQL-Datenbank ausgeführt werden können.

**Ergebnis:**
GORM erzeugt Zeitstempel über eine UTC-`NowFunc`. `main.go` öffnet die bestehende PostgreSQL-Verbindung und baut daraus `CarrierRepository` und `CarrierService`. `app.NewRouter` registriert Health- und Carrier-Routen; damit sind `GET /api/carriers`, `GET /api/carriers/:id` nach der späteren Ergänzung und `POST /api/carriers` über denselben produktiven Server erreichbar.

**Prüfung:**
Die Komponentenverdrahtung ist durch App-Tests abgesichert. `go test ./...` und `go vet ./...` sind erfolgreich. Die lokale `.env` wurde ohne Ausgabe sensibler Werte auf Host, Port, Datenbankname und TLS-Modus geprüft. Der anschließende manuelle POST gegen die laufende DB war erfolgreich.

**Commits:**

- `47ffd76 chore: configure gorm UTC timestamps`
- `28ae2a5 feat: wire carrier API components`
- `300672b merge remote integration updates`

## 2026-06-19 – POST-Response auf Location-Header reduziert

**Prompt an Agent:**
Beim erfolgreichen POST den Datensatz nicht erneut im Response-Body zurückgeben, sondern nur `201 Created` und die neue Ressourcenadresse im `Location`-Header senden.

**Ergebnis:**
Der POST-Handler verwendet die von PostgreSQL erzeugte Carrier-ID nur noch für `/api/carriers/{id}` und antwortet ohne JSON-Body. Der Service gibt den gespeicherten Carrier intern weiterhin zurück, weil dessen generierte ID für den Header benötigt wird. GET-Antworten bleiben unverändert.

**Prüfung:**
Der Handler-Test prüft Statuscode, Location-Header und einen vollständig leeren Response-Body. Der lokale API-Vertrag wurde entsprechend angepasst.

**Commit:**
`3d35eb8 fix: return empty body after carrier creation`

## 2026-06-19 – Carrier-Detailabfrage ergänzt

**Prompt an Agent:**
Prüfen, ob Carrier per ID gelesen werden können, und den gepullten Detail-GET vor dem nächsten Schritt vollständig reviewen.

**Ergebnis:**
`GET /api/carriers/:id` liest einen Carrier einschließlich Command Center und sortierter Aircrafts. Ungültige IDs liefern `400 Bad Request`, unbekannte IDs `404 Not Found` und technische Fehler kontrolliert `500 Internal Server Error`.

**Prüfung:**
Repository, Service, Handler und Fakes verwenden gemeinsam `GetByID(ctx, id)`. Tests decken Erfolg, ungültige ID, nicht gefunden und internen Fehler ab. `go test ./...` und `go vet ./...` sind erfolgreich.

**Commits:**

- `b65bbd1 feat: add carrier detail lookup`
- `893c821 feat: add carrier detail endpoint`
- `55aad1a docs: add carrier detail request`

## 2026-06-19 – Strukturiertes HTTP-Logging ergänzt

**Prompt an Agent:**
Nach erfolgreichem Review des Detail-GET das Logging mit Middleware korrekt umsetzen.

**Ergebnis:**
Eine Gin-Middleware protokolliert jeden abgeschlossenen Request mit Methode, Route, Status, Dauer und Client-IP über Zap. Unerwartete Handlerfehler werden über den Gin-Context an die Middleware gemeldet und auf Error-Level inklusive Ursache protokolliert. Eine eigene Recovery-Middleware fängt Panics ab, protokolliert Panic und Stacktrace und sendet eine kontrollierte `500`-Antwort ohne interne Details. Request-Bodies und Secrets werden nicht geloggt.

**Prüfung:**
Middleware-Tests prüfen Info-Logging erfolgreicher Requests, Error-Logging bei `500` und Panic-Recovery. Der vollständige Testlauf und `go vet ./...` sind erfolgreich.

**Commit:**
`7e343e3 feat: add structured HTTP logging middleware`

## 2026-06-19 – Carrier-Routen mit Keycloak OIDC abgesichert

**Prompt an Agent:**
Den bereits laufenden Keycloak-Server des vorherigen Projekts weiterverwenden und die Carrier-Routen mit Bearer-Token-Prüfung sowie rollenbasierter Autorisierung absichern. Health soll öffentlich bleiben, GET soll für USER und ADMIN erlaubt sein und POST nur für ADMIN.

**Ergebnis:**
Die Anwendung lädt Issuer, Client-ID und den Pfad zum lokalen CA-Zertifikat aus der Umgebung. Beim Start werden die OIDC-Metadaten und JWKS von Keycloak geladen. Die Middleware prüft Signatur, Issuer und Ablaufzeit. Als Clientbindung akzeptiert sie `aud=python-client` oder, kompatibel zur vorhandenen Keycloak-Konfiguration, `azp=python-client`. Anschließend wertet sie die Client-Rollen aus `resource_access.python-client.roles` aus. `/health` bleibt ohne Token erreichbar; beide GET-Routen akzeptieren USER oder ADMIN, während POST die ADMIN-Rolle verlangt. Ungültige oder fehlende Tokens liefern `401`, fehlende Berechtigungen `403`. Tokens und Secrets werden weder gespeichert noch geloggt.

**Prüfung:**
Middleware-, Verifier- und App-Tests sichern Authentifizierung, Clientbindung, Rollen-Normalisierung und die gesamte Routenmatrix ab. `go test ./...`, `go vet ./...` und `git diff --check` sind erfolgreich. Die bestehende Keycloak-Konfiguration muss dadurch nicht angepasst werden.

Der abschließende Live-Test war ebenfalls erfolgreich: Der Server konnte die OIDC-Konfiguration des bestehenden Keycloak-Servers laden. Ein authentifizierter GET-Request in Bruno las Carrier aus PostgreSQL. Ein POST mit ADMIN-Rolle erzeugte einen neuen Carrier, antwortete mit `201 Created`, leerem Body und dem erwarteten `Location`-Header.

**Commit:**
`a89c35e feat: secure carrier routes with Keycloak OIDC`

### Agent-Log Jérôme Meier

## 2026-06-19 15:10 +02:00 - PostgreSQL beim Serverstart pruefen

**Prompt an Agent:**
Der Nutzer hat freigegeben, die passenden Dateien fuer den naechsten Schritt anzulegen und den kompletten docs-Ordner in die .gitignore aufzunehmen.

**Ergebnis des Agenten:**
`cmd/api/main.go` laedt jetzt die vollstaendige Anwendungskonfiguration, oeffnet beim Start die PostgreSQL-Verbindung ueber GORM und schliesst den Verbindungspool beim Beenden. `.gitignore` ignoriert nun den kompletten `docs/`-Ordner.

**Eigene Pruefung:**
Noch offen: `go test ./...` und ein lokaler Serverstart, sobald `go` in der Shell verfuegbar ist und die lokale `.env` eine gueltige `DATABASE_URL` enthaelt.

**Commit:**
Noch offen. Vorschlag: `chore: connect postgres on startup`

## 2026-06-19 15:29 +02:00 - Carrier API-DTOs und Dokumentation

**Prompt an Agent:**
Der Nutzer hat die aktualisierte Arbeitsaufteilung bestaetigt und beauftragt, vorerst Person B umzusetzen: API-DTOs und Dokumentation, logisch in kleine Commit-Pakete geteilt.

**Ergebnis des Agenten:**
Angelegt wurden die Create-Request-DTOs mit JSON- und Validator-Tags, Response-DTOs inklusive Fehlerformat, Validierungstests fuer positive und negative Faelle, API-Dokumentation sowie Beispielrequests fuer Liste, gueltigen POST und ungueltigen POST.

**Eigene Pruefung:**
`gofmt` und `go test ./...` konnten in dieser Shell nicht ausgefuehrt werden, weil `go`/`gofmt` nicht im PATH gefunden werden. Die Tests muessen lokal nachgeholt werden.

**Commit:**
Noch offen. Vorschlaege: `feat: add carrier API DTOs` und `docs: document carrier API requests`

## 2026-06-19 16:01 +02:00 - Carrier REST-Handler mit Fake-Service

**Prompt an Agent:**
Der Nutzer hat beauftragt, fuer Person B die Dateien `internal/carrier/router/mapper.go`, `handler.go` und `handler_test.go` umzusetzen: Model-zu-Response-Mapping, GET- und POST-Handler, Request-Binding und Validierung, Ablehnung unbekannter JSON-Felder, Fehlerantworten und Statuscodes sowie Fake-Service fuer Handler-Tests. Dabei soll `int64` fuer IDs beachtet werden.

**Ergebnis des Agenten:**
Angelegt wurden Mapper fuer Carrier-, Command-Center- und Aircraft-Responses, ein Gin-Handler mit strengem JSON-Decoding, Validator-Pruefung und einheitlichen Fehlerantworten sowie Handler-Tests mit Fake-Service fuer GET, POST, ungueltiges JSON, unbekannte Felder, Validierungsfehler, Namenskonflikt und interne Fehler.

**Eigene Pruefung:**
`gofmt` und `go test ./...` konnten in dieser Shell nicht ausgefuehrt werden, weil `go`/`gofmt` hier nicht im PATH gefunden werden. Die Pruefung muss lokal nachgeholt werden.

**Commit:**
Noch offen. Vorschlag: `feat: add carrier REST handlers`

## 2026-06-19 16:23 +02:00 - Carrier API verdrahten

**Prompt an Agent:**
Der Nutzer hat beauftragt, nach dem Zusammenfuehren der Person-A- und Person-B-Arbeit die Integration umzusetzen: GORM-Zeitstempel mit UTC-`NowFunc` konfigurieren, Repository und Service mit `postgres.DB` erzeugen, den Service an den Gin-Router uebergeben, Carrier-Routen unter `/api` registrieren, Health-Test anpassen und den vollstaendigen Start vorbereiten.

**Ergebnis des Agenten:**
`internal/database/postgres.go` nutzt nun eine GORM-Konfiguration mit UTC-`NowFunc`. `internal/app/app.go` nimmt einen Carrier-Service entgegen und registriert Carrier-Routen unter `/api`. `cmd/api/main.go` erzeugt Repository und Service aus der bestehenden PostgreSQL-Verbindung und uebergibt den Service an den Router. Tests fuer UTC-NowFunc und Routenregistrierung wurden ergaenzt.

**Eigene Pruefung:**
`gofmt` und `go test ./...` konnten in dieser Shell nicht ausgefuehrt werden, weil `go`/`gofmt` hier nicht im PATH gefunden werden. Die Pruefung muss lokal nachgeholt werden. Ein vollstaendiger Start sollte lokal mit laufender PostgreSQL-Datenbank und gueltiger `.env` getestet werden.

**Commit:**
Noch offen. Vorschlaege: `chore: configure gorm UTC timestamps` und `feat: wire carrier API components`

## 2026-06-19 16:35 +02:00 - Schema-qualifizierte Carrier-Tabellen

**Prompt an Agent:**
Nach dem Live-Test lieferten GET und gueltiger POST `internal_error`. Die Serverlogs zeigten `ERROR: relation "carrier" does not exist`, daher sollte der GORM-Tabellenname gegen das bestehende PostgreSQL-Schema geprueft und korrigiert werden.

**Ergebnis des Agenten:**
Die `TableName()`-Methoden der Carrier-Modelle verwenden jetzt schema-qualifizierte Tabellennamen: `carrier.carrier`, `carrier.command_center` und `carrier.aircraft`. Ein kleiner Model-Test prueft diese Namen.

**Eigene Pruefung:**
Noch offen: lokal `gofmt`, `go test ./...`, Server neu starten und GET/POST erneut gegen die laufende Datenbank testen.

**Commit:**
Noch offen. Vorschlag: `fix: qualify carrier table names with schema`

## 2026-06-19 16:50 +02:00 - Carrier per ID lesen

**Prompt an Agent:**
Der Nutzer hat den aktuellen Stand nach dem Zusammenfuehren gepullt und beauftragt, zusaetzlich einen GET-Request mit ID zu ergaenzen.

**Ergebnis des Agenten:**
Repository, Service und Router wurden um `GetByID` erweitert. Der neue Endpunkt `GET /api/carriers/:id` validiert positive numerische IDs, liefert `200 OK` mit Carrier-Response, `400` fuer ungueltige IDs, `404` fuer nicht gefundene Carrier und `500` fuer unerwartete Fehler. Tests und eine Bruno/HTTP-Request-Datei wurden ergaenzt.

**Eigene Pruefung:**
Noch offen: lokal `gofmt`, `go test ./...`, Server neu starten und `GET /api/carriers/{id}` gegen die laufende Datenbank testen.

**Commit:**
Noch offen. Vorschlaege: `feat: add carrier detail lookup` und `docs: add carrier detail request`

## 2026-06-19 17:20 +02:00 - OIDC mit bestehendem Keycloak vorbereiten

**Prompt an Agent:**
Der Nutzer hat beauftragt, die bestehende Keycloak-Installation aus `carrier-api` weiterzuverwenden und OIDC in der Go-API optional zu integrieren.

**Ergebnis des Agenten:**
Aus dem alten Projekt wurden Realm `python`, Client-ID `python-client` und der lokale Token-Endpunkt `http://localhost:8880/realms/python/protocol/openid-connect/token` ermittelt. Die Go-Konfiguration liest nun optionale OIDC-Variablen. Ein neues `internal/auth`-Package enthaelt OIDC-Verifier und Gin-Middleware. `/health` bleibt oeffentlich, `/api` wird bei aktivierter OIDC-Konfiguration geschuetzt. Request-Dateien enthalten Bearer-Token-Platzhalter und es gibt eine Token-Request-Datei ohne echte Secrets.

**Eigene Pruefung:**
Noch offen: lokal `go get github.com/coreos/go-oidc/v3`, `gofmt`, `go test ./...`, Start mit `OIDC_ENABLED=true` und Test mit echtem Keycloak-Token. Die Sandbox-Shell findet `go` nicht im PATH.

**Commit:**
Noch offen. Vorschlaege: `chore: add oidc dependency`, `chore: add oidc configuration`, `feat: add oidc auth middleware`, `feat: protect carrier API routes`, `docs: add authenticated requests`

## 2026-06-19 18:08 +02:00 - Einfacher Integrationstest

**Prompt an Agent:**
Der Nutzer hat nach dem Zusammenfuehren der Auth-/Rollenpruefung einen einfachen Integrationstest fuer den aktuellen Stand angefordert.

**Ergebnis des Agenten:**
Angelegt wurde `test/integration_test.go`. Der Test erzeugt den kompletten Gin-Router mit Fake-Service und Fake-Token-Verifier und prueft ohne echte Datenbank oder Keycloak: `/api/carriers` ohne Token liefert `401`, mit User-Rolle liefert die Liste `200`, und `POST /api/carriers` mit User-Rolle liefert `403`.

**Eigene Pruefung:**
Noch offen: lokal `gofmt` und `go test ./...` ausfuehren.

**Commit:**
Noch offen. Vorschlag: `test: add authenticated API integration test`
