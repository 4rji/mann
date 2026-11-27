<!-- es: Inspección de datos FTP con TShark -->
# FTP Payload Inspection with Tshark

Reconstruct binary content from FTP traffic captured to a pcap.

<!-- es: Extraer los datos de la captura -->
## Extract payload
`tshark -r 0.pcap -Y 'ftp' -Tfields -e tcp.payload 2>/dev/null | xxd -ps -r`

Filter FTP frames, dump the TCP payload as hex, then decode hex back to binary.
