# passcheck

Check passwords against [Pwned Passwords](https://haveibeenpwned.com/Passwords)

## install
```bash
go install github.com/mokytis/passcheck@main
```

## usage

### passwords from stdin

```bash
$ cat passwords.txt
password1
averysecurepassword
thispasswordissecret
secret
passcheck4lyfe

$ cat passwords.txt | passcheck
thispasswordissecret
password1
secret
```

### passwords from a file

```bash
$ cat passwords.txt
password1
averysecurepassword
thispasswordissecret
secret
passcheck4lyfe

$ passcheck < passwords.txt
thispasswordissecret
password1
secret
```

### show the password count

The `-c` flag shows you how many times each password has been pwned

```bash
$ passcheck -c < passwords.txt
password1:2413945
secret:235493
thispasswordissecret:14
```

### we must go faster

The `-w` flag lets you increase the amount of go workers

```bash
$ passcheck -w 30 < passwords.txt
```
