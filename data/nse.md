# NSE Script Search

Discover Nmap NSE scripts by category.

## List NSE categories
`locate .nse | xargs grep 'categories' | grep -oP '".*?"' | sort -u`

Scrape every category string from the local NSE script collection. Run `sudo updatedb` first if `locate` returns nothing.
