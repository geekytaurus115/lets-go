### Email Verification Tool

#### Points to Know

- `net.LookupMX()` to check for MX (Mail Exchange) records associated with the domain. If records are found, it sets the `hasMX` flag to true.
- `net.LookupTXT()` to retrieve TXT records associated with the domain. It searches through these records to find an SPF record (starting with "v=spf1"). If found, it sets the `hasSPF` flag to true and stores the SPF record.
- the "\_dmarc" domain and uses `net.LookupTXT()` to retrieve TXT records associated with it. It searches through these records to find a DMARC record (starting with "v=DMARC1"). If found, it sets the `hasDMARC` flag to true and stores the DMARC record.
