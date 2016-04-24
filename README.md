# sendgmail -- CLI tool to send email via Gmail

## For what? 

Sending alert by Email is useful, but there're servers without root/sudo previledge. 

Or you wouldn't want to setup smtp server on your desktop PC or laptop. 

Those are the times you'd need this tool! 


## How to use 

It works like sendmail/mailx. For example: 

	$ echo "this is mail body" | sendgmail -s "Mail subject" -t "Mail TO" -g "Your Gmail acount to send the mail" -p "Your Gmail Password"

## Options 

Here's the list of command line options 

```
Usage of ./sendgmail:
  -c 'mail cc'
  -f 'mail from'
  -g 'Gmail account email is sent from'
  -p 'Gmail password'
  -rm Raw mode
  -s 'subject'
  -t 'mail to'
  -v Verbose
```

## Raw mode 

When you want to send Email with attachment using sendmail/mailx command,
you'd need to customize mail header and include encoded data. 

Raw mode behaves just like sendmail/mailx. 

For example, mail body is like this:  

encoded data 

	$ base64 /tmp/test.txt
	 dGhpcyBpcyB0ZXN0Cgo=
	
	
and the mail body to send 
	 
	$ cat /tmp/testmail.txt
	From: Sender Name <senderaccount@gmail.com>
	To: Recpt Name <destinationaccount@gmail.com>
	Subject: Test mail
	MIME-Version: 1.0
	Content-Type: multipart/mixed; boundary="12345678"; charset="us-ascii"
	Content-Transfer-Encoding: 7bit
	
	--12345678
	MIME-Version: 1.0
	Content-Type: text/plain; charset="us-ascii"
	Content-Transfer-Encoding: 7bit
	
	This is test message.
	from linux mail command
	
	--12345678
	MIME-Version: 1.0
	Content-Type: text/plain; name="test.txt"
	Content-Transfer-Encoding: base64
	Content-Disposition: attachment; filename="test.txt"
	
	dGhpcyBpcyB0ZXN0Cgo=
	
	--12345678


You send Email with "-rm" option.

	$ cat /tmp/testmail.txt | sendgmail -rm -t destinationaccount@gmail.com -g senderaccount@gmail.com -p senderpassword


## Default Gmail account/password 

Your default Gmail account and password can be set to "sendgmail.conf", which looks like this.

	{
	"Gmailuser":"your_Gmail_account",
	"Gmailpass":"your_Gmail_password"
	}

Rename sendmail.conf-sample to sendmail.conf, edit, and place it in the same directory as sendgmail. 

Alternatively, you can hardcode those by editing the lines below: 

	const (
		default_smtpauthuser = "" // Default(Your) Gmail account
		default_smtpauthpass = "" // Default(Your) Gmail password
	)

Note: 
In case multiple Gmail accounts/passwords are supplied, those are used accordingly to the following priority: 

command line > conf > const (hardcoded values)


## Testing environment

 * OSX 10.10 (Mac, x64) 
 * Rasbian (Debian 7, ARM)

