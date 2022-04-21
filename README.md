## GoTerraClouDNS

Terraform record dumping tool for ClouDNS zones.

Requires a `sub-auth-id` and `auth-password` which can be set up here: https://www.cloudns.net/api-settings/

This tool will create a terraform file with DNS resources and the corresponding `tf import` commands.

```
Usage of ./GoTerraClouDNS:
  -auth-password string
    	Auth Password
  -out string
    	Output file path (default "import.tf")
  -sub-auth-id int
    	Sub Auth ID
  -zone string
    	Zone name
```

Example output:
```
terraform import cloudns_dns_record.example_com_257659547 cloudns_dns_record/example_com_257659547
terraform import cloudns_dns_record.example_com_257659534 cloudns_dns_record/example_com_257659534
terraform import cloudns_dns_record.example_com_257659541 cloudns_dns_record/example_com_257659541
[...]
```

Example import.tf:
```
resource "cloudns_dns_record" "example_com_257659539" {
  # www.example.com 86400 in CNAME example.com
  name  = "www"
  zone  = "example.com"
  type  = "CNAME"
  value = "example.com"
  ttl   = "86400"
}

resource "cloudns_dns_record" "example_com_257659533" {
  # .example.com 86400 in A 192.0.2.16
  name  = ""
  zone  = "example.com"
  type  = "A"
  value = "192.0.2.16"
  ttl   = "86400"
}
```