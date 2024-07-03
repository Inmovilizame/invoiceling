Output create_client.gif

Require echo

Set Shell zsh
Set FontSize 24
Set Width 1200
Set Height 600

Type "cd /tmp/Invoicing"
Enter
Sleep 0.5
Type "invoiceling client create \"
Enter
Type '--name "Company Name" \'
Enter
Type "--vat_id ES12345678Z \"
Enter
Type '--address1 "Streeet 10" \'
Enter
Type '--address2 "City, Country, Zip Code" \'
Enter
Type "--id client_id"
Enter
Sleep 1
Type "invoiceling client"
Enter
Sleep 1
Type "ls -las client"
Enter
Sleep 1
Type "cat client/client_id.json"
Enter
Sleep 3
