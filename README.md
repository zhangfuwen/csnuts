csnuts
======
A blog program based on golang and google app engine.

Upload: 
        git push -u origin master
local run:
        dev_server.py .
download from GAE:_
        appcfg.py download_app -A usrccapp -V 1-12 .
update to GAE:
        app_cfg update ./csnuts

To install sessions:
cd csnuts
git clone https://github.com/gorilla/sessions.git

