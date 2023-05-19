## mso-pdf-renderer
#### Remote Microsoft Office PDF Renderer

---

### Introduction
Mobile versions of Microsoft Office sometimes renderers documents differently than the desktop versions.
This is a problem when you need to render a document on a server and then display it on a mobile device.
This project is a simple web service that renders Microsoft Office documents to PDF.

### Description
This is a simple web service that renders Microsoft Office documents to PDF.
It is based on the remote CSCRIPT execution.
It is designed to be run on a server with Microsoft Office installed.

### Requirements

 - Microsoft Windows
 - Microsoft Office
 - Windows Auto-Logon Configured

Note that the service will not work if its not ran by an actual user.
Thats why we need to configure the auto-logon.
IDK why i cannot call microsoft office as a service ~~Maybe i just sucks~~

### Deployment

 - Pull the code and compile it with golang
 - Copy the compiled binary with ```static``` and ```scripts``` folder to the server
 - Configure the auto-logon and autorun the binary
 - Its done

### Something else

 - It works anyways, but its junk XD