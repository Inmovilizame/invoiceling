Output create_invoice.gif

Require echo

Set Shell zsh
Set FontSize 24
Set Width 1200
Set Height 600

Type "cd /tmp/Invoicing"
Enter
Sleep 0.5
Type "invoiceling client"
Enter
Sleep 0.5
Type "invoiceling invoice create \"
Enter
Type '--client client_id \'
Enter
Type "--due 28 \"
Enter
Type '--vat 21 \'
Enter
Type '--retention 15'
Enter
Sleep 1
Type "invoiceling invoice"
Enter
Sleep 1
Type "ls -las invoice"
Enter
Sleep 1
Type "cat invoice/F24-001.json"
Enter
Sleep 3





