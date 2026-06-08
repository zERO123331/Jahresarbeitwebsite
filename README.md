
Zum Ausführen der Website wird [Docker](https://docs.docker.com/get-started/get-docker/) benötigt.

Um das Programm mit allen Dependencies zu starten, muss dieser Befehl in dem Folder Projekts ausgeführt werden:
```
docker compose up --build
```
Das restliche Setup sollte automatisch ablaufen
sobald in der Konsole diese Zeile:
```
web-1                   | time=2026-06-08T17:01:43.989Z level=INFO msg=request ip=[::1]:43048 proto=HTTP/1.1 method=GET uri=/healthcheck
```
zusehen ist (der timestamp wird anders sein), sollte die Website gestartet sein. Dann ist die Website unter localhost:4000 aufrufbar.
Die Website ist zwar HTTP allerdings, da es eine localhost connection ist, sollte das keine Fehler erzeugen.
