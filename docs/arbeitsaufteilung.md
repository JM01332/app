# Arbeitsaufteilung und gemeinsame Verträge

Diese Datei legt fest, wie zwei Personen parallel an der Carrier-REST-API arbeiten. Vor jeder Phase müssen die zugehörigen gemeinsamen Entscheidungen getroffen sein. Offene Punkte werden nicht stillschweigend von einer Person allein festgelegt.

## 1. Bekannte technische Grundlage

- Go, Gin, GORM und PostgreSQL
- Bestehende Carrier-Datenbank aus dem vorherigen FastAPI-Projekt
- PostgreSQL läuft weiterhin dort über Docker und ist lokal über `localhost:5432` erreichbar
- Datenbank, Schema und Rolle heißen `carrier`
- Vorhandene Tabellen: `carrier`, `command_center` und `aircraft`
- Keine DB-, Compose- oder Zertifikatsdateien in dieses Repository kopieren
- Kein `AutoMigrate` und keine Schemaänderung durch die Go-Anwendung
- Zugangsdaten ausschließlich in der ignorierten `.env`

## 2. Pflichtentscheidungen vor dem Coding

### Endpunkte

- [ ] GET-Pfad: `________________________________________`
- [ ] POST-Pfad: `_______________________________________`
- [ ] Nur Listen-GET oder zusätzlich Detail-GET festgelegt
- [ ] Verhalten bei leerer Liste festgelegt
- [ ] Sortierung der Liste festgelegt

### POST-Inhalt

- [ ] Wird `commandCenter` zwingend mit angelegt?
- [ ] Werden `aircrafts` direkt mit angelegt?
- [ ] Darf die Aircraft-Liste leer sein?
- [ ] Pflichtfelder und maximale Textlängen festgelegt
- [ ] Unbekannte JSON-Felder ablehnen oder ignorieren?
- [ ] Carrier und Beziehungen gemeinsam in einer Transaktion speichern?

Diskussionsentwurf für den Request:

```json
{
  "name": "",
  "nation": "",
  "carrierType": "AIRCRAFT_CARRIER",
  "commandCenter": {
    "codeName": "",
    "securityLevel": 1
  },
  "aircrafts": [
    {
      "model": "",
      "manufacturer": ""
    }
  ]
}
```

> Dieses JSON wird erst nach gemeinsamer Bestätigung zum verbindlichen Vertrag.

### Response und Fehler

- [ ] Response-Felder und einheitliches `camelCase` festgelegt
- [ ] Beziehungen in GET-Antworten: ja / nein
- [ ] `version`, `erzeugt` und `aktualisiert` in Responses: ja / nein
- [ ] Format für Validierungsfehler festgelegt
- [ ] Ungültiges JSON: `400 Bad Request`
- [ ] Doppelter Carrier-Name: `409 Conflict`
- [ ] Fachlich ungültige Felder: `422 Unprocessable Content`
- [ ] Unerwartete Fehler: `500 Internal Server Error` ohne interne Details

## 3. Verbindlicher Datenbankvertrag

### `carrier`

| DB-Spalte | Typ | Regel |
|---|---|---|
| `id` | `INTEGER` | Identity ab 1000, DB-generiert |
| `version` | `INTEGER` | Pflichtfeld, Standardwert 0 |
| `name` | `TEXT` | Pflichtfeld, eindeutig |
| `nation` | `TEXT` | Pflichtfeld |
| `carrier_type` | Enum | `AIRCRAFT_CARRIER` oder `HELICOPTER_CARRIER` |
| `erzeugt` | `TIMESTAMP` | Pflichtfeld, von der Anwendung zu setzen |
| `aktualisiert` | `TIMESTAMP` | Pflichtfeld, von der Anwendung zu setzen |

### `command_center`

| DB-Spalte | Typ | Regel |
|---|---|---|
| `id` | `INTEGER` | Identity ab 1000, DB-generiert |
| `code_name` | `TEXT` | Pflichtfeld |
| `security_level` | `INTEGER` | Pflichtfeld, Wert 1 bis 5 |
| `carrier_id` | `INTEGER` | Eindeutiger 1:1-Fremdschlüssel auf `carrier` |

### `aircraft`

| DB-Spalte | Typ | Regel |
|---|---|---|
| `id` | `INTEGER` | Identity ab 1000, DB-generiert |
| `model` | `TEXT` | Pflichtfeld |
| `manufacturer` | `TEXT` | Pflichtfeld |
| `carrier_id` | `INTEGER` | n:1-Fremdschlüssel auf `carrier` |

### Vor dem GORM-Modell entscheiden

- [ ] Struct- und Feldnamen abgestimmt
- [ ] Bestehende Tabellen- und Spaltennamen werden exakt abgebildet
- [ ] `CarrierType` wird als eigener Go-Stringtyp modelliert
- [ ] 1:1-Beziehung zu `CommandCenter` festgelegt
- [ ] 1:n-Beziehung zu `Aircrafts` festgelegt
- [ ] Werden Beziehungen bei GET mit `Preload` geladen?
- [ ] Wer setzt `erzeugt` und `aktualisiert`?
- [ ] Wird ein vollständiger POST in einer Transaktion gespeichert?

Wichtig: Das bestehende SQL-Schema definiert für `erzeugt` und `aktualisiert` keinen DB-Standardwert. Die Anwendung muss beide Werte beim Insert setzen.

## 4. Gemeinsamer Service-Vertrag

Vor Parallelphase 2 Methodennamen und Typen gemeinsam festlegen:

```go
type CarrierService interface {
    List(ctx context.Context) ([]model.Carrier, error)
    Create(ctx context.Context, input CreateCarrierInput) (model.Carrier, error)
}
```

- [ ] Rückgabe als Wert oder Pointer festgelegt
- [ ] Package für `CreateCarrierInput` festgelegt
- [ ] Bekannte Service-Fehler festgelegt
- [ ] Darstellung eines Namenskonflikts festgelegt
- [ ] Zuständigkeit für Model-zu-Response-Mapping festgelegt

Das Router-Package darf ein kleines Interface für genau die benötigten Methoden definieren. Go erkennt implizit, ob der echte Service dieses Interface erfüllt. Dadurch kann der Router mit einem Fake-Service ohne PostgreSQL getestet werden.

## 5. Parallelphase 1

### Person A – GORM-Modelle

Exklusive Dateien:

```text
internal/carrier/model/carrier.go
internal/carrier/model/carrier_type.go
internal/carrier/model/command_center.go
internal/carrier/model/aircraft.go
```

- [ ] Structs für alle drei Tabellen erstellen
- [ ] Primär- und Fremdschlüssel abbilden
- [ ] Spaltennamen über GORM-Tags abbilden
- [ ] Enum als eigenen Go-Typ definieren
- [ ] Beziehungen und Zeitstempel abbilden
- [ ] Keine JSON-/Handler-Verantwortung in DB-Modelle mischen
- [ ] Keine Migration ausführen
- [ ] `gofmt` und `go test ./...` ausführen

Commit-Vorschlag: `feat: add carrier GORM models`

### Person B – API-DTOs und Dokumentation

Exklusive Dateien:

```text
internal/carrier/router/create_request.go
internal/carrier/router/carrier_response.go
internal/carrier/router/validation_test.go
docs/api.md
requests/*
```

- [ ] Create-Request nach dem vereinbarten JSON definieren
- [ ] Response-DTO definieren
- [ ] JSON- und Validator-Tags ergänzen
- [ ] Sicherheitsstufe auf 1 bis 5 begrenzen
- [ ] Carrier-Typ auf bekannte Enum-Werte begrenzen
- [ ] Positive und negative Validierungsfälle testen
- [ ] Endpunkte und Beispielrequests dokumentieren
- [ ] Noch keinen DB-Zugriff implementieren
- [ ] `app.go` und `main.go` nicht verändern
- [ ] `gofmt` und `go test ./...` ausführen

Commit-Vorschlag: `feat: add carrier API DTOs`

## 6. Übergabe vor Parallelphase 2

- [ ] Beide Änderungen gegenseitig geprüft
- [ ] `go test ./...` nach dem Zusammenführen erfolgreich
- [ ] Modelle stimmen mit `docs/db.md` und dem SQL-Schema überein
- [ ] DTOs stimmen mit dem vereinbarten JSON überein
- [ ] Service-Methoden und Fehler gemeinsam festgelegt

## 7. Parallelphase 2

### Person A – Repository und Service

Exklusive Dateien:

```text
internal/carrier/service/repository.go
internal/carrier/service/service.go
internal/carrier/service/errors.go
```

- [ ] Carrier-Liste mit GORM lesen
- [ ] Vereinbarte Beziehungen laden
- [ ] Carrier und Beziehungen atomar anlegen
- [ ] Request-Context mit GORM `WithContext` verwenden
- [ ] Eindeutigkeitsverletzung für `name` erkennen
- [ ] DB-Fehler in kontrollierte Service-Fehler übersetzen
- [ ] Tests ergänzen

Commit-Vorschlag: `feat: add carrier repository and service`

### Person B – REST-Handler

Exklusive Dateien:

```text
internal/carrier/router/handler.go
internal/carrier/router/handler_test.go
internal/carrier/router/mapper.go
```

- [ ] GET- und POST-Handler implementieren
- [ ] Request binden und validieren
- [ ] Service über das vereinbarte Interface ansprechen
- [ ] Modell in Response-DTO mappen
- [ ] Statuscodes und Fehlerformat umsetzen
- [ ] Handler mit Fake-Service testen
- [ ] Keine echte DB im Handler-Test verwenden

Commit-Vorschlag: `feat: add carrier REST handlers`

## 8. Gemeinsame Integration

Nur eine Person bearbeitet während der Integration `app.go` und `main.go`.

- [ ] Verantwortlich für `internal/app/app.go`: ____________________
- [ ] Verantwortlich für `cmd/api/main.go`: ________________________
- [ ] DB-Verbindung beim Start öffnen und beim Beenden schließen
- [ ] Repository, Service und Handler zusammensetzen
- [ ] Carrier-Routen registrieren
- [ ] `/health` weiterhin ohne DB-Abfrage betreiben
- [ ] Echten GET gegen die laufende DB prüfen
- [ ] Gültigen und ungültigen POST prüfen
- [ ] Namenskonflikt prüfen
- [ ] `go test ./...` ausführen

Commit-Vorschläge:

```text
feat: wire carrier API components
test: add carrier API integration tests
```

## 9. Regeln gegen Konflikte

- Eine Datei hat während einer Parallelphase genau eine verantwortliche Person.
- JSON, Service-Methoden und Fehler werden vor der Implementierung festgelegt.
- Änderungen an `go.mod`, `go.sum`, `app.go` und `main.go` vorher abstimmen.
- Vor Übergaben immer `gofmt` und `go test ./...` ausführen.
- Keine `.env`, Verbindungs-URLs, Passwörter oder Tokens committen.
- Keine bestehenden Tabellen löschen, migrieren oder automatisch verändern.

## 10. Definition of Ready

Die parallele Implementierung beginnt erst, wenn:

- [ ] GET- und POST-Pfade feststehen
- [ ] POST- und Response-JSON feststehen
- [ ] Pflichtfelder und Validierungsregeln feststehen
- [ ] Umgang mit Command Center und Aircrafts feststeht
- [ ] Zeitstempel und Transaktion festgelegt sind
- [ ] Fehlerformat und Statuscodes feststehen
- [ ] Service-Methoden feststehen
- [ ] Dateiverantwortung namentlich eingetragen ist
- [ ] Beide Personen vom gleichen geprüften Stand arbeiten
