#!/bin/bash

echo "Docs..."
cd ..
pandoc Dokumentation.md -f markdown -t latex --pdf-engine=xelatex -s -o Dokumentation.pdf --include-in-header=.docs/header.tex --lua-filter=.docs/files.lua -V fontsize=12pt -M lang:de

for aufgabe in "a1-Woerter-aufraeumen" "a2-Dreieckspuzzle" "a3-Tobis-Turnier" "a5-Wichteln" ; do
echo "$aufgabe"
cd "$aufgabe"
echo "Build..."
go build main.go
cd beispieldaten
if [ "$1" == "run" ]; then
echo "Run..."
for i in *.txt; do (../main "$i") > "../loesungen/${i}"; done
fi
cd ..
echo "Docs..."
pandoc Dokumentation.md -f markdown -t latex --pdf-engine=xelatex -s -o Dokumentation.pdf --include-in-header=../.docs/header.tex --lua-filter=../.docs/files.lua -V fontsize=12pt -M lang:de
cd ..
done
echo "Zip..."
zip -r "../bwinf39-runde1.zip" "$(pwd)" > /dev/null
echo "Done"