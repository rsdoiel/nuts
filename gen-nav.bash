#!/bin/bash

GIT_REPO="https://github.com/rsdoiel/nuts"

# root folder
# docs
# how to
for D in "" "docs/" "how-to/"; do
	echo "Writing ${D}nav.md"
	if [[ "$D" == "" ]]; then
		RELPATH=""
	else
		RELPATH=$(reldocpath "${D}" .)
	fi
	mkpage nav.tmpl relroot="text:${RELPATH}" \
		readme="text:index.html" \
		license="text:license.html" \
		install="text:install.html" \
		docs="text:docs/" \
		howto="text:how-to/" \
		gitrepo="text:${GIT_REPO}" \
		>"${D}nav.md"
done
