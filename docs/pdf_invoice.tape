Output pdf_invoice.gif

Require echo

Set Shell zsh
Set FontSize 24
Set Width 1200
Set Height 600

Type "cd /tmp/Invoicing"
Enter
Sleep 0.5
Type "invoiceling invoice"
Enter
Sleep 0.5
Type "invoiceling pdf \"
Enter
Type '--invoice F24-001'
Enter
Type "ls -las pdf"
Enter
Sleep 1
Type "invoiceling pdf \"
Enter
Type '--invoice F24-001 \'
Enter
Type '--draft'
Enter
Type "ls -las pdf"
Enter
Sleep 3





