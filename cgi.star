now = time.now()
text = '''
<!DOCTYPE html>
<html>
<head>
    <title>My Homepage</title>
</head>
<body>
    <h1>Welcome to my homepage!</h1>
    <p>Current time is {}.</p>
    <p>Args is <code>{}</code>.</p>
    <pre>Your header: {}</pre>
    <p>This is a simple CGI script written in Starlark.</p>
</body>
</html>
'''.format(now, sys.argv, json.dumps(request.headers, indent=2)).strip()

response.set_html(text)
