# PostgreSQL-Anbindung

## Bestehende Datenbank

Dieses Go-Projekt verwendet die bereits eingerichtete PostgreSQL-Datenbank aus dem vorherigen FastAPI-Projekt. Der PostgreSQL-Container, seine Docker-Volumes, das Schema und die vorhandenen Daten bleiben dort bestehen. In dieses Repository werden keine Datenbankdateien, Compose-Dateien, Zertifikate oder Zugangsdaten kopiert.

Die Go-Anwendung greift ausschließlich als Client über GORM auf die laufende Datenbank zu. `AutoMigrate` wird nicht verwendet, damit das bestehende Schema nicht unbeabsichtigt verändert wird.

## Datenbank starten

Im vorherigen `carrier-api`-Projekt in das Verzeichnis `extras/compose/postgres` wechseln und den vorhandenen Container starten:

```powershell
docker compose up -d db
docker compose ps
```

Der Container veröffentlicht PostgreSQL auf `localhost:5432`. Er kann nach der Arbeit ohne Löschen der bestehenden Docker-Volumes beendet werden:

```powershell
docker compose down
```

## Lokale Konfiguration

Im Go-Projekt `.env.example` nach `.env` kopieren und `PASSWORD` durch das lokale Passwort der bestehenden Carrier-Datenbank ersetzen:

```dotenv
PORT=8080
DATABASE_URL=postgres://carrier:PASSWORD@localhost:5432/carrier?sslmode=require
```

Die Datei `.env` ist über `.gitignore` ausgeschlossen und darf nicht committed werden. Auch vollständige Verbindungs-URLs dürfen nicht geloggt oder in Dokumentationsdateien übernommen werden.

`sslmode=require` erzwingt eine verschlüsselte Verbindung zum TLS-fähigen lokalen PostgreSQL-Container. Eine Prüfung des Serverzertifikats würde zusätzlich einen lokalen Zertifikatspfad benötigen und wird für dieses Workshop-Projekt zunächst nicht eingerichtet.

## Vorhandenes Schema

Die Datenbank und das PostgreSQL-Schema heißen `carrier`. Der bestehende Suchpfad der Rolle `carrier` zeigt auf dieses Schema.

### Tabelle `carrier`

| Spalte | PostgreSQL-Typ | Hinweise |
|---|---|---|
| `id` | `INTEGER` | Identity ab 1000, Primärschlüssel |
| `version` | `INTEGER` | Pflichtfeld, Standardwert 0 |
| `name` | `TEXT` | Pflichtfeld, eindeutig |
| `nation` | `TEXT` | Pflichtfeld |
| `carrier_type` | `carrier_type` | Enum |
| `erzeugt` | `TIMESTAMP` | Pflichtfeld |
| `aktualisiert` | `TIMESTAMP` | Pflichtfeld |

Der Enum `carrier_type` erlaubt:

- `AIRCRAFT_CARRIER`
- `HELICOPTER_CARRIER`

### Tabelle `command_center`

| Spalte | PostgreSQL-Typ | Hinweise |
|---|---|---|
| `id` | `INTEGER` | Identity ab 1000, Primärschlüssel |
| `code_name` | `TEXT` | Pflichtfeld |
| `security_level` | `INTEGER` | Pflichtfeld, Wert von 1 bis 5 |
| `carrier_id` | `INTEGER` | Eindeutiger Fremdschlüssel auf `carrier` |

Zwischen `carrier` und `command_center` besteht eine 1:1-Beziehung.

### Tabelle `aircraft`

| Spalte | PostgreSQL-Typ | Hinweise |
|---|---|---|
| `id` | `INTEGER` | Identity ab 1000, Primärschlüssel |
| `model` | `TEXT` | Pflichtfeld |
| `manufacturer` | `TEXT` | Pflichtfeld |
| `carrier_id` | `INTEGER` | Fremdschlüssel auf `carrier` |

Zwischen `carrier` und `aircraft` besteht eine 1:n-Beziehung. Datensätze in `command_center` und `aircraft` werden durch `ON DELETE CASCADE` entfernt, wenn der zugehörige Carrier gelöscht wird.

## Verbindungsprüfung

Der Go-Code in `internal/database/postgres.go` öffnet die Verbindung mit GORM und ruft anschließend `PingContext` auf. Dadurch wird nicht nur ein Verbindungsobjekt erzeugt, sondern die Erreichbarkeit des PostgreSQL-Servers tatsächlich geprüft. Der darunterliegende Verbindungspool wird beim Beenden der Anwendung mit `Close` freigegeben.

Die Live-Prüfung erfolgt erst, nachdem der Container läuft und eine lokale `.env` mit dem korrekten Passwort vorhanden ist.
