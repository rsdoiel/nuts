#!/bin/bash

GIT_REPO="https://github.com/rsdoiel/nuts"

function write_nav() {
	D=$(dirname "$1")
	echo "Writing ${D}/nav.md"
	if [[ "$D" == "." ]]; then
		RELPATH=""
	else
		RELPATH=$(reldocpath "${D}/" .)
	fi
	mkpage nav.tmpl relroot="text:${RELPATH}" \
		readme="text:index.html" \
		license="text:license.html" \
		install="text:install.html" \
		docs="text:docs/" \
		howto="text:how-to/" \
		gitrepo="text:${GIT_REPO}" \
		>"${D}/nav.md"
}

# root folder
write_nav README.md
# docs/ and how-to/
findfile -c index.md | while read FNAME; do
	write_nav "$FNAME"
done
