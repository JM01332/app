# AGENTS.md

## Kontext und Abgabe

Dieses Repository entsteht im **Programmierworkshop am 19.06.2026**.

Gegeben sind:

- **Plattform:** Go
- **Datenbank:** der PostgreSQL-Server aus den vorherigen Abgaben
- **OIDC mit Keycloak:** optional
- **Abgabe:** die ausgefüllte, vorgegebene `README.md` aus dem ILIAS-Ordner per E-Mail an den Dozenten

Die konkrete fachliche Aufgabenstellung (Ressource, Felder und genaue API-Anforderungen) wird noch bereitgestellt. Bis dahin keine fachlichen Modelle oder Endpunkte erfinden.

## Fester Technologie-Stack

Verwende für diese Aufgabe den folgenden Stack, sofern die spätere Fachaufgabe nichts Widersprüchliches vorgibt:

| Zweck | Vorgabe |
|---|---|
| Sprache / Plattform | Go |
| REST-Framework | Gin (`github.com/gin-gonic/gin`) |
| ORM für PostgreSQL | GORM v2 (`gorm.io/gorm`, `gorm.io/driver/postgres`) |
| Validierung | `github.com/go-playground/validator/v10` |
| Konfiguration | `github.com/joho/godotenv` für lokale `.env`-Dateien; Umgebungsvariablen bleiben die eigentliche Quelle |
| Logging | `go.uber.org/zap` |
| Integrationstest | Go-Standardpakete `testing` und `net/http/httptest`; Assertions mit `github.com/stretchr/testify` |
| Optionales OIDC | `github.com/coreos/go-oidc` |

Keine zusätzliche Bibliothek einführen, wenn Gin, GORM, Validator, Zap oder die Go-Standardbibliothek die Anforderung bereits sinnvoll abdecken.

## Go-Lernkontext und Erklärpflicht

Die Beteiligten verwenden Go in diesem Projekt zum ersten Mal. Der Agent soll deshalb nicht nur funktionierenden Code liefern, sondern die Arbeit mit Go nachvollziehbar vermitteln.

- Im gesamten Projekt etablierte Go-Best-Practices und idiomatische Go-Konventionen verwenden.
- Bei jeder Änderung kurz erklären, warum die gewählte Lösung für Go passend ist.
- Neu verwendete Go-Sprachmittel und Standardbibliotheksfunktionen verständlich erklären, beispielsweise Packages, Structs, Interfaces, Methoden, Pointer, Fehlerbehandlung, `defer`, Goroutines oder Contexts.
- Bei Framework- oder Bibliotheksfunktionen erklären, welche Aufgabe sie übernehmen und wie sie in den Ablauf der Anwendung passen.
- Erklärungen am Wissensstand von Go-Einsteigern ausrichten, Fachbegriffe beim ersten Auftreten erläutern und kleine konkrete Beispiele bevorzugen.
- Keine unnötig fortgeschrittenen Sprachmittel oder Abstraktionen einsetzen. Falls eine komplexere Lösung nötig ist, deren Nutzen und Alternative benennen.
- Best Practices nicht nur behaupten, sondern anhand der konkreten Änderung begründen.
- Den Code trotzdem übersichtlich halten; Lernhinweise gehören primär in die Kommunikation und Dokumentation, nicht als übermäßige Kommentare in jede Codezeile.

## Initialisierung

Wenn noch kein Go-Starterprojekt vorhanden ist, initialisiere das Projekt schlank und nachvollziehbar:

```bash
git status --short
go mod init <modulpfad-aus-dem-git-repository>
go get github.com/gin-gonic/gin
go get gorm.io/gorm gorm.io/driver/postgres
go get github.com/go-playground/validator/v10
go get github.com/joho/godotenv
go get go.uber.org/zap
go get github.com/stretchr/testify
```

`github.com/coreos/go-oidc` nur installieren, wenn OIDC/Keycloak tatsächlich umgesetzt wird.

Falls ein Starterprojekt vorhanden ist:

- zuerst `README.md`, `go.mod`, vorhandene Packages und `.gitignore` lesen,
- vorhandene Konventionen übernehmen,
- nicht durch ein neues `go mod init` oder einen generierten Projektbaum ersetzen.

Nach der Initialisierung sofort mit `go run ./cmd/api` oder dem tatsächlich vorgesehenen Startbefehl prüfen, ob der Server grundsätzlich läuft. Der konkrete Port kommt aus `PORT` oder wird dokumentiert.

## Arbeitsreihenfolge nach Ausgabe der Fachaufgabe

1. Lies Aufgabenstellung, `README.md`, `go.mod`, `.gitignore` und den aktuellen Git-Status.
2. Ermittle Ressource, Felder, DB-Tabelle bzw. Tabellen, GET- und POST-Anforderungen sowie Sonderfälle.
3. Prüfe die vorhandene PostgreSQL-Verbindung aus den vorherigen Abgaben. Verwende nur Konfiguration/Umgebungsvariablen; keine lokalen Pfade oder Zugangsdaten in Code oder Git.
4. Erstelle eine kurze Umsetzungsskizze mit kleinen, überprüfbaren Schritten.
5. Implementiere zunächst einen minimal startbaren Gin-Server und die DB-Verbindung.
6. Implementiere `GET` zum Lesen und `POST` zum Neuanlegen.
7. Ergänze Validierung nur für `POST`.
8. Ergänze einen einfachen Integrationstest mit `httptest`.
9. Aktualisiere README, Agent-Log und ggf. API-Dokumentation fortlaufend.
10. Führe manuelle Checks und `go test ./...` aus und dokumentiere das Ergebnis ehrlich.

## Zielstruktur

Die Struktur darf an die Fachdomäne angepasst werden, soll aber klar und Go-typisch bleiben:

```text
.
├── cmd/
│   └── api/
│       └── main.go              # Startpunkt
├── internal/
│   ├── app/
│   │   └── app.go               # Router-Setup, Middleware, Fehlerbehandlung
│   ├── config/
│   │   └── config.go            # Environment und zentrale Konfiguration
│   ├── database/
│   │   └── postgres.go          # GORM-/PostgreSQL-Verbindung
│   └── <domain>/
│       ├── model.go             # GORM-Modelle und API-DTOs
│       ├── repository.go        # DB-Zugriffe (nur falls es die Klarheit verbessert)
│       ├── service.go           # Fachlogik
│       └── handler.go           # Gin-Handler, HTTP und Binding
├── test/
│   └── integration_test.go      # einfacher HTTP-Integrationstest, falls passend
├── docs/                        # Zusatzdokumentation bei Bedarf
├── .env.example                 # nur Variablennamen und Beispielwerte, keine Secrets
├── go.mod
├── go.sum
└── README.md                    # vorgegebene Abgabe-README
```

Keine unnötige Abstraktion erzeugen. Bei nur einer Ressource dürfen Repository und Service zusammengefasst werden, solange HTTP, Fachlogik und Datenbankzugriff nachvollziehbar getrennt bleiben.

## REST-Schnittstelle

Implementiere ausschließlich die verlangte Ressource und die geforderten Operationen. Als erwarteter Kern gelten:

- `GET /api/<ressource>`: vorhandene Datensätze lesen,
- `POST /api/<ressource>`: einen neuen Datensatz anlegen.

Regeln:

- Ressourcen im Plural benennen.
- Gin-Handler bleiben dünn: Request binden, validieren, Service aufrufen, HTTP-Antwort senden.
- GORM-Modelle nicht ungeprüft direkt nach außen geben; Request- und Response-DTOs verwenden, wenn Datenbank- und API-Form auseinanderlaufen.
- Erfolgreiches `GET`: `200 OK`.
- Erfolgreiches `POST`: `201 Created`, möglichst mit `Location`-Header.
- Fehlende oder ungültige JSON-Daten: `400 Bad Request`.
- Validierungsfehler beim POST: `422 Unprocessable Content` mit feldbezogenen, verständlichen Fehlern.
- Fachliche Konflikte, etwa Duplikate: `409 Conflict`, sofern die Fachaufgabe dies erfordert.
- Unerwartete Fehler: `500 Internal Server Error`, ohne interne DB- oder Stack-Details im Response-Body.

## Validierung beim Neuanlegen

Für `POST`:

- Request mit Gin binden, z. B. `ShouldBindJSON`.
- Validierungsregeln an einem dedizierten Create-DTO mit Validator-Tags definieren, z. B. `required`, `min`, `max`, `email`, soweit fachlich passend.
- `validator/v10` für feldbezogene Fehler nutzen; Gin verwendet diese Bibliothek auch für viele Binding-/Validation-Fälle, die Fehlerantwort soll aber von der Anwendung bewusst und einheitlich formatiert werden.
- Keine Validierung für Leseoperationen ergänzen, außer einfache technische Prüfung von zwingenden Query-Parametern ist explizit gefordert.

## PostgreSQL und GORM

- PostgreSQL über GORM v2 und `gorm.io/driver/postgres` anbinden.
- Konfiguration über Umgebungsvariablen, beispielsweise `DATABASE_URL` oder die vom bestehenden DB-Projekt vorgegebene Konfiguration.
- Lokale `.env` nur für Entwicklung verwenden; `.env` bleibt in `.gitignore`.
- `.env.example` darf nur Platzhalter enthalten.
- Vor einem `AutoMigrate` prüfen, ob die bestehende Datenbank bzw. das Schema aus vorherigen Abgaben dadurch nicht unbeabsichtigt verändert wird. Migrationen nur ausführen, wenn die Fachaufgabe und der aktuelle DB-Zustand dies zulassen.
- Bestehende Tabellen und Daten respektieren. Keine DB-Ordner aus dem früheren Projekt in dieses Repository kopieren.

## Konfiguration und Logging

- `godotenv` ausschließlich zum lokalen Laden einer optionalen `.env` nutzen; die Anwendung soll auch mit gesetzten Systemumgebungsvariablen funktionieren.
- Konfiguration zentral in `internal/config` lesen und früh validieren.
- Zap als strukturierten Logger einrichten; mindestens Start, DB-Verbindungsfehler und unerwartete Serverfehler loggen.
- Keine Passwörter, vollständigen Connection-Strings oder Tokens loggen.

## Optional: OIDC mit Keycloak

OIDC/Keycloak ist **optional**. Erst die Pflichtteile (Server, DB, GET, POST, POST-Validierung, README und Test) fertigstellen.

Nur bei ausreichender Zeit:

- `github.com/coreos/go-oidc` nutzen,
- Issuer, Client-ID und ggf. Audience ausschließlich über Environment-Variablen konfigurieren,
- Bearer-Token über Middleware prüfen,
- mindestens eine geschützte Route dokumentieren,
- keine Keycloak-Secrets, Realm-Exports mit Zugangsdaten oder Tokens committen.

Falls OIDC nicht umgesetzt wird, README ehrlich mit „nicht umgesetzt (optional, Zeitpriorität auf Pflichtteile)“ ausfüllen.

## Einfacher Integrationstest

Ein einfacher Test ist Teil der geplanten Abgabe:

- Router/App so aufbauen, dass er in einem Test ohne echten Listener erzeugt werden kann.
- `net/http/httptest` und `testing` verwenden.
- `testify/require` oder `testify/assert` nur für gut lesbare Assertions einsetzen.
- Mindestens einen HTTP-Fall testen, bevorzugt `GET /api/<ressource>` oder einen POST-Validierungsfehler, der ohne produktive Datenbank robust ausführbar ist.
- Keine aufwendige neue Testinfrastruktur oder Testcontainer einführen, außer die Aufgabenstellung verlangt sie ausdrücklich.
- Test mit `go test ./...` ausführen und Ergebnis im Agent-Log festhalten.

## Verbindliche README-Abgabe

Die vorgegebene `README.md` ist die Hauptabgabe und muss vollständig ausgefüllt werden:

```md
# Programmierworkshop am 19.6.2026

## Namen

## Link zum Git-Repository

## KI-Werkzeuge

### Agenten

### Chat-URLs, z.B. https://chatgpt.com

## Frameworks und Bibliotheken

### REST-Schnittstelle (Lesen und Neuanlegen)

### Validierung (nur Neuanlegen)

### OR-Mapping (für PostgreSQL)

### Optional: OIDC mit Keycloak

### Einfacher Integrationstest

## Prompts/Requests an KI-Agent/en
```

Inhaltlich soll sie mindestens festhalten:

- Namen aller Beteiligten und Repository-Link,
- genutzte KI-Werkzeuge und Chat-URLs,
- Gin für die REST-Schnittstelle,
- `validator/v10` für POST-Validierung,
- GORM v2 + PostgreSQL-Treiber für ORM,
- ob OIDC/Keycloak umgesetzt wurde oder nicht,
- den verwendeten Integrationstest,
- die wichtigsten Prompts/Requests an Agenten, inklusive eigener Prüfung und Entscheidungen.

## Protokollierung der KI-Nutzung

Zusätzlich zur README die lokale Datei `docs/agent-log.md` pflegen. Sie ist persönliche Dokumentation, wird über `.gitignore` ausgeschlossen und nicht committed. Nach jedem größeren Schritt notieren:

- Datum/Uhrzeit,
- Ziel des Arbeitsschritts,
- verwendeten Prompt oder präzise Zusammenfassung,
- Ergebnis des Agenten,
- eigene Prüfung, Korrektur oder Entscheidung,
- zugehörigen Git-Commit.

Beispiel:

```md
## 2026-06-19 – POST-Endpunkt und Validierung

**Prompt an Agent:**
„Implementiere POST /api/books mit Gin, GORM und einem Create-DTO. Validiere die Felder vor dem DB-Zugriff und antworte bei Validierungsfehlern mit 422. Ändere keine fremden Dateien.“

**Eigene Prüfung:**
Ein gültiger Request liefert 201. Ein Body ohne Pflichtfeld liefert 422 mit Feldname.

**Commit:**
`feat: add book creation endpoint with validation`
```

Keine Secrets, privaten Zugangsdaten oder vollständigen lokalen Pfade in README oder Agent-Log schreiben.

## Git- und Commit-Regeln

### Änderungsfreigabe und Commit-Verantwortung

- Der Agent darf Dateien nur verändern, nachdem der Nutzer die konkrete Änderung ausdrücklich freigegeben hat.
- Reine Lese-, Analyse- und Prüfaktionen dürfen ohne Änderungsfreigabe durchgeführt werden.
- Eine allgemeine Aufgabenbesprechung gilt nicht automatisch als Freigabe zur Implementierung; vor dem Schreiben nennt der Agent kurz die vorgesehenen Dateien und Änderungen und fragt nach Zustimmung.
- Der Agent erstellt keine Git-Commits.
- Nach einem abgeschlossenen und geprüften Arbeitsschritt weist der Agent den Nutzer darauf hin, dass ein Commit sinnvoll ist, nennt die betroffenen Dateien und schlägt eine Conventional-Commit-Nachricht vor.
- Der Nutzer prüft und erstellt den Commit selbst.

Vor jedem Commit:

```bash
git status --short
```

Regeln:

- Nur Dateien committen, die zum aktuellen Schritt gehören.
- Kleine, lauffähige und erklärbare Commits verwenden.
- Keine `.env`, Datenbank-Credentials, Tokens oder fremden Änderungen committen.
- Conventional-Commits-Stil verwenden.

Empfohlene Commit-Reihenfolge:

```text
docs: record workshop task and Go stack
chore: initialize Go module
chore: configure Gin application
chore: add PostgreSQL configuration
feat: add resource read endpoint
feat: add resource creation endpoint
feat: validate creation request
test: add API integration smoke test
docs: complete workshop README
```

## Abschlusschecks

Vor Abgabe mindestens ausführen und dokumentieren:

```bash
gofmt -w .
go test ./...
git status --short
```

Zusätzlich manuell prüfen:

- Anwendung startet.
- PostgreSQL ist erreichbar.
- `GET /api/<ressource>` liefert plausibel Daten.
- `POST /api/<ressource>` legt mit gültigem Body einen Datensatz an.
- Ungültiger POST liefert kontrolliert `422`.
- README enthält alle geforderten Abschnitte.
- Agent-Log und Git-Historie sind nachvollziehbar.
- Optionale OIDC-Funktion ist entweder getestet oder ehrlich als nicht umgesetzt dokumentiert.
