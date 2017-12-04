#!/bin/bash

# root nav
echo "Writing nav.md"
mkpage nav.tmpl relroot="text:" \
    install="text:install.html" \
    docs="text:docs/peanut/" \
    gitrepo="text:https://github.com/rsdoiel/nuts" \
    >"nav.md"

# docs
# docs/peanuts
for D in "docs/" "docs/peanut/"; do
    echo "Writing ${D}nav.md"
    RELPATH=$(reldocpath "${D}" .)
    mkpage nav.tmpl relroot="text:${RELPATH}" \
    install="text:install.html" \
    docs="text:docs/peanuts/" \
    gitrepo="text:https://github.com/rsdoiel/nuts" \
    >"${D}nav.md"
done
